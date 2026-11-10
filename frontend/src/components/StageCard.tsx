import { Link } from "react-router-dom";
import { StageImage } from "./StageMedia";
import { StageActions } from "./StageActions";
import type { HypertensionStage } from "../types";
interface Props { stage: HypertensionStage }
const shortTitle = (title: string) => title.replace(" стадия гипертонии", " стадия");
export function StageCard({ stage }: Props) {
  return <article className="stage-card">
    <Link className="stage-card-link" to={`/stages/${stage.id}`}>
      <StageImage stage={stage}/>
      <div className="stage-card-copy"><h2>{shortTitle(stage.title)}</h2>
        <p>САД {stage.systolic_bp} / ДАД {stage.diastolic_bp}</p></div>
    </Link>
    <StageActions stage={stage}/>
  </article>;
}
