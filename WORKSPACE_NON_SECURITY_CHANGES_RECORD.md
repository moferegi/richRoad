# 工作区非安全改动完整记录（补充）

更新时间：2026-05-20

适用范围：本工作区当前变更中，已确认不属于安全主线的内容。

说明：
- 本文是对 `WORKSPACE_SECURITY_IMPROVEMENTS_RECORD.md` 的补充。
- 当前已确认的非安全主线改动，集中在 geo 数据生成链路与 CI 工程化流程。

---

## 1. 变更总览（Geo 数据链路）

Geo 数据链路改动共 3 个文件：
- server/source/geo/generate_iso_geo_sql.mjs
- server/source/geo/geo_language_country_templates.sql
- server/source/geo/geo_zh_translate_cache.json

文件体积（当前工作区）：
- generate_iso_geo_sql.mjs: 31,677 bytes
- geo_language_country_templates.sql: 28,122,185 bytes
- geo_zh_translate_cache.json: 1,200,708 bytes

---

## 2. 文件逐条记录

### 2.1 server/source/geo/generate_iso_geo_sql.mjs

改了什么：
- 新增地理数据 SQL 生成脚本（Node.js）。
- 支持从多数据源组装国家/省市区层级数据，并统一产出多语言字段。
- 支持翻译缓存文件读写（geo_zh_translate_cache.json）。

有什么用：
- 让 geo 数据更新从“手工维护 SQL”变为“脚本生成”。
- 统一不同语言字段的回退策略，减少手工补漏。

能防止什么问题：
- 手工改 SQL 造成格式不一致、字段漏填、批量变更难回溯。

小白怎么操作：
1. 在受控环境执行脚本，不要直接在生产机临时跑。
2. 先检查生成结果再导入数据库。
3. 提交前确认 geo_zh_translate_cache.json 仅包含可公开地名映射。

---

### 2.2 server/source/geo/geo_language_country_templates.sql

改了什么：
- 新增大体量 geo 模板 SQL（含国家/省市区多级数据）。
- SQL 头部说明明确了数据来源和规则：
  - 基于 GeoNames ADM1-ADM5 + 真实地名兜底
  - level3-5 使用真实地名
  - name_i18n 缺失按英文回退

有什么用：
- 可一次性初始化或重建 geo 相关数据。
- 方便新环境对齐同一套基础地理库。

能防止什么问题：
- 环境间 geo 基础数据不一致导致的筛选、展示、统计偏差。

小白怎么操作：
1. 导入前先备份现有 geo_areas / geo_cities / geo_provinces。
2. 在测试库先导入验证，再推进到生产。
3. 关注导入耗时与锁表影响，尽量在低峰操作。

---

### 2.3 server/source/geo/geo_zh_translate_cache.json

改了什么：
- 新增地名翻译缓存映射（英文地名 -> 中文）。
- 为脚本批量生成 SQL 提供缓存命中，减少重复翻译请求。

有什么用：
- 提升脚本生成速度与稳定性。
- 同一地名在不同批次生成时保持翻译一致。

能防止什么问题：
- 每次重跑脚本都重新翻译，导致结果漂移或速率受限。

小白怎么操作：
1. 把它当“构建缓存”管理，不要手工大范围改写。
2. 如需修正少量翻译，按 key 精确修改并复跑校验。
3. 避免混入乱码或错误编码，提交前抽样检查 JSON 可解析。

---

## 3. 关联关系与边界说明

- 以上 3 个文件是同一条数据链路：
  - `generate_iso_geo_sql.mjs` 负责生成
  - `geo_zh_translate_cache.json` 负责翻译缓存
  - `geo_language_country_templates.sql` 负责最终导入模板

- 这条链路属于“基础数据构建与多语言地理数据维护”，不属于权限、鉴权、限流、上传、WS、导出等安全主线。

---

## 4. 上线与回滚建议（非安全）

上线建议：
1. 在测试环境完整跑通“生成 -> 导入 -> 页面验证”。
2. 验证多语言展示（尤其 zh / en / zh-TW）是否符合预期。
3. 再进行生产导入窗口安排。

回滚建议：
1. 回滚到导入前备份的 geo 三表数据。
2. 同时回退本次 SQL 与缓存文件版本。
3. 重新验证 geo 级联选择与检索结果。

---

## 5. 与安全文档的关系

- 安全主线记录见：WORKSPACE_SECURITY_IMPROVEMENTS_RECORD.md
- 本文仅记录非安全主线改动，不覆盖安全策略与开关。

---

