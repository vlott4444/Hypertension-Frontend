import { useEffect, useState, type FormEvent } from "react";
import type { StageFilters } from "../types";

interface Props {
  filters: StageFilters;
  onApply: (filters: StageFilters) => void;
  onReset: () => void;
  disabled: boolean;
}
export function StageFilter({ filters, onApply, onReset, disabled }: Props) {
  const [draftMin, setDraftMin] = useState(filters.minBP);
  const [draftMax, setDraftMax] = useState(filters.maxBP);
  const [draftTitle, setDraftTitle] = useState(filters.title);
  useEffect(() => {
    setDraftMin(filters.minBP);
    setDraftMax(filters.maxBP);
    setDraftTitle(filters.title);
  }, [filters.minBP, filters.maxBP, filters.title]);
  const apply = (e: FormEvent) => {
    e.preventDefault();
    if (draftMin <= draftMax) onApply({ minBP: draftMin, maxBP: draftMax, title: draftTitle });
  };
  return <section className="search-zone"><form className="sbp-slider-form" onSubmit={apply}>
    <label className="filter-title-label" htmlFor="title_search">Поиск по названию</label>
    <input className="title-search" id="title_search" type="search" value={draftTitle}
      onChange={e => setDraftTitle(e.target.value)} placeholder="Например: нормальное" disabled={disabled}/>
    <div className="slider-header"><label htmlFor="sbp_min">Систолическое давление</label>
      <strong className="slider-current-value">{draftMin}–{draftMax}</strong></div>
    <input id="sbp_min" type="range" min={100} max={200} value={draftMin} disabled={disabled}
      onChange={e => setDraftMin(Number(e.target.value))} aria-label="Минимальное систолическое давление"/>
    <input id="sbp_max" type="range" min={100} max={200} value={draftMax} disabled={disabled}
      onChange={e => setDraftMax(Number(e.target.value))} aria-label="Максимальное систолическое давление"/>
    <div className="slider-scale"><span>100</span><span>200 мм рт. ст.</span></div>
    {draftMin > draftMax && <p className="filter-error">Минимальное давление больше максимального.</p>}
    <div className="slider-actions">
      <button className="slider-submit" type="submit" disabled={disabled || draftMin > draftMax}>Применить</button>
      <button className="slider-reset" type="button" onClick={() => {
        setDraftMin(100); setDraftMax(200); setDraftTitle(""); onReset();
      }} disabled={disabled}>Сбросить</button>
    </div>
  </form></section>;
}
