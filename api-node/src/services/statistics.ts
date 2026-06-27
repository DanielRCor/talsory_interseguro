import { DIAGONAL_TOLERANCE } from "../constants";
import { ApiError } from "../errors";
import { Matrix, StatisticsRequestBody, StatisticsResponse } from "../types";

function validateMatrix(name: string, matrix: unknown): asserts matrix is Matrix {
  if (!Array.isArray(matrix) || matrix.length === 0) {
    throw new ApiError(422, "VALIDATION_ERROR", `${name} matrix is required`, [
      `${name} must be a non-empty array of arrays`,
    ]);
  }

  if (!Array.isArray(matrix[0]) || matrix[0].length === 0) {
    throw new ApiError(422, "VALIDATION_ERROR", `${name} matrix must have non-empty rows`, [
      `${name} row 0 must contain numeric values`,
    ]);
  }

  const expectedLength = matrix[0].length;

  matrix.forEach((row, rowIndex) => {
    if (!Array.isArray(row)) {
      throw new ApiError(422, "VALIDATION_ERROR", `${name} matrix must be rectangular`, [
        `${name} row ${rowIndex} is not an array`,
      ]);
    }

    if (row.length !== expectedLength) {
      throw new ApiError(422, "VALIDATION_ERROR", `${name} matrix must be rectangular`, [
        `${name} row ${rowIndex} has length ${row.length} but expected ${expectedLength}`,
      ]);
    }

    row.forEach((value, columnIndex) => {
      if (typeof value !== "number" || !Number.isFinite(value)) {
        throw new ApiError(422, "VALIDATION_ERROR", `${name} matrix contains invalid values`, [
          `${name}[${rowIndex}][${columnIndex}] must be a finite number`,
        ]);
      }
    });
  });
}

function isDiagonalMatrix(matrix: Matrix): boolean {
  if (matrix.length !== matrix[0].length) {
    return false;
  }

  for (let rowIndex = 0; rowIndex < matrix.length; rowIndex += 1) {
    for (let columnIndex = 0; columnIndex < matrix[rowIndex].length; columnIndex += 1) {
      if (rowIndex !== columnIndex && Math.abs(matrix[rowIndex][columnIndex]) > DIAGONAL_TOLERANCE) {
        return false;
      }
    }
  }

  return true;
}

function flattenMatrices(matrices: Matrix[]): number[] {
  return matrices.flatMap((matrix) => matrix.flat());
}

export function calculateStatistics(body: StatisticsRequestBody): StatisticsResponse {
  if (!body.matrices) {
    throw new ApiError(422, "VALIDATION_ERROR", "Matrices payload is required", [
      "matrices must include q and r",
    ]);
  }

  validateMatrix("q", body.matrices.q);
  validateMatrix("r", body.matrices.r);

  const matrices: Array<[string, Matrix]> = [
    ["q", body.matrices.q],
    ["r", body.matrices.r],
  ];

  const values = flattenMatrices(matrices.map(([, matrix]) => matrix));
  const sum = values.reduce((accumulator, value) => accumulator + value, 0);
  const average = sum / values.length;
  const diagonalMatrices = matrices
    .filter(([, matrix]) => isDiagonalMatrix(matrix))
    .map(([name]) => name);

  return {
    max: Math.max(...values),
    min: Math.min(...values),
    average,
    sum,
    hasDiagonalMatrix: diagonalMatrices.length > 0,
    diagonalMatrices,
  };
}
