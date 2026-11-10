import { useAppSelector } from "../store/store";

export function BusyOverlay() {
  const stagesLoading = useAppSelector(s => s.hypertensionStages.loading);
  const feedLoading = useAppSelector(s => s.hypertensionFeed.loading);
  const authBusy = useAppSelector(s => s.hypertensionAuth.busy);
  if (!stagesLoading && !feedLoading && !authBusy) return null;
  return <div className="busy-overlay" role="status" aria-live="polite" aria-label="Ожидается ответ сервера">
    <div className="busy-box"><span className="busy-spinner"/>Загрузка...</div>
  </div>;
}
