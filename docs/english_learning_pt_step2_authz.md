# PT完善 Step2 统一授权引擎与用户粒度授权

更新时间: 2026-07-01

## 1. 本步目标
- 解决 R09: 权限口径统一（免费期/免费分钟/VIP/资源授权）。
- 解决 R10: 新用户学习资产创建时自动发放全站免费期。
- 解决 R11: 支持用户粒度授权（分类/剧集/单集）。

## 2. 本步交付
- 新增用户资源授权模型：`UserLearningEntitlement`。
- 新增统一授权服务：`LearningAuthzService`。
- 中间件改造：`LearningAuth` 改为统一服务判定并注入用户ID上下文。
- 视频单集详情权限改造：`findVideoEpisode` 按资源粒度返回 `hasFullAuth`。
- 新增用户授权管理 API：
  - `POST /englishLearning/asset/grantEntitlement`
  - `DELETE /englishLearning/asset/revokeEntitlement`
  - `GET /englishLearning/asset/getEntitlementList`
- 新增系统配置项：`learning_new_user_free_hours`（默认 24 小时）。

## 3. 授权优先级（统一口径）
- 全局授权（任一满足即全量）：
  - 免费期内（`freeTimeExpire > now`）
  - 免费分钟余额 > 0
  - VIP = true
- 资源授权（在全局授权不满足时）：
  - 单集授权命中
  - 或单集所属剧集授权命中
- 资源收费判定：
  - 若资源免费（`price=0` 且 `needVip=false`），直接全量。
  - 若资源受限且无命中授权，返回试看态（`hasFullAuth=false`）。

## 4. 新增与改造文件
- 模型与请求：
  - `server/plugin/english_learning/model/auth_model.go`
  - `server/plugin/english_learning/model/request/authz.go`
- 服务层：
  - `server/plugin/english_learning/service/learning_authz.go`
  - `server/plugin/english_learning/service/enter.go`
  - `server/plugin/english_learning/service/user_learning_asset.go`
  - `server/plugin/english_learning/service/checkin.go`
  - `server/plugin/english_learning/service/user_data.go`
- API/中间件/路由：
  - `server/plugin/english_learning/middleware/learning_auth.go`
  - `server/plugin/english_learning/api/content.go`
  - `server/plugin/english_learning/api/user_learning_asset.go`
  - `server/plugin/english_learning/router/user_learning_asset.go`
- 初始化与配置：
  - `server/plugin/english_learning/initialize/gorm.go`
  - `server/initialize/other.go`

## 5. 验收结果
- 自动迁移覆盖：`UserLearningEntitlement` 已纳入插件 `AutoMigrate`。
- 新用户资产默认免费期：`EnsureUserAsset` 创建资产时自动写入 `freeTimeExpire`。
- 用户粒度授权可运维：授权/撤销/查询 API 已打通。
- 播放授权可落地：`findVideoEpisode` 返回 `hasFullAuth` 时已包含资源粒度判定。

## 6. 测试结果
- 已通过：`cd server && go test ./plugin/english_learning/...`
- 新增测试：
  - `server/plugin/english_learning/service/learning_authz_service_test.go`
  - `server/plugin/english_learning/api/user_learning_asset_entitlement_api_test.go`

## 7. 结论
- PT-2 已完成，可进入 PT-3（Web 运营端全链路 CRUD 补齐）。
