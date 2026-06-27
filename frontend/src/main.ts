import "./style.css";

type Matrix = number[][];

type ServiceHealth = {
  status: string;
  service: string;
};

type ApiError = {
  error?: {
    code?: string;
    message?: string;
    details?: string[];
  };
};

const goApiBase = import.meta.env.VITE_GO_API_BASE ?? "http://localhost:8080";
const nodeApiBase = import.meta.env.VITE_NODE_API_BASE ?? "http://localhost:3000";

const app = document.querySelector<HTMLDivElement>("#app");

if (!app) {
  throw new Error("App root not found");
}

app.innerHTML = `
  <div class="shell">
    <header class="hero">
      <div class="hero-copy">
        <p class="eyebrow">Interseguro Coding Challenge</p>
        <h1>Matrix Lab</h1>
        <p class="lede">
          A focused browser client for testing QR factorization, downstream statistics,
          and the optional rotation flow without leaving localhost.
        </p>
      </div>
      <div class="hero-panel">
        <div class="status-card" data-health="go">
          <span>Go API</span>
          <strong>Checking...</strong>
        </div>
        <div class="status-card" data-health="node">
          <span>Node API</span>
          <strong>Checking...</strong>
        </div>
      </div>
    </header>

    <main class="grid">
      <section class="panel control-panel">
        <div class="panel-heading">
          <p class="eyebrow">Request Builder</p>
          <h2>Drive the APIs from one place</h2>
        </div>

        <label class="field">
          <span>JWT token (optional)</span>
          <input id="jwtToken" type="text" placeholder="Paste bearer token if ENABLE_AUTH=true" />
        </label>

        <label class="field">
          <span>Matrix JSON</span>
          <textarea id="matrixInput" spellcheck="false">[[1,2],[3,4],[5,6]]</textarea>
        </label>

        <div class="action-grid">
          <button id="analyzeButton" class="action primary">Run QR + Statistics</button>
          <button id="rotateClockwiseButton" class="action secondary">Rotate Clockwise</button>
          <button id="rotateCounterButton" class="action secondary">Rotate Counterclockwise</button>
          <button id="refreshHealthButton" class="action ghost">Refresh Health</button>
        </div>

        <p class="hint">
          QR requires <code>rows &gt;= columns</code>. Rotation accepts any rectangular matrix.
        </p>
      </section>

      <section class="panel result-panel">
        <div class="panel-heading">
          <p class="eyebrow">Response</p>
          <h2>Live API output</h2>
        </div>
        <div id="flash" class="flash" hidden></div>
        <pre id="responseOutput" class="response">(results will appear here)</pre>
      </section>

      <section class="panel quick-panel">
        <div class="panel-heading">
          <p class="eyebrow">Playbook</p>
          <h2>Suggested checks</h2>
        </div>
        <ul class="quick-list">
          <li>Use the default 3x2 matrix to confirm the full Go → Node flow.</li>
          <li>Try <code>[[1,2,3],[4,5,6]]</code> to see QR validation reject wide matrices.</li>
          <li>Enable JWT in both APIs and paste the token above to test protected routes.</li>
          <li>Switch to rotation when you want to demo the optional interpretation of the prompt.</li>
        </ul>
      </section>
    </main>
  </div>
`;

const matrixInput = document.querySelector<HTMLTextAreaElement>("#matrixInput");
const jwtTokenInput = document.querySelector<HTMLInputElement>("#jwtToken");
const responseOutput = document.querySelector<HTMLElement>("#responseOutput");
const flash = document.querySelector<HTMLElement>("#flash");
const analyzeButton = document.querySelector<HTMLButtonElement>("#analyzeButton");
const rotateClockwiseButton = document.querySelector<HTMLButtonElement>("#rotateClockwiseButton");
const rotateCounterButton = document.querySelector<HTMLButtonElement>("#rotateCounterButton");
const refreshHealthButton = document.querySelector<HTMLButtonElement>("#refreshHealthButton");

function showFlash(kind: "success" | "error", message: string) {
  if (!flash) {
    return;
  }

  flash.hidden = false;
  flash.className = `flash ${kind}`;
  flash.textContent = message;
}

function setResponse(value: unknown) {
  if (!responseOutput) {
    return;
  }

  responseOutput.textContent = JSON.stringify(value, null, 2);
}

function parseMatrix(): Matrix {
  if (!matrixInput) {
    throw new Error("Matrix input not found");
  }

  const parsed = JSON.parse(matrixInput.value);
  if (!Array.isArray(parsed) || !parsed.every((row) => Array.isArray(row))) {
    throw new Error("Matrix must be a JSON array of arrays");
  }

  return parsed as Matrix;
}

function buildHeaders(): HeadersInit {
  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };

  const token = jwtTokenInput?.value.trim();
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  return headers;
}

async function parseJsonResponse<T>(response: Response): Promise<T> {
  const payload = (await response.json()) as T & ApiError;
  if (!response.ok) {
    const message = payload?.error?.message ?? `Request failed with status ${response.status}`;
    throw new Error(message);
  }
  return payload;
}

async function runAnalyze() {
  try {
    const matrix = parseMatrix();
    const response = await fetch(`${goApiBase}/api/v1/qr/analyze`, {
      method: "POST",
      headers: buildHeaders(),
      body: JSON.stringify({ matrix }),
    });

    const payload = await parseJsonResponse<Record<string, unknown>>(response);
    setResponse(payload);
    showFlash("success", "QR factorization and statistics completed.");
  } catch (error) {
    setResponse({ error: String(error) });
    showFlash("error", error instanceof Error ? error.message : "Unexpected error");
  }
}

async function runRotate(direction: "clockwise" | "counterclockwise") {
  try {
    const matrix = parseMatrix();
    const response = await fetch(`${goApiBase}/api/v1/matrix/rotate`, {
      method: "POST",
      headers: buildHeaders(),
      body: JSON.stringify({ matrix, direction }),
    });

    const payload = await parseJsonResponse<Record<string, unknown>>(response);
    setResponse(payload);
    showFlash("success", `Rotation completed (${direction}).`);
  } catch (error) {
    setResponse({ error: String(error) });
    showFlash("error", error instanceof Error ? error.message : "Unexpected error");
  }
}

async function updateHealthCard(selector: string, url: string) {
  const card = document.querySelector<HTMLElement>(selector);
  if (!card) {
    return;
  }

  const strong = card.querySelector("strong");
  try {
    const response = await fetch(url);
    const payload = await parseJsonResponse<ServiceHealth>(response);
    card.dataset.state = "ok";
    if (strong) {
      strong.textContent = `${payload.status} · ${payload.service}`;
    }
  } catch (error) {
    card.dataset.state = "error";
    if (strong) {
      strong.textContent = error instanceof Error ? error.message : "Unavailable";
    }
  }
}

async function refreshHealth() {
  await Promise.all([
    updateHealthCard('[data-health="go"]', `${goApiBase}/health`),
    updateHealthCard('[data-health="node"]', `${nodeApiBase}/health`),
  ]);
}

analyzeButton?.addEventListener("click", () => {
  void runAnalyze();
});

rotateClockwiseButton?.addEventListener("click", () => {
  void runRotate("clockwise");
});

rotateCounterButton?.addEventListener("click", () => {
  void runRotate("counterclockwise");
});

refreshHealthButton?.addEventListener("click", () => {
  void refreshHealth();
});

void refreshHealth();
