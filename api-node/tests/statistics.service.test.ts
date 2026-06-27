import { ApiError } from "../src/errors";
import { calculateStatistics } from "../src/services/statistics";

describe("calculateStatistics", () => {
  it("returns aggregate statistics for q and r", () => {
    const result = calculateStatistics({
      matrices: {
        q: [
          [1, 0],
          [0, 1],
        ],
        r: [
          [2, 3],
          [0, 4],
        ],
      },
    });

    expect(result).toEqual({
      max: 4,
      min: 0,
      average: 1.375,
      sum: 11,
      hasDiagonalMatrix: true,
      diagonalMatrices: ["q"],
    });
  });

  it("supports negative and decimal values", () => {
    const result = calculateStatistics({
      matrices: {
        q: [[-1.5, 0]],
        r: [[2.5, 4]],
      },
    });

    expect(result.max).toBe(4);
    expect(result.min).toBe(-1.5);
    expect(result.sum).toBe(5);
    expect(result.average).toBe(1.25);
    expect(result.hasDiagonalMatrix).toBe(false);
  });

  it("flags diagonal matrices using tolerance", () => {
    const result = calculateStatistics({
      matrices: {
        q: [
          [1, 1e-10],
          [0, 2],
        ],
        r: [
          [2, 0],
          [0, 3],
        ],
      },
    });

    expect(result.hasDiagonalMatrix).toBe(true);
    expect(result.diagonalMatrices).toEqual(["q", "r"]);
  });

  it("rejects missing matrices payload", () => {
    expect(() => calculateStatistics({})).toThrow(ApiError);
  });

  it("rejects non-numeric matrix values", () => {
    expect(() =>
      calculateStatistics({
        matrices: {
          q: [[1, Number.NaN]],
          r: [[2, 3]],
        },
      }),
    ).toThrow(ApiError);
  });
});
