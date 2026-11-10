import { useEffect, useState } from "react";
import { StageVideo } from "../components/StageMedia";
import { loadHypertensionFeed } from "../store/feedSlice";
import { toggleHypertensionLike } from "../store/stagesSlice";
import { useAppDispatch, useAppSelector } from "../store/store";

function heading(title: string) {
  const match = title.match(/^(I|II|III) стадия/);
  return match ? <>{match[1]} стадия<br/>гипертонической<br/>болезни</> : title;
}
export function FeedPage() {
  const dispatch = useAppDispatch();
  const [expanded, setExpanded] = useState(false);
  const { current: stage, loading, error, source } = useAppSelector(s => s.hypertensionFeed);
  const isAuthenticated = useAppSelector(s => s.hypertensionAuth.isAuthenticated);
  useEffect(() => {
    document.title = "Лента — стадии гипертонии";
    void dispatch(loadHypertensionFeed({ next: false }));
    return () => { document.title = "Стадии гипертонии"; };
  }, [dispatch, isAuthenticated]);
  const next = () => { setExpanded(false); void dispatch(loadHypertensionFeed({ next: true })); };
  if (!stage) return <main className="feed-shell feed-empty">
    {error && <p className="api-warning" role="alert">{error}</p>}
    {loading && <p>Загрузка ленты...</p>}
  </main>;
  return <main className="feed-shell"><section className="feed-media">
    <StageVideo key={stage.id} stage={stage}/><div className="feed-shade"/>
    {source === "mock" && <div className="feed-demo-label">mock — сервер недоступен</div>}
    <input type="checkbox" className="desc-toggle" checked={expanded} onChange={e => setExpanded(e.target.checked)} id="desc-toggle"/>
    <aside className="feed-actions">
      <button type="button" className="feed-like" aria-label={stage.is_liked ? "Убрать лайк" : "Поставить лайк"}
        disabled={!isAuthenticated || !!stage.is_mock || loading}
        onClick={() => { void dispatch(toggleHypertensionLike({ id: stage.id, wasLiked: !!stage.is_liked })); }}>
        <span className="heart-icon">{stage.is_liked ? "♥" : "♡"}</span><span>{stage.likes_count || 0}</span>
      </button>
      <button type="button" className="next-button" disabled={loading} onClick={next}>Далее</button>
    </aside>
    <section className="feed-copy"><h1>{heading(stage.title)}</h1><div className="feed-divider"/>
      <div className="feed-description"><p className="desc-text">{stage.description}</p>
        <label htmlFor="desc-toggle" className="more-button"><span className="more-label-open">ещё</span>
          <span className="more-label-close">свернуть</span></label>
      </div>
      {error && <p className="api-warning" role="alert">{error}</p>}
    </section>
  </section></main>;
}
