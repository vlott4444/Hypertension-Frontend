// Zero-dependency smoke check of generator and bundled API contract.
import { readFile } from "node:fs/promises";
const spec = JSON.parse(await readFile("openapi/lab4.contract.json", "utf8"));
const generated = await readFile("src/api/generated/HypertensionStagesApi.ts", "utf8");
const required = [
  "listHypertensionStages", "getHypertensionFeed", "getHypertensionFeedById",
  "getHypertensionDraft", "createHypertensionStage", "publishHypertensionStage",
  "deleteHypertensionStage", "likeHypertensionStage",
];
for (const name of required) if (!generated.includes(`${name}(`)) throw Error(`${name} missing from generated file`);
if (Object.keys(spec.paths).length !== 10) throw Error("Unexpected OpenAPI snapshot path count");
console.log(`Smoke check passed: ${required.length} Axios operations, ${Object.keys(spec.paths).length} OpenAPI paths.`);
