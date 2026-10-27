import { useState } from "react";
interface Props { minBP:number; maxBP:number; onMinChange:(value:number)=>void; onMaxChange:(value:number)=>void; onReset:()=>void; }
export function StageFilter({minBP,maxBP,onMinChange,onMaxChange,onReset}:Props) {
  const [draftMin,setDraftMin]=useState(minBP);
  const [draftMax,setDraftMax]=useState(maxBP);
  const apply=(e:React.FormEvent)=>{e.preventDefault();onMinChange(draftMin);onMaxChange(draftMax);};
  const reset=()=>{setDraftMin(100);setDraftMax(200);onReset();};
  return <section className="search-zone"><form className="sbp-slider-form" onSubmit={apply}>
    <div className="slider-header"><label htmlFor="sbp_min">Систолическое давление</label><strong className="slider-current-value">{draftMin}–{draftMax} мм рт. ст.</strong></div>
    <input id="sbp_min" type="range" min={100} max={200} value={draftMin} onChange={e=>setDraftMin(Number(e.target.value))} aria-label="Минимальное систолическое давление"/>
    <input id="sbp_max" type="range" min={100} max={200} value={draftMax} onChange={e=>setDraftMax(Number(e.target.value))} aria-label="Максимальное систолическое давление"/>
    <div className="slider-scale"><span>100</span><span>200</span></div>
    <div className="slider-actions"><button className="slider-submit" type="submit">Применить</button><button className="slider-reset" type="button" onClick={reset}>Сбросить</button></div>
  </form></section>;
}
