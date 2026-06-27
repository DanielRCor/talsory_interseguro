import request from "supertest";
import { createApp } from "../src/app";

describe("node-statistics-api", () => {
  const app = createApp();

  it("returns health status", async () => {
    const response = await request(app).get("/health");

    expect(response.status).toBe(200);
    expect(response.body).toEqual({
      status: "ok",
      service: "node-statistics-api",
    });
  });

  it("returns statistics for a valid payload", async () => {
    const response = await request(app)
      .post("/api/v1/statistics")
      .send({
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

    expect(response.status).toBe(200);
    expect(response.body.max).toBe(4);
    expect(response.body.min).toBe(0);
    expect(response.body.sum).toBe(11);
    expect(response.body.hasDiagonalMatrix).toBe(true);
    expect(response.body.diagonalMatrices).toEqual(["q"]);
  });

  it("returns 422 for invalid payloads", async () => {
    const response = await request(app)
      .post("/api/v1/statistics")
      .send({
        matrices: {
          q: [[1, 2]],
        },
      });

    expect(response.status).toBe(422);
    expect(response.body.error.code).toBe("VALIDATION_ERROR");
  });
});
