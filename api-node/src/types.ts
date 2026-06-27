export type Matrix = number[][];

export interface StatisticsRequestBody {
  matrices?: {
    q?: Matrix;
    r?: Matrix;
  };
}

export interface StatisticsResponse {
  max: number;
  min: number;
  average: number;
  sum: number;
  hasDiagonalMatrix: boolean;
  diagonalMatrices: string[];
}

export interface ApiErrorShape {
  error: {
    code: string;
    message: string;
    details?: string[];
  };
}
