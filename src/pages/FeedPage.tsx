import { useEffect, useState } from "react";
import { MOCK_STAGES } from "../data/mock";
import { StageVideo } from "../components/StageMedia";
const items=MOCK_STAGES.filter(s=>s.status==="Published");
function heading(title:string){const numeral=title.match(/^(I|II|III) стадия/);return numeral? <>{numeral[1]} стадия<br/>гипертонической<br/>болезни</>:title;}
export function FeedPage(){const [index,setIndex]=useState(0);const [expanded,setExpanded]=useState(false);const stage=items[index];
 useEffect(()=>{document.title="Лента — стадии гипертонии";return ()=>{document.title="Стадии гипертонии"}},[]);
 const next=()=>{setExpanded(false);setIndex(i=>(i+1)%items.length)};
 return <main className="feed-shell"><section className="feed-media"><StageVideo key={stage.id} stage={stage}/><div className="feed-shade"/>
 <input type="checkbox" className="desc-toggle" checked={expanded} onChange={e=>setExpanded(e.target.checked)} id="desc-toggle"/>
 <aside className="feed-actions"><div className="feed-like"><span className="heart-icon">♡</span><span>0</span></div><button type="button" className="next-button" onClick={next}>Далее</button></aside>
 <section className="feed-copy"><h1>{heading(stage.title)}</h1><div className="feed-divider"/><div className="feed-description"><p className="desc-text">{stage.description}</p><label htmlFor="desc-toggle" className="more-button"><span className="more-label-open">ещё</span><span className="more-label-close">свернуть</span></label></div></section>
 </section></main>;
}
