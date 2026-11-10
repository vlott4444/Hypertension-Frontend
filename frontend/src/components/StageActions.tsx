import type { HypertensionStage } from "../types";
import { deleteHypertensionStage, toggleHypertensionLike } from "../store/stagesSlice";
import { useAppDispatch, useAppSelector } from "../store/store";

export function StageActions({ stage, onDeleted }: { stage: HypertensionStage; onDeleted?: () => void }) {
  const dispatch = useAppDispatch();
  const isAuthenticated = useAppSelector(s => s.hypertensionAuth.isAuthenticated);
  const loading = useAppSelector(s => s.hypertensionStages.loading);
  const canLike = isAuthenticated && !stage.is_mock;
  const canDelete = canLike && stage.is_owner === true;
  return <div className="stage-actions">
    <button type="button" className="like-button" disabled={!canLike || loading}
      title={canLike ? "Поставить или убрать лайк" : "Чтобы лайкнуть, войдите"}
      aria-label={stage.is_liked ? "Убрать лайк" : "Поставить лайк"}
      onClick={() => { void dispatch(toggleHypertensionLike({ id: stage.id, wasLiked: !!stage.is_liked })); }}>
      {stage.is_liked ? "♥" : "♡"} {stage.likes_count ?? 0}
    </button>
    {canDelete && <button type="button" className="stage-delete-button" disabled={loading}
      onClick={async () => {
        if (!window.confirm("Удалить свою карточку?")) return;
        const result = await dispatch(deleteHypertensionStage(stage.id));
        if (deleteHypertensionStage.fulfilled.match(result)) onDeleted?.();
      }}>Удалить</button>}
  </div>;
}
