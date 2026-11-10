import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import {
  createHypertensionStage, fetchHypertensionDraft, fetchHypertensionStages,
  publishHypertensionStage,
} from "../store/stagesSlice";
import { useAppDispatch, useAppSelector } from "../store/store";

export function AddPage() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const authenticated = useAppSelector(s => s.hypertensionAuth.isAuthenticated);
  const { loading, error, draft } = useAppSelector(s => s.hypertensionStages);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [sbp, setSbp] = useState(140);
  const [dbp, setDbp] = useState(90);
  const [image, setImage] = useState<File | null>(null);
  const [video, setVideo] = useState<File | null>(null);
  const [note, setNote] = useState("");
  useEffect(() => {
    if (authenticated) void dispatch(fetchHypertensionDraft());
  }, [authenticated, dispatch]);
  const create = async (publish: boolean) => {
    if (!authenticated) return;
    setNote("");
    const form = new FormData();
    form.append("title", title.trim());
    form.append("description", description.trim());
    form.append("systolic_bp", String(sbp));
    form.append("diastolic_bp", String(dbp));
    if (image) form.append("image", image);
    if (video) form.append("video", video);
    const result = await dispatch(createHypertensionStage(form));
    if (!createHypertensionStage.fulfilled.match(result)) return;
    if (!publish) { setNote(`Черновик №${result.payload.id} сохранён.`); return; }
    const published = await dispatch(publishHypertensionStage(result.payload.id));
    if (!publishHypertensionStage.fulfilled.match(published)) return;
    void dispatch(fetchHypertensionStages());
    navigate("/stages");
  };
  const publishDraft = async () => {
    if (!draft || draft.status !== "Draft") return;
    const result = await dispatch(publishHypertensionStage(draft.id));
    if (publishHypertensionStage.fulfilled.match(result)) {
      void dispatch(fetchHypertensionStages());
      navigate("/stages");
    }
  };
  if (!authenticated) return <main className="form-page"><h1>Добавление стадии</h1>
    <p>Создавать и публиковать услуги могут только авторизованные пользователи.</p>
    <Link className="publish-button add-login-link" to="/login">Войти</Link>
  </main>;
  return <main className="draft-shell"><div className="draft-content">
    <h1>Новая стадия гипертонии</h1>
    {draft && draft.status === "Draft" && <div className="existing-draft">
      <strong>Ваш черновик №{draft.id}: {draft.title}</strong>
      <button type="button" className="publish-button" disabled={loading} onClick={() => { void publishDraft(); }}>Опубликовать черновик</button>
    </div>}
    <form onSubmit={e => { e.preventDefault(); void create(false); }}>
      <div className="draft-field"><label htmlFor="new_title">Название</label>
        <input id="new_title" type="text" value={title} onChange={e => setTitle(e.target.value)} disabled={loading} required maxLength={160}/></div>
      <div className="draft-field text-field"><label htmlFor="new_description">Описание</label>
        <textarea id="new_description" value={description} onChange={e => setDescription(e.target.value)} disabled={loading} required rows={4}/></div>
      <div className="draft-field pressure-field"><label htmlFor="new_sbp">Систолическое АД</label>
        <div className="value-row"><input id="new_sbp" type="number" min={40} max={300} value={sbp} onChange={e => setSbp(Number(e.target.value))} disabled={loading} required/>
          <span>мм рт. ст.</span></div></div>
      <div className="draft-field pressure-field"><label htmlFor="new_dbp">Диастолическое АД</label>
        <div className="value-row"><input id="new_dbp" type="number" min={20} max={250} value={dbp} onChange={e => setDbp(Number(e.target.value))} disabled={loading} required/>
          <span>мм рт. ст.</span></div></div>
      <div className="draft-field"><label htmlFor="new_image">Изображение</label>
        <input id="new_image" type="file" accept="image/*" disabled={loading} onChange={e => setImage(e.target.files?.[0] || null)}/>
        <small>Загрузится через Go API в MinIO.</small></div>
      <div className="draft-field"><label htmlFor="new_video">Видео</label>
        <input id="new_video" type="file" accept="video/*" disabled={loading} onChange={e => setVideo(e.target.files?.[0] || null)}/></div>
      {error && <p className="api-warning" role="alert">{error}</p>}
      {note && <p className="api-connected" role="status">{note}</p>}
      <button type="submit" className="publish-button" disabled={loading}>Сохранить черновик</button>
      <button type="button" className="publish-button publish-secondary" disabled={loading || !title.trim() || !description.trim()}
        onClick={() => { void create(true); }}>Создать и опубликовать</button>
    </form>
  </div></main>;
}
