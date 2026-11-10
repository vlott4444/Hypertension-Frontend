import { useEffect } from "react";
import { StageFilter } from "../components/StageFilter";
import { StageCard } from "../components/StageCard";
import { ApiStatus } from "../components/ApiStatus";
import { fetchHypertensionStages } from "../store/stagesSlice";
import { resetFilters, setPressureRange, setTitleFilter } from "../store/filterSlice";
import { useAppDispatch, useAppSelector } from "../store/store";
import type { StageFilters } from "../types";

export function GridPage() {
  const dispatch = useAppDispatch();
  const { items, loading } = useAppSelector(s => s.hypertensionStages);
  const filters = useAppSelector(s => s.hypertensionFilters);
  const isAuthenticated = useAppSelector(s => s.hypertensionAuth.isAuthenticated);
  useEffect(() => {
    document.title = "Плитка — стадии гипертонии";
    void dispatch(fetchHypertensionStages());
    return () => { document.title = "Стадии гипертонии"; };
  }, [dispatch, isAuthenticated]);
  const apply = (f: StageFilters) => {
    dispatch(setPressureRange({ minBP: f.minBP, maxBP: f.maxBP }));
    dispatch(setTitleFilter(f.title));
    // Read the updated Redux filter from the thunk's getState() after both sync actions.
    void dispatch(fetchHypertensionStages());
  };
  const reset = () => { dispatch(resetFilters()); void dispatch(fetchHypertensionStages()); };
  return <main className="grid-shell">
    <StageFilter filters={filters} onApply={apply} onReset={reset} disabled={loading}/>
    <ApiStatus/>
    <section className="stage-grid" aria-label="Опубликованные услуги">
      {items.map(stage => <StageCard key={stage.id} stage={stage}/>)}
      {items.length === 0 && !loading && <p className="empty-result">Карточки не найдены.</p>}
    </section>
  </main>;
}
