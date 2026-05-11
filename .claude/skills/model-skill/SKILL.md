---
name: model-skill
description: "试衣模型用途与计费统一规范。用于处理 tryon/refiner/parsing/beautify 的配置绑定、调用链、扣退币口径、统计归因与 Uni/Web 展示一致性。"
---

# Model Skill（试衣模型规范）

## 适用场景
- 调整试衣模型配置结构（tryon_models）
- 修复“模型用途错配”导致的调用/计费/统计异常
- 修复“额外消耗与单次消耗混用”
- 修复 Uni/Web 模型消耗展示与后端不一致

## 模型用途定义（强约束）
1. `tryon`：主生成模型，负责基础试衣。
2. `refiner`：图片精修协同模型，生成后追加执行。
3. `parsing`：分割协同模型，单上装/单下装场景可开启。
4. `beautify`：美肤协同模型（可选）。

## 配置与绑定规则（tryon_models）
1. 主模型通过 `refinerModelKey`、`parsingModelKey`、`beautifyModelKey` 绑定协同模型。
2. 协同模型参数（url/token/provider/model 等）必须来自自身配置，不跨模型借用。
3. 允许主模型声明 `supportsRefiner/supportsParsing`；如有 `*ModelKey`，以绑定关系优先。

## 计费口径（必须一致）
1. `tryon` 基础消耗：仅取主模型 `cost`（或 `costPoints`）。
2. `refiner` 额外消耗：仅取 `refinerExtraCost` 或 `refinerExtraPoints`。
3. `parsing` 额外消耗：仅取 `parsingExtraCost` 或 `parsingExtraPoints`。
4. **禁止回退**：`refiner/parsing` 在 extraCost 为 `0` 时，仍然按 `0` 计算，不能回退到 `cost`。
5. 退款必须按 `operation_type` 拆分（`tryon/refiner/parsing/beautify`），避免混账。

## 后端实现要点
1. `tryonModelConfig.refinerExtraCostValue()` 只读 `RefinerExtraCost/RefinerExtraPoints`。
2. `tryonModelConfig.parsingCostValue()` 只读 `ParsingExtraCost/ParsingExtraPoints`。
3. 主调用与协同调用的 `operation_type` 必须准确写入扣费记录。
4. 下装图片链路必须统一 URL 前缀处理，避免 `modelDownloadRefuse`（`templateImageLower` 同样处理）。

## 前端实现要点（Uni）
1. 模型解析需将 `parsingModelKey -> parsingExtraCost` 映射到主模型。
2. 总消耗计算：`base + (refinerEnabled ? refinerExtra : 0) + (parsingEnabled ? parsingExtra : 0)`。
3. 界面必须显示协同额外消耗，`0` 也要展示。
4. Uni 所有用户可见文案必须走 i18n（含 toast/hint/按钮/说明）。

## 统计口径
1. 统计维度至少包含：`modelUsage`、`provider`、`model`。
2. `tryon` 基础消耗回填（历史数据兼容）建议：
- `base = GREATEST(cost_points - refund_points - refiner_cost_points - parsing_cost_points, 0)`。
3. 禁止把协同消耗归入 `tryon` 基础统计。

## 最小回归清单
1. 同时上传上装+下装：任务请求与模型调用都能收到双图。
2. 单上装/单下装 + 开启 parsing：正确扣 `parsingExtraCost`，`0` 不回退。
3. 开启 refiner：正确扣 `refinerExtraCost`，`0` 不回退。
4. 失败退款：按 `operation_type` 分项回退正确。
5. Uni 试衣间：总价与明细一致，额外消耗 `0` 可见。
6. 管理端统计：用途分组、供应商分组、消耗分组一致。

## 推荐自检命令
```bash
# 后端
cd server
gofmt -w service/client/tryon_task.go
go test ./service/client/...

# Uni
cd ../uni
npm run build:h5
```
