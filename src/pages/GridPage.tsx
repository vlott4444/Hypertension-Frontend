import { useEffect, useState } from "react";
import { MOCK_STAGES } from "../data/mock";
import { StageFilter } from "../components/StageFilter";
import { StageCard } from "../components/StageCard";
export function GridPage(){
 const [minBP,setMinBP]=useState(100);const [maxBP,setMaxBP]=useState(200);
 useEffect(()=>{document.title="Плитка — стадии гипертонии";return ()=>{document.title="Стадии гипертонии"}},[]);
 const results=MOCK_STAGES.filter(s=>s.status==="Published"&&s.systolic_bp>=minBP&&s.systolic_bp<=maxBP);
 return <main className="grid-shell"><StageFilter minBP={minBP} maxBP={maxBP} onMinChange={setMinBP} onMaxChange={setMaxBP} onReset={()=>{setMinBP(100);setMaxBP(200)}}/>
 {minBP>maxBP && <p className="filter-error">Минимальное давление больше максимального</p>}
 <section className="stage-grid">{results.map(stage=><StageCard key={stage.id} stage={stage}/>)}{results.length===0&&<p className="empty-result">Для выбранного САД карточка не найдена.</p>}</section></main>
}
