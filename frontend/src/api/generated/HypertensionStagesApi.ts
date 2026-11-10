/* AUTO-GENERATED from openapi/lab4.contract.json.
 * Regenerate via npm run generate-api. DO NOT EDIT BY HAND.
 * Paths come directly from OpenAPI; all methods use shared Axios instance.
 */
import type { AxiosResponse } from "axios";
import { http } from "../http";

export type HypertensionStagesQuery = { min_sbp?: number; max_sbp?: number; title?: string };
export type LikeValue = { value: 0 | 1 };
export type CreateHypertensionStageForm = FormData;

export class HypertensionStagesApi {
  listHypertensionStages(query: HypertensionStagesQuery = {}): Promise<AxiosResponse<unknown>> {
    return http.get("/hypertension_stages", { params: query });
  }
  getHypertensionFeed(): Promise<AxiosResponse<unknown>> {
    return http.get("/feed");
  }
  getHypertensionFeedById(id: number, next = false): Promise<AxiosResponse<unknown>> {
    return http.get(`/feed/${id}`, { params: { next } });
  }
  getHypertensionDraft(): Promise<AxiosResponse<unknown>> {
    return http.get("/draft");
  }
  createHypertensionStage(data: CreateHypertensionStageForm): Promise<AxiosResponse<unknown>> {
    return http.post("/hypertension_stages", data);
  }
  publishHypertensionStage(id: number): Promise<AxiosResponse<unknown>> {
    return http.put(`/hypertension_stages/${id}/publish`);
  }
  deleteHypertensionStage(id: number): Promise<AxiosResponse<unknown>> {
    return http.delete(`/hypertension_stages/${id}`);
  }
  likeHypertensionStage(id: number, data: LikeValue): Promise<AxiosResponse<unknown>> {
    return http.post(`/hypertension_stages/${id}/like`, data);
  }
}

export const hypertensionStagesApi = new HypertensionStagesApi();
