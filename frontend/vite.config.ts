import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      // AJAX: http://localhost:3000/api/... -> Gin http://localhost:8080/api/...
      "/api": { target: process.env.VITE_PROXY_API_TARGET || "http://localhost:8080", changeOrigin: true, configure: proxy => {
        proxy.on("error", (_error, _req, res) => {
          if (res && "writeHead" in res && "end" in res && !res.headersSent) {
            res.writeHead(503, { "Content-Type": "application/json", "X-Lab6-Proxy-Offline": "1" });
            res.end(JSON.stringify({ error: "Go API is offline" }));
          }
        });
      } },
      // S3/MinIO: http://localhost:3000/s3/bucket/file -> http://localhost:9000/bucket/file
      "/s3": {
        target: process.env.VITE_PROXY_S3_TARGET || "http://localhost:9000",
        changeOrigin: true,
        rewrite: (path: string) => path.replace(/^\/s3(?=\/|$)/, ""),
      },
    },
  },
});
