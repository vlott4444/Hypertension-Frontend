import { useEffect } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { StageImage, StageVideo } from "../components/StageMedia";
import { StageActions } from "../components/StageActions";
import { fetchHypertensionStage } from "../store/stagesSlice";
import { useAppDispatch, useAppSelector } from "../store/store";

export function DetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const { selected: stage, loading, error } = useAppSelector(s => s.hypertensionStages);
  const isAuthenticated = useAppSelector(s => s.hypertensionAuth.isAuthenticated);
  const numberId = Number(id);
  useEffect(() => {
    if (Number.isInteger(numberId) && numberId > 0) void dispatch(fetchHypertensionStage(numberId));
  }, [dispatch, numberId, isAuthenticated]);
  if ((!stage || stage.id !== numberId) && loading) return <main className="stage-detail"><p>Загрузка карточки...</p></main>;
  if (!stage || stage.id !== numberId) return <main className="stage-detail">
    <Link className="back-link" to="/stages">← К плитке</Link>
    <h1>Карточка не найдена</h1><p className="api-warning">{error}</p>
  </main>;
  return <main className="stage-detail"><Link className="back-link" to="/stages">← Назад к плитке</Link>
    <StageImage stage={stage} className="detail-cover"/>
    <div className="detail-chip">САД {stage.systolic_bp} / ДАД {stage.diastolic_bp}</div>
    <h1>{stage.title}</h1><p>{stage.description}</p>
    <StageActions stage={stage} onDeleted={() => navigate("/stages")}/>
    {error && <p className="api-warning" role="alert">{error}</p>}
    <h2>Видео</h2><div className="detail-video"><StageVideo stage={stage}/></div>
  </main>;
}
