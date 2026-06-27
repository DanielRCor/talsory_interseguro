import express, { NextFunction, Request, Response } from "express";
import { errorHandler } from "./middleware/error-handler";
import { calculateStatistics } from "./services/statistics";

export function createApp() {
  const app = express();

  app.use(express.json());

  app.get("/health", (_request, response) => {
    response.json({
      status: "ok",
      service: "node-statistics-api",
    });
  });

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
