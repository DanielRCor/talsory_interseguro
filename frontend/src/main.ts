import "./style.css";

type Matrix = number[][];

type ServiceHealth = {
  status: string;
  service: string;
};

type DemoTokenResponse = {
  token: string;
  tokenPreview: string;
  expiresAt: string;
  expiresInSeconds: number;
  enabled: boolean;
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
        <p class="eyebrow">Prueba Tecnica Interseguro</p>
        <h1>Laboratorio de Matrices</h1>
        <p class="lede">
          Cliente web para probar la factorizacion QR, las estadisticas, la rotacion
          y las peticiones autenticadas sin salir de localhost.
        </p>
      </div>
      <div class="hero-panel">
        <div class="status-card" data-health="go">
          <span>Go API</span>
          <strong>Verificando...</strong>
        </div>
        <div class="status-card" data-health="node">
          <span>Node API</span>
          <strong>Verificando...</strong>
        </div>
      </div>
    </header>

    <main class="grid">
      <section class="panel control-panel">
        <div class="panel-heading">
          <p class="eyebrow">Constructor de Peticiones</p>
          <h2>Controla las APIs desde un solo lugar</h2>
        </div>

        <label class="field">
          <span>Token JWT</span>
          <input id="jwtToken" type="text" placeholder="Usa Generar JWT Demo o pega tu propio bearer token" />
        </label>

        <div class="action-grid compact">
          <button id="generateTokenButton" class="action secondary">Generar JWT Demo</button>
          <button id="clearTokenButton" class="action ghost">Limpiar Token</button>
        </div>

        <label class="field">
          <span>Matriz en JSON</span>
          <textarea id="matrixInput" spellcheck="false">[[1,2],[3,4],[5,6]]</textarea>
        </label>

        <div class="action-grid">
          <button id="analyzeButton" class="action primary">Ejecutar QR + Estadisticas</button>
          <button id="rotateClockwiseButton" class="action secondary">Rotar a la Derecha</button>
          <button id="rotateCounterButton" class="action secondary">Rotar a la Izquierda</button>
          <button id="refreshHealthButton" class="action ghost">Actualizar Estado</button>
        </div>

        <p class="hint">
          QR requiere <code>filas &gt;= columnas</code>. La rotacion acepta cualquier matriz rectangular.
        </p>
      </section>

      <section class="panel result-panel">
        <div class="panel-heading">
          <p class="eyebrow">Respuesta</p>
          <h2>Salida en vivo de las APIs</h2>
        </div>
        <div id="flash" class="flash" hidden></div>
        <pre id="responseOutput" class="response">(aqui apareceran los resultados)</pre>
      </section>

      <section class="panel quick-panel">
        <div class="panel-heading">
          <p class="eyebrow">Guia Rapida</p>
          <h2>Pruebas sugeridas</h2>
        </div>
        <ul class="quick-list">
          <li>Genera primero un JWT demo si quieres usar el modo autenticado por defecto.</li>
          <li>Usa la matriz 3x2 por defecto para confirmar el flujo completo Go -> Node.</li>
          <li>Prueba <code>[[1,2,3],[4,5,6]]</code> para ver como QR rechaza matrices anchas.</li>
          <li>Usa rotacion cuando quieras demostrar la interpretacion opcional del enunciado.</li>
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
const generateTokenButton = document.querySelector<HTMLButtonElement>("#generateTokenButton");
const clearTokenButton = document.querySelector<HTMLButtonElement>("#clearTokenButton");

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

function formatDemoTokenResponse(payload: DemoTokenResponse) {
  return {
    autenticacionActiva: payload.enabled,
    expiraEnSegundos: payload.expiresInSeconds,
    expiraEn: payload.expiresAt,
    tokenPreview: payload.tokenPreview,
    nota: "El JWT completo ya fue cargado en el campo superior.",
  };
}

async function generateDemoToken() {
  try {
    const response = await fetch(`${goApiBase}/auth/demo-token`, {
      method: "POST",
    });

    const payload = await parseJsonResponse<DemoTokenResponse>(response);
    if (jwtTokenInput) {
      jwtTokenInput.value = payload.token;
    }

    setResponse(formatDemoTokenResponse(payload));
    showFlash("success", `JWT demo cargado. Expira en ${payload.expiresInSeconds} segundos.`);
  } catch (error) {
    setResponse({ error: String(error) });
    showFlash("error", error instanceof Error ? error.message : "Fallo la generacion del token");
  }
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
    showFlash("success", "Factorizacion QR y estadisticas completadas.");
  } catch (error) {
    setResponse({ error: String(error) });
    showFlash("error", error instanceof Error ? error.message : "Error inesperado");
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
    showFlash("success", `Rotacion completada (${direction}).`);
  } catch (error) {
    setResponse({ error: String(error) });
    showFlash("error", error instanceof Error ? error.message : "Error inesperado");
  }
}

async function updateHealthCard(selector: string, url: string): Promise<ServiceHealth> {
  const card = document.querySelector<HTMLElement>(selector);
  if (!card) {
    throw new Error(`Health card not found for ${selector}`);
  }

  const strong = card.querySelector("strong");
  try {
    const response = await fetch(url);
    const payload = await parseJsonResponse<ServiceHealth>(response);
    card.dataset.state = "ok";
    if (strong) {
      strong.textContent = `${payload.status} - ${payload.service}`;
    }
    return payload;
  } catch (error) {
    card.dataset.state = "error";
    const message = error instanceof Error ? error.message : "No disponible";
    if (strong) {
      strong.textContent = message;
    }
    throw new Error(message);
  }
}

async function refreshHealth() {
  try {
    const [goHealth, nodeHealth] = await Promise.all([
      updateHealthCard('[data-health="go"]', `${goApiBase}/health`),
      updateHealthCard('[data-health="node"]', `${nodeApiBase}/health`),
    ]);

    setResponse({
      goApi: goHealth,
      nodeApi: nodeHealth,
    });
    showFlash("success", "Estado de salud actualizado.");
  } catch (error) {
    setResponse({ error: String(error) });
    showFlash("error", error instanceof Error ? error.message : "Fallo la verificacion de salud");
  }
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

generateTokenButton?.addEventListener("click", () => {
  void generateDemoToken();
});

clearTokenButton?.addEventListener("click", () => {
  if (jwtTokenInput) {
    jwtTokenInput.value = "";
  }
  showFlash("success", "Token limpiado.");
});

void refreshHealth();