## 6. 2026-05-20 CI 工程化改进（非安全）

### 6.1 .github/workflows/ci.yaml

改了什么：
- 将 backend/devops-test/devops-prod 的 Go 版本从 `1.22` 统一提升到 `1.24.2`，与项目 `server/go.mod` 对齐。
- backend 任务新增测试门禁：在 `go build -race` 前执行
  1. `go test ./... -run TestDoesNotExist -count=1`
  2. `go test ./service/shop ./service/client ./task ./api/v1/client -count=1`
- frontend 任务在 web build 后新增 uni H5 build（`npm run build:h5`）。

有什么用：
- 降低“本地通过但 CI 因 Go 版本偏差失败”的概率。
- 让后端 CI 从“只构建”升级为“构建 + 测试”，回归问题更早暴露。
- 将 uni 纳入默认 CI 构建链路，避免仅 web 可发布而 uni 产物潜在失配。

能防止什么问题：
- Go 版本不一致造成的语法/依赖兼容问题延后暴露。
- 安全与业务改动在 PR 阶段缺少最小测试闸门。
- Uni 端在无人值守构建中长期失去验证，发布时才发现问题。

小白怎么操作：
1. 提交 PR 后观察 CI 三个关键任务：frontend、backend、release/devops 相关链路。
2. 若 backend 失败，先本地执行与 CI 相同的两条 `go test` 命令复现。
3. 若 uni build 失败，先在 `uni` 目录执行 `npm install && npm run build:h5` 修复再推送。

### 6.2 前端 lint 基础能力补齐（渐进式，不阻断发布）

改了什么：
- `web/package.json` 新增 `lint` 与 `lint:fix` 脚本。
- `uni/package.json` 新增 `lint` 与 `lint:fix` 脚本，并补充 eslint 相关 devDependencies。
- 新增 `uni/eslint.config.mjs`，对 `src` 目录启用 Vue 基础规则，并默认忽略 `src/uni_modules` 等第三方目录。
- `frontend` CI 任务新增 web/uni lint 步骤，采用 advisory 模式（`continue-on-error: true`），不阻断现有构建发布。
- 修复本轮改动文件中的 lint 明确问题：
  - `web/src/components/exportExcel/exportExcel.vue`
  - `web/src/components/exportExcel/exportTemplate.vue`
  两处 `content-disposition` 文件名解析正则去除无效转义。

有什么用：
- 先把 lint 能力接入流水线并持续暴露问题，再分批治理历史存量，避免一次性切换导致 CI 全面红灯。
- 让 web 与 uni 两端都具备统一入口（`npm run lint`），便于团队逐步收敛代码风格。

能防止什么问题：
- 没有 lint 入口导致低级语法/规则问题长期积累。
- 前端质量问题只在人工验收阶段暴露，回归成本高。

小白怎么操作：
1. 本地先执行：`cd web && npm run lint`，修复自己改动相关报错。
2. 再执行：`cd uni && npm run lint`，优先处理业务目录 `uni/src/pages` 与 `uni/src/components`。
3. CI 里看到 advisory lint 报错时，不会阻断构建，但应在后续提交中逐步清理。

当前基线（2026-05-20，本地实测）：
- web：`157 errors`（主要集中在历史 `no-unused-vars`、`no-empty`、`no-mixed-spaces-and-tabs`、`vue/no-parsing-error`）。
- uni：`130 errors`（主要集中在历史 `no-unused-vars`、`no-empty`、`no-mixed-spaces-and-tabs`、局部语法错误）。

说明：
- 本轮已修复新增改动中的明确 lint 问题（exportExcel/exportTemplate 两个组件正则转义问题已清零）。
- 历史存量问题采用“先纳入 CI 可见、后分批治理”的策略，避免一次性切换导致发布链路中断。

### 6.3 前端 lint 存量分批治理进展（web）

改了什么：
- 按目录分批清理 web 历史 lint 存量，先后完成：
  1. `src/view/client` 全目录问题清零。
  2. `src/view/shop` 全目录问题清零（含批量未使用导入、模板语法、空块处理）。
  3. `src/components`、`src/view/systemTools`、`src/view/superAdmin`、`src/plugin/geo` 等剩余高频项清理。
- 对 `src/components/multilingual/multi-lang-editor.vue` 中“按设计原位编辑传入多语言对象”的场景，补充了 eslint 规则豁免注释，避免与现有调用方式冲突。

有什么用：
- 将 web lint 从“可见但高噪音”推进到“可直接作为质量门禁”的状态。
- 为后续把 CI 中 web lint 从 advisory 升级为 required 提供基础。

