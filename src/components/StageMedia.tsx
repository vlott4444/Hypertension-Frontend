import { useState } from "react";
import type { HypertensionStage } from "../types";
export const DEFAULT_IMAGE = "/default-stage.png";
export function StageImage({stage,className=""}:{stage:HypertensionStage;className?:string}){
 const [failed,setFailed]=useState(false);
 const src=(!stage.image_name || failed)?DEFAULT_IMAGE:`/media/${stage.image_name}`;
 return <img className={className} src={src} alt={stage.title} onError={()=>setFailed(true)}/>;
}
export function StageVideo({stage}:{stage:HypertensionStage}){
 const [failed,setFailed]=useState(false);
 if(!stage.video_name || failed)return <div className="video-fallback"><StageImage stage={stage}/><span>Видео отсутствует — показано изображение по умолчанию</span></div>;
 return <video key={stage.id} className="feed-video" src={`/media/${stage.video_name}`} autoPlay muted loop playsInline controls onError={()=>setFailed(true)} poster={stage.image_name?`/media/${stage.image_name}`:DEFAULT_IMAGE}/>;
}
