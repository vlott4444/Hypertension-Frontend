import type { HypertensionStage, StageFilters, StageStatus } from "../types";

const record = (v: unknown): Record<string, unknown> =>
  v && typeof v === "object" && !Array.isArray(v) ? v as Record<string, unknown> : {};

const first = (o: Record<string, unknown>, ...keys: string[]): unknown => {
  for (const k of keys) if (o[k] !== undefined && o[k] !== null) return o[k];
  return undefined;
};
const str = (v: unknown): string => typeof v === "string" ? v : "";
const num = (v: unknown, fallback = 0): number => {
  const parsed = Number(v);
  return Number.isFinite(parsed) ? parsed : fallback;
};
const bool = (v: unknown): boolean => v === true || v === 1 || v === "true" || v === "1";

export function normalizeStage(raw: unknown): HypertensionStage {
  const o = record(raw);
  const statusValue = str(first(o, "status", "Status")).toLowerCase();
  const status: StageStatus = statusValue === "draft" || statusValue === "черновик" ? "Draft" :
    statusValue === "deleted" || statusValue === "удален" ? "Deleted" : "Published";
  return {
    id: num(first(o, "id", "ID", "stage_id")),
    title: str(first(o, "title", "Title", "name")),
    description: str(first(o, "description", "Description")),
    systolic_bp: num(first(o, "systolic_bp", "SystolicBP", "systolicBP", "sbp")),
    diastolic_bp: num(first(o, "diastolic_bp", "DiastolicBP", "diastolicBP", "dbp")),
    status,
    image_name: str(first(o, "image_name", "ImageName", "imageName")),
    video_name: str(first(o, "video_name", "VideoName", "videoName")),
    image_url: str(first(o, "image_url", "ImageURL", "imageUrl")),
    video_url: str(first(o, "video_url", "VideoURL", "videoUrl")),
    is_owner: bool(first(o, "is_owner", "IsOwner", "isOwner")),
    is_liked: bool(first(o, "is_liked", "IsLiked", "isLiked")),
    likes_count: num(first(o, "likes_count", "LikesCount", "likesCount", "like_count", "likes")),
  };
}

export function unpackList(response: unknown): HypertensionStage[] {
  const o = record(response);
  const data = first(o, "stages", "items", "services", "hypertension_stages", "data", "results");
  const items = Array.isArray(response) ? response : Array.isArray(data) ? data : null;
  if (!items) throw new Error("Сервер вернул неизвестный формат списка услуг.");
  return items.map(normalizeStage).filter(item => item.id > 0);
}

export function unpackStage(response: unknown): HypertensionStage {
  const o = record(response);
  const nested = first(o, "stage", "item", "service", "hypertension_stage", "data", "draft");
  const result = normalizeStage(nested && typeof nested === "object" ? nested : response);
  if (!result.id) throw new Error("Сервер вернул карточку без id.");
  return result;
}

export function applyFilters(list: HypertensionStage[], f: StageFilters): HypertensionStage[] {
  const title = f.title.trim().toLowerCase();
  return list.filter(s => s.status === "Published" && s.systolic_bp >= f.minBP &&
    s.systolic_bp <= f.maxBP && (!title || s.title.toLowerCase().includes(title)));
}