能防止什么问题：
- 历史未使用变量、空块、模板语法问题持续回流。
- 团队改动混入低级错误而未被及时发现。

本轮结果（2026-05-20，本地实测）：
- web：`157 errors -> 123 errors -> 41 errors -> 0 errors`。
- 命令：`cd web && npm run lint`，退出码 `0`。

说明：
- 目前前端 lint 存量主要剩余在 uni 侧（web 已清零）。
- 下一阶段建议按同样策略推进 uni 分批治理，再考虑将 CI lint 闸门从 advisory 调整为 required。

### 6.4 前端 lint 存量分批治理进展（uni）

改了什么：
- 对 uni 历史 lint 存量按“先降噪、后清零”方式处理：
  1. `uni/eslint.config.mjs` 新增 `src/js_sdk/**` 忽略，避免将第三方/SDK 历史脚本噪音纳入业务治理批次。
  2. 修复 `src/App.vue`、`src/common/*`、`src/components/*`、`src/pages/*`、`src/pinia/*`、`src/utils/*` 中的高频问题：`no-unused-vars`、`no-empty`、`no-mixed-spaces-and-tabs`、`no-unreachable`、`vue/valid-v-for`、`vue/require-valid-default-prop`、局部 parsing error。
  3. 对 tryon 页面 H5 分支中“条件编译 return”导致的 `no-unreachable`，改为运行时环境判断，保持行为一致并通过静态检查。

有什么用：
- 将 uni lint 从“高噪音告警”推进到“可作为门禁”的可执行状态。
- 降低后续 uni 业务迭代中回归到历史低级错误的风险。

能防止什么问题：
- 旧代码中的空 catch、混合缩进、未使用变量持续积累。
- 条件编译语句在静态分析阶段引入误报并掩盖真实问题。

本轮结果（2026-05-20，本地实测）：
- uni：`130 errors -> 41 errors -> 3 errors -> 0 errors`。
- 命令：`cd uni && npm run lint`，退出码 `0`。
- 回归：`cd uni && npm run build:h5` 构建完成（存在 Sass `@import` 废弃与 chunk 提示，非阻断）。

说明：
- 至此 web/uni 两端 lint 均已清零。
- 下一阶段可评估将 CI 中前端 lint 从 advisory 升级为 required（建议先 web+uni 同步启用，或分阶段启用并观察 1-2 个迭代周期）。

### 6.5 前端 CI lint 门禁升级（required）

改了什么：
- 在 `.github/workflows/ci.yaml` 中移除前端 lint 步骤的 `continue-on-error: true`：
  1. `Lint web (advisory)` 调整为 `Lint web`。
  2. `Lint uni (advisory)` 调整为 `Lint uni`。

有什么用：
- 前端 lint 从“可见但不阻断”升级为“质量门禁”，PR 阶段可直接阻断低级问题入主干。

能防止什么问题：
- lint 回归问题在构建阶段被忽略，导致问题延后到联调或上线阶段。

本轮结果（2026-05-20，本地实测）：
- `cd web && npm run lint`：退出码 `0`。
- `cd uni && npm run lint`：退出码 `0`。

说明：
- 在当前 web/uni 均已清零的前提下切换 required，风险可控。

### 6.6 前端 CI 安装与缓存优化

改了什么：
- 将 frontend 任务中的 Node 安装步骤由 `actions/setup-node@v1` 升级为 `actions/setup-node@v4`。
- 启用 npm 缓存：`cache: npm`，并配置 `cache-dependency-path` 覆盖：
  1. `web/package-lock.json`
  2. `uni/package-lock.json`
- 将 web/uni lint 步骤中的依赖安装从 `npm install` 调整为 `npm ci --no-audit --no-fund`。

有什么用：
- 提升 CI 安装阶段稳定性与可重复性（锁文件严格安装）。
- 通过缓存 npm 依赖下载结果缩短重复构建耗时。

能防止什么问题：
- `npm install` 在锁文件漂移或次要版本变化时引发的构建不一致。
- 前端双工程（web/uni）每次全量重新拉取依赖导致的流水线时长抖动。

说明：
- 本优化不改变构建产物逻辑，仅优化安装与缓存策略。

### 6.7 devops 发布链路 Node 安装策略统一

