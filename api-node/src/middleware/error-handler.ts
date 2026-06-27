import { NextFunction, Request, Response } from "express";
import { ApiError } from "../errors";
import { ApiErrorShape } from "../types";

export function errorHandler(
  error: Error,
  _request: Request,
  response: Response<ApiErrorShape>,
  _next: NextFunction,
): void {
  if (error instanceof ApiError) {
    response.status(error.statusCode).json({
      error: {
        code: error.code,
        message: error.message,
        details: error.details,
      },
    });
    return;
  }

  response.status(500).json({
    error: {
      code: "INTERNAL_SERVER_ERROR",
      message: "Unexpected server error",
    },
  });
}
