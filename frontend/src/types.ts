export type StageStatus = "Published" | "Draft" | "Deleted";

export interface HypertensionStage {
  id: number;
  title: string;
  description: string;
  systolic_bp: number;
  diastolic_bp: number;
  status: StageStatus;
  image_name: string;
  video_name: string;
  image_url?: string;
  video_url?: string;
  is_owner?: boolean;
  is_liked?: boolean;
  likes_count?: number;
  /** Explicitly marks local examples; backend items never acquire this flag. */
  is_mock?: boolean;
}

export interface StageFilters {
  minBP: number;
  maxBP: number;
  title: string;
}