改了什么：
- 在 `.github/workflows/ci.yaml` 的 `devops-test` 与 `devops-prod` 两个 job 中：
  1. `actions/setup-node@v2.1.2` 升级为 `actions/setup-node@v4`。
  2. 增加 npm 缓存配置：`cache: npm` 与 `cache-dependency-path: web/package-lock.json`。
  3. `Build-Node` 从 `yarn install && yarn run build` 改为 `npm ci --no-audit --no-fund && npm run build`。

有什么用：
- 让“前端检查链路(frontend job)”与“发布链路(devops jobs)”保持一致的依赖安装策略。
- 提升发布构建的可复现性，并减少重复下载依赖造成的耗时。

能防止什么问题：
- 不同 job 使用不同包管理命令导致的安装差异与偶发构建不一致。
- 发布链路长时间重复安装依赖带来的时长波动。

说明：
- 本次未改动 docker job 内对 Makefile 的 `sed` 注入脚本（其内部仍包含 yarn 字样，属于镜像构建特定流程）。

### 6.8 前端 CI 并行化（web/uni 拆分）

改了什么：
- 在 `.github/workflows/ci.yaml` 中将原有单一 `frontend` job 拆分为两个并行执行 job：
  1. `frontend-web`：执行 web 的 lint + build。
  2. `frontend-uni`：执行 uni 的 lint + build:h5。
- 新增一个汇总 `frontend` job（保留原门禁名称语义），依赖 `frontend-web` 与 `frontend-uni`，仅用于汇总状态。
- `devops-test`、`release-please` 仍保持依赖 `frontend`，无需调整其工作流依赖写法。

有什么用：
- web 与 uni 前端检查可并行运行，缩短前端阶段总耗时。
- 保留 `frontend` 汇总 job 可降低 required check 名称迁移风险。

能防止什么问题：
- 单 job 串行执行导致的等待时间过长、反馈回路变慢。
- 直接改 required check 名称带来的分支保护配置漂移风险。

说明：
- 本次为 CI 编排优化，不涉及业务代码与构建产物逻辑变更。

### 6.9 后端与发布链路 Go 缓存化

改了什么：
- 在 `.github/workflows/ci.yaml` 中将 Go 运行时安装从 `actions/setup-go@v1` 升级为 `actions/setup-go@v5`，覆盖：
  1. `backend` job
  2. `devops-test` job
  3. `devops-prod` job
- 为上述 Go 步骤统一启用模块缓存：
  - `cache: true`
  - `cache-dependency-path: server/go.sum`
- `backend` job 的依赖预下载由历史 `go get -v -t -d ./...` 调整为：
  1. `go mod download`
  2. `go mod verify`

有什么用：
- 提升 Go 依赖安装阶段的可重复性与缓存命中率。
- 缩短 backend 与发布链路的模块下载时间。

能防止什么问题：
- 旧式 `go get` 在 CI 中可能引入不必要的依赖解析波动。
- 多条流水线重复下载相同 Go 模块导致时长抖动。

说明：
- 本次未调整发布构建中的 `go mod tidy` 习惯（仍保持原有行为），优先保证流程兼容。

### 6.10 CI 并发互斥与基础 Action 版本统一

改了什么：
- 在 `.github/workflows/ci.yaml` 顶层新增并发控制：
  - `concurrency.group: ci-${{ github.workflow }}-${{ github.ref }}`
  - `concurrency.cancel-in-progress: true`
- 将 workflow 内所有 `actions/checkout@v2` 统一升级为 `actions/checkout@v4`。

有什么用：
- 同一分支连续 push 时，自动取消旧流水线，减少排队与重复执行耗时。
- checkout action 版本统一，降低旧版本维护风险并提升一致性。

能防止什么问题：
- 多次快速提交触发的历史流水线堆积，导致反馈延迟。
- 不同 job 混用 checkout 旧版本造成的行为差异。

说明：
- 本次为 CI 流程治理，不涉及业务代码逻辑。

### 6.11 CI 作业超时上限治理

改了什么：
- 在 `.github/workflows/ci.yaml` 的关键 job 增加 `timeout-minutes` 上限：
  1. `init`: 5 分钟
  2. `frontend-web`: 25 分钟
  3. `frontend-uni`: 25 分钟
  4. `frontend`(汇总): 5 分钟
  5. `backend`: 35 分钟
  6. `devops-test`: 60 分钟
  7. `release-pr`: 20 分钟
  8. `release-please`: 20 分钟
  9. `devops-prod`: 60 分钟
  10. `docker`: 90 分钟

有什么用：
- 限制异常卡住任务的最长占用时长，缩短失败反馈时间。
- 结合并发取消机制，进一步降低 runner 资源浪费。

