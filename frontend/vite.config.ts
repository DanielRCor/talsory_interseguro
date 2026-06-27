import { defineConfig } from "vite";

export default defineConfig({
  server: {
    host: "0.0.0.0",
    port: 5173,
    proxy: {
      "/go-api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/go-api/, ""),
      },
      "/node-api": {
        target: "http://localhost:3000",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/node-api/, ""),
      },
    },
  },
  preview: {
    host: "0.0.0.0",
    port: 5173,
  },
});
