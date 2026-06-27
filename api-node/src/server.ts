import { createApp } from "./app";

const port = Number(process.env.PORT || process.env.NODE_API_PORT || 3000);
const app = createApp();

app.listen(port, () => {
  console.log(`node-statistics-api listening on port ${port}`);
});