能防止什么问题：
- 外部网络抖动、远端部署命令阻塞导致作业长时间不退出。
- 无上限等待引起的流水线排队和反馈延迟。

说明：
- 本次仅调整 CI 运行策略，不涉及业务代码与构建逻辑。

### 6.12 发布链路远程命令重试机制

改了什么：
- 在 `.github/workflows/ci.yaml` 的 `devops-test` 与 `devops-prod` 两个 `restart` 步骤中新增轻量 `retry()` 函数（最多 3 次，线性退避 5s/10s）。
- 将以下远程网络敏感命令改为重试执行：
  1. 上传 `web/dist` 的 `scp`
  2. 上传 `web/ser` 的 `scp`
  3. 清理远端 `resource` 目录的 `ssh`
  4. 上传 `server/resource` 的 `scp`
  5. 触发远端 `restart.sh` 的 `ssh`

有什么用：
- 降低网络抖动导致的偶发部署失败，提升发布链路稳定性。
- 在短暂失败场景下自动自愈，减少人工重跑次数。

能防止什么问题：
- 单次 `scp/ssh` 因瞬时连接失败直接导致整条发布 job 失败。
- 发布阶段偶发失败带来的无效重试和等待成本。

说明：
- 本次未改变部署目标路径和重启逻辑，仅增强失败重试能力。

### 6.13 发布链路连接超时与阶段日志增强

改了什么：
- 在 `.github/workflows/ci.yaml` 的 `devops-test` 与 `devops-prod` 两个 `restart` 步骤中新增统一连接参数数组：
  - `SSH_OPTS=(-o StrictHostKeyChecking=no -o ConnectTimeout=10 -o ConnectionAttempts=2)`
- `ssh-keyscan` 增加超时参数：`-T 10`。
- 为远程部署关键阶段增加可读日志：
  1. upload web dist
  2. upload web binary
  3. clean remote resources
  4. upload server resources
  5. restart remote service

有什么用：
- 遇到网络不稳定时更快失败返回，避免单次连接长时间卡住。
- 通过阶段日志快速定位失败步骤，降低排障成本。

能防止什么问题：
- ssh/scp 默认连接等待过长导致的发布阶段“假卡死”。
- 失败日志信息不足，难以判断是上传、清理还是重启环节异常。

说明：
- 本次仅增强网络连接控制与日志可观测性，不改变发布业务逻辑。

### 6.14 发布链路重试逻辑脚本化复用

改了什么：
- 新增公共脚本 `.github/scripts/ci_retry.sh`，封装：
  1. `retry()`：支持最多重试次数与递增等待，失败时输出退出码。
  2. `run_with_retry()`：输出阶段名后执行带重试命令。
- `devops-test` 与 `devops-prod` 的 `restart` 步骤改为 `source .github/scripts/ci_retry.sh`，移除内联重复重试函数。
- 现在重试日志会包含“第几次重试/总次数、退出码、等待秒数、命令内容”。

有什么用：
- 将重试策略集中到单文件，后续调整次数或等待策略只改一处。
- 发布失败时日志信息更完整，便于快速判断是瞬时故障还是持续故障。

能防止什么问题：
- 两条发布链路各自维护重试函数导致配置漂移。
- 失败日志缺少重试上下文，排障时难以复现和归因。

说明：
- 本次为脚本复用与可观测性增强，不改变现有部署步骤顺序和目标路径。

### 6.15 发布链路重试参数按环境可配置

改了什么：
- 在 `.github/scripts/ci_retry.sh` 新增 `print_retry_policy()`，统一输出当前重试策略。
- 在 `.github/workflows/ci.yaml` 的 `devops-test` `restart` 步骤注入：
  1. `RETRY_MAX_ATTEMPTS: "3"`
  2. `RETRY_BASE_SLEEP_SECONDS: "5"`
- 在 `devops-prod` `restart` 步骤注入：
  1. `RETRY_MAX_ATTEMPTS: "5"`
  2. `RETRY_BASE_SLEEP_SECONDS: "5"`
- 两条链路在执行远程命令前都会打印当前重试策略（`print_retry_policy`）。

有什么用：
- 测试与生产可按稳定性诉求设置不同重试强度，避免“一刀切”。
- 日志直接展示当前策略值，便于排障和审计。

能防止什么问题：
- 测试环境重试过多拖慢反馈，或生产环境重试过少导致偶发失败。
- 仅看日志无法确认当次执行到底使用了哪套重试参数。

说明：
- 本次仅调整重试参数配置与日志输出，不改变部署命令本身。
