import axios from "axios";

export const TOKEN_STORAGE_KEY = "hypertension.jwt";

/** Same-origin /api is forwarded by the Vite reverse proxy; no browser-side CORS workaround. */
export const http = axios.create({ baseURL: "/api", timeout: 60000 });

http.interceptors.request.use(config => {
  const token = localStorage.getItem(TOKEN_STORAGE_KEY);
  if (token && !/^\/auth\/(?:login|register)$/.test(config.url || "")) {
    config.headers.set("Authorization", `Bearer ${token}`);
  }
  return config;
});

export function apiError(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const data: unknown = error.response?.data;
    if (data && typeof data === "object") {
      const record = data as Record<string, unknown>;
      for (const key of ["error", "message", "detail"]) {
        if (typeof record[key] === "string") return record[key];
      }
    }
    if (error.response?.status === 401) return "Необходимо войти в аккаунт (401).";
    if (error.response?.status === 403) return "Недостаточно прав (403).";
    if (!error.response) return "Бэкенд недоступен. Проверьте, что Go-сервис запущен на порту 8080.";
    return `Ошибка сервера: HTTP ${error.response.status}.`;
  }
  return error instanceof Error ? error.message : "Неизвестная ошибка";
}

/** Offline mock is permitted only for guests and only for transport errors. */
export function isTransportError(error: unknown): boolean {
  return axios.isAxiosError(error) && (!error.response || error.response.headers?.["x-lab6-proxy-offline"] === "1");
}
