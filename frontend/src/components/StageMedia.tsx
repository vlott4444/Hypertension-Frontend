import { useEffect, useState } from "react";
import type { HypertensionStage } from "../types";

export const DEFAULT_IMAGE = "/default-stage.png";

/** Convert backend MinIO URLs to the second same-origin Vite proxy route. */
export function mediaUrl(stage: HypertensionStage, kind: "image" | "video"): string {
  const remote = kind === "image" ? stage.image_url : stage.video_url;
  const name = kind === "image" ? stage.image_name : stage.video_name;
  if (remote) {
    try {
      const u = new URL(remote, window.location.href);
      if (u.protocol === "http:" || u.protocol === "https:") {
        if (u.port === "9000" || /(?:^|[.-])minio(?:$|[.-])|(?:^|[.-])silo(?:$|[.-])/i.test(u.hostname)) {
          return `/s3${u.pathname}${u.search}`;
        }
        // Non-MinIO URLs should remain unchanged.
        return remote;
      }
    } catch { /* Handle relative media paths below. */ }
    if (remote.startsWith("/s3/")) return remote;
    if (remote.startsWith("/hypertension-media/")) return `/s3${remote}`;
    if (remote.startsWith("/")) return remote;
  }
  if (!name) return "";
  if (/^https?:\/\//i.test(name)) return mediaUrl({ ...stage, image_url: kind === "image" ? name : "", video_url: kind === "video" ? name : "" }, kind);
  return stage.is_mock ? `/media/${encodeURIComponent(name)}` : `/s3/hypertension-media/${encodeURIComponent(name)}`;
}
const localUrl = (stage: HypertensionStage, kind: "image" | "video") => {
  const name = kind === "image" ? stage.image_name : stage.video_name;
  return name ? `/media/${encodeURIComponent(name)}` : "";
};

export function StageImage({ stage, className = "" }: { stage: HypertensionStage; className?: string }) {
  const firstUrl = mediaUrl(stage, "image");
  const [src, setSrc] = useState(firstUrl || DEFAULT_IMAGE);
  useEffect(() => { setSrc(firstUrl || DEFAULT_IMAGE); }, [firstUrl]);
  const fallback = () => {
    const local = localUrl(stage, "image");
    setSrc(current => current !== local && local ? local : DEFAULT_IMAGE);
  };
  return <img className={className} src={src} alt={stage.title} onError={fallback}/>;
}

export function StageVideo({ stage }: { stage: HypertensionStage }) {
  const firstUrl = mediaUrl(stage, "video");
  const [src, setSrc] = useState(firstUrl);
  useEffect(() => { setSrc(firstUrl); }, [firstUrl]);
  const onError = () => {
    const local = localUrl(stage, "video");
    setSrc(current => current !== local && local ? local : "");
  };
  if (!src) return <div className="video-fallback"><StageImage stage={stage}/><span>Видео отсутствует — показано изображение</span></div>;
  return <video key={`${stage.id}:${src}`} className="feed-video" src={src} autoPlay muted loop playsInline controls
    onError={onError} poster={mediaUrl(stage, "image") || DEFAULT_IMAGE}/>;
}
