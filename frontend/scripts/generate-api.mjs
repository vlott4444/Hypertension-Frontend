/** OpenAPI -> TypeScript Axios service generator (no generator-only dependencies).
 * `npm run generate-api` generates from documented Lab4 route snapshot;
 * `npm run generate-api:server` reads the actual Go Swagger on localhost:8080.
 * Method names are intentionally stable to keep thunk call-sites readable.
 */
import { readFile, writeFile, mkdir } from "node:fs/promises";
import { resolve } from "node:path";

const useServer = process.argv.includes("--from-server");
const sourceUrl = process.env.SWAGGER_URL || "http://localhost:8080/swagger/openapi.json";
const spec = useServer ? await (async () => {
  const response = await fetch(sourceUrl);
  if (!response.ok) throw new Error(`Swagger ${sourceUrl}: ${response.status}`);
  return response.json();
})() : JSON.parse(await readFile(resolve("openapi/lab4.contract.json"), "utf8"));

const paths = spec.paths || {};
const route = (method, matcher) => {
  const found = Object.entries(paths).find(([url, methods]) => matcher(url) && methods?.[method]);
  if (!found) throw new Error(`No ${method.toUpperCase()} route for ${matcher.toString()} in Swagger`);
  return found[0].replace(/^\/api(?=\/|$)/, "");
};
const get = pattern => route("get", url => pattern.test(url));
const post = pattern => route("post", url => pattern.test(url));
const output = `/* AUTO-GENERATED from ${useServer ? sourceUrl : "openapi/lab4.contract.json"}.
 * Regenerate via npm run generate-api${useServer ? ":server" : ""}. DO NOT EDIT BY HAND.
 * Paths come directly from OpenAPI; all methods use shared Axios instance.
 */
import type { AxiosResponse } from "axios";
import { http } from "../http";

export type HypertensionStagesQuery = { min_sbp?: number; max_sbp?: number; title?: string };
export type LikeValue = { value: 0 | 1 };
export type CreateHypertensionStageForm = FormData;

export class HypertensionStagesApi {
  listHypertensionStages(query: HypertensionStagesQuery = {}): Promise<AxiosResponse<unknown>> {
    return http.get(${JSON.stringify(get(/^\/api\/hypertension_stages\/?$/))}, { params: query });
  }
  getHypertensionFeed(): Promise<AxiosResponse<unknown>> {
    return http.get(${JSON.stringify(get(/^\/api\/feed\/?$/))});
  }
  getHypertensionFeedById(id: number, next = false): Promise<AxiosResponse<unknown>> {
    return http.get(\`${get(/^\/api\/feed\/\{[^}]+\}$/).replace(/\{[^}]+\}/, "${id}")}\`, { params: { next } });
  }
  getHypertensionDraft(): Promise<AxiosResponse<unknown>> {
    return http.get(${JSON.stringify(get(/^\/api\/draft\/?$/))});
  }
  createHypertensionStage(data: CreateHypertensionStageForm): Promise<AxiosResponse<unknown>> {
    return http.post(${JSON.stringify(post(/^\/api\/hypertension_stages\/?$/))}, data);
  }
  publishHypertensionStage(id: number): Promise<AxiosResponse<unknown>> {
    return http.put(\`${route("put", u => /^\/api\/hypertension_stages\/\{[^}]+\}\/publish$/.test(u)).replace(/\{[^}]+\}/, "${id}")}\`);
  }
  deleteHypertensionStage(id: number): Promise<AxiosResponse<unknown>> {
    return http.delete(\`${route("delete", u => /^\/api\/hypertension_stages\/\{[^}]+\}$/.test(u)).replace(/\{[^}]+\}/, "${id}")}\`);
  }
  likeHypertensionStage(id: number, data: LikeValue): Promise<AxiosResponse<unknown>> {
    return http.post(\`${post(/^\/api\/hypertension_stages\/\{[^}]+\}\/like$/).replace(/\{[^}]+\}/, "${id}")}\`, data);
  }
}

export const hypertensionStagesApi = new HypertensionStagesApi();
`;
await mkdir(resolve("src/api/generated"), { recursive: true });
await writeFile(resolve("src/api/generated/HypertensionStagesApi.ts"), output);
console.log(`Generated Axios service from ${useServer ? sourceUrl : "openapi/lab4.contract.json"}`);
