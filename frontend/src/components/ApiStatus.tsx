import { useAppSelector } from "../store/store";

export function ApiStatus() {
  const { source, error } = useAppSelector(s => s.hypertensionStages);
  if (error) return <p className="api-warning" role="alert">{error}</p>;
  if (source === "mock") return <p className="mock-notice">Демо-режим: сервер недоступен, показаны mock-карточки для гостя.</p>;
  if (source === "backend") return <p className="api-connected">Данные загружены из Go API</p>;
  return null;
}
