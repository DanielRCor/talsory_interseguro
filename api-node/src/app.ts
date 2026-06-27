import cors from "cors";
import express, { NextFunction, Request, Response } from "express";
import jwt from "jsonwebtoken";
import { ApiError } from "./errors";
import { errorHandler } from "./middleware/error-handler";
import { calculateStatistics } from "./services/statistics";

interface AppOptions {
  enableAuth?: boolean;
  jwtSecret?: string;
}

function readBoolean(rawValue: string | undefined, fallback: boolean): boolean {
  if (rawValue === undefined) {
    return fallback;
  }

  return rawValue === "true";
}

function resolveOptions(options?: AppOptions): Required<AppOptions> {
  return {
    enableAuth: options?.enableAuth ?? readBoolean(process.env.ENABLE_AUTH, false),
    jwtSecret: options?.jwtSecret ?? process.env.JWT_SECRET ?? "change-me-only-if-auth-enabled",
  };
}

export function createApp(options?: AppOptions) {
  const app = express();
  const resolvedOptions = resolveOptions(options);

  app.use(cors());
  app.use(express.json());

  app.get("/health", (_request, response) => {
    response.json({
      status: "ok",
      service: "node-statistics-api",
    });
  });

  if (resolvedOptions.enableAuth) {
    app.use("/api/v1", (request, _response, next) => {
      const authorization = request.header("Authorization");
      if (!authorization?.startsWith("Bearer ")) {
        next(new ApiError(401, "UNAUTHORIZED", "Missing or invalid bearer token"));
        return;
      }

      const token = authorization.slice("Bearer ".length);
      try {
        jwt.verify(token, resolvedOptions.jwtSecret, { algorithms: ["HS256"] });
        next();
      } catch (error) {
        next(new ApiError(401, "UNAUTHORIZED", "Invalid bearer token", [String(error)]));
      }
    });
  }

  app.post("/api/v1/statistics", (request, response, next) => {
    try {
      const statistics = calculateStatistics(request.body);
      response.json(statistics);
    } catch (error) {
      next(error);
    }
  });

  app.use((error: Error, request: Request, response: Response, next: NextFunction) => {
    errorHandler(error, request, response, next);
  });

  return app;
}
