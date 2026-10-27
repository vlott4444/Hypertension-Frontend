import { Link, useParams } from "react-router-dom";
import { MOCK_STAGES } from "../data/mock";
import { StageImage, StageVideo } from "../components/StageMedia";
export function DetailPage(){const {id}=useParams();const stage=MOCK_STAGES.find(s=>s.id===Number(id)&&s.status==="Published");
 if(!stage)return <main className="stage-detail"><Link className="back-link" to="/stages">← К плитке</Link><h1>Карточка не найдена</h1></main>;
 return <main className="stage-detail"><Link className="back-link" to="/stages">← Назад к плитке</Link><StageImage stage={stage} className="detail-cover"/><div className="detail-chip">САД {stage.systolic_bp} / ДАД {stage.diastolic_bp}</div><h1>{stage.title}</h1><p>{stage.description}</p><h2>Видео</h2><div className="detail-video"><StageVideo stage={stage}/></div></main>;
}
