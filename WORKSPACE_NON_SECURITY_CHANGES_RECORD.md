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

### 6.16 perf-comment-pass-uni-first：Uni 浏览历史渲染微优化

改了什么：
- 本批次命名为 `perf-comment-pass-uni-first`，用于后续继续按 `uni -> server -> web` 顺序推进“性能微优化 + 可读注释”。
- 在 `uni/src/pages/browseHistory/index.vue` 中，将浏览历史列表的图片地址、分组日期、展示时间、排序时间戳、稳定 key 统一提前归一化。
- 在 `uni/src/pages/tryon/history.vue` 中，将试衣历史列表的缩略图、展示时间、排序时间戳、记录类型统一提前归一化。
- 在 `uni/src/pages/collect/collect.vue` 中，将收藏列表的图片地址和稳定 key 统一提前归一化。
- 在 `uni/src/pages/integral/integral.vue` 中，将积分/试衣币记录的展示原因、展示时间和稳定 key 统一提前归一化。
- 在 `uni/src/pages/invite/index.vue` 中，将邀请好友列表的头像、展示名称、展示时间和稳定 key 统一提前归一化。
- 在 `uni/src/pages/coupon/index.vue` 中，将优惠券名称、描述、背景样式、折扣文案、门槛文案和稳定 key 统一提前归一化。
- 在 `uni/src/pages/kefu/index.vue` 中，将客服列表的展示名称、头像、头像兜底样式、状态文案和稳定 key 统一提前归一化。
- 在 `uni/src/pages/kefu/chat.vue` 中，将客服聊天页客服头像 URL 收敛到只读展示对象，避免每条客服消息重复解析头像地址。
- 在 `uni/src/pages/address/address.vue` 中，将地址列表的电话展示、地址详情和稳定 key 统一提前归一化，并复用原地区翻译重建入口。
- 在 `uni/src/pages/rechargeRecord/index.vue` 中，将试衣币充值记录的状态文案、展示时间、金额文案、支付方式文案和稳定 key 统一提前归一化。
- 在 `uni/src/pages/presale/list.vue` 中，将预售商品列表的图片地址、价格文案、进度样式和稳定 key 统一提前归一化，倒计时仍按秒动态刷新。
- 在 `uni/src/pages/order/order.vue` 中，将订单列表的订单时间、状态文案、商品缩略图、单品名称/描述、金额文案和稳定 key 统一提前归一化。
- 在 `uni/src/pages/evaluate/evaluate.vue` 中，将评价列表的头像、展示日期、图片预览 URL 列表和稳定 key 统一提前归一化。
- 在 `uni/src/pages/evaluate/evaluate.vue` 中，为评价列表图片展示项补充原始索引字段，模板预览不再依赖循环索引。
- 在 `uni/src/pages/orderDetail/orderDetail.vue` 中，将订单详情商品行的缩略图、名称、描述、规格文案和稳定 key 统一提前归一化。
- 在 `uni/src/pages/orderDetail/orderDetail.vue` 中，将订单详情商品行的本地化单价文案收敛到只读展示字段。
- 在 `uni/src/pages/orderDetail/orderDetail.vue` 中，将订单详情状态文案/class 和信息卡支付方式文案收敛到只读展示对象。
- 在 `uni/src/pages/evaluate/addEvaluate.vue` 中，将评价新增/查看页商品头部的图片、名称和规格文案收敛到缓存型展示对象。
- 在 `uni/src/pages/goodsDetails/components/goods-swiper.vue` 中，将商品轮播条目的稳定 key 和图片预览 URL 列表统一归一化缓存。
- 在 `uni/src/pages/goodsDetails/components/goods-swiper.vue` 中，为商品轮播条目补充原始索引字段，轮播点击预览入口不再依赖循环索引。
- 在 `uni/src/pages/goodsDetails/components/goods-swiper.vue` 中，将商品轮播文字叠加层的位置、颜色和字号样式收敛到只读展示字段。
- 在 `uni/src/pages/tabBar/components/swiper.vue` 中，将首页轮播的图片地址、遮罩文案、遮罩样式和稳定 key 统一提前归一化。
- 在 `uni/src/pages/tabBar/components/seckilling.vue` 中，将首页热卖横滑商品的图片地址、标题、价格文案和稳定 key 统一提前归一化。
- 在 `uni/src/pages/tabBar/components/category-page.vue` 中，将分类详情页横栏分类、商品卡片、标签样式和稳定 key 统一提前归一化。
- 在 `uni/src/pages/tabBar/components/category-page.vue` 中，为分类详情页横栏分类补充原始索引字段，模板选中和点击不再依赖循环索引。
- 在 `uni/src/pages/tabBar/category/category.vue` 中，将主分类页左侧分类、分区标题、商品卡片、首个标签和稳定 key 统一提前归一化。
- 在 `uni/src/pages/tabBar/category/category.vue` 中，为主分类页左侧分类补充原始索引字段，模板选中和点击不再依赖循环索引。
- 在 `uni/src/pages/signIn/signIn.vue` 中，将签到记录展示日期、记录稳定 key 和日历单元格 key 统一提前归一化。
- 在 `uni/src/pages/order/order-confirm.vue` 中，为订单确认页 mock 商品行补充稳定 key，避免继续使用索引 key。
- 在 `uni/src/pages/evaluate/evaluate-img.vue` 中，将评价图片缩略图 URL、预览 URL 列表和稳定 key 统一归一化缓存。
- 在 `uni/src/pages/evaluate/evaluate-img.vue` 中，为评价图片组件补充原始索引字段，模板预览入口不再依赖循环索引。
- 在 `uni/src/pages/player/index.vue` 中，将播放页封面图 URL 缓存为 `coverImageUrl`，供 hero、video poster 和无视频占位复用。
- 在 `uni/src/pages/myModel/index.vue` 中，将我的模特列表的图片地址、展示名称和稳定 key 统一提前归一化。
- 在 `uni/src/pages/myModel/index.vue` 中，将我的模特裁剪弹窗的舞台、裁剪框和隐藏 canvas 尺寸样式收敛到只读展示对象。
- 在 `uni/src/pages/myCloth/index.vue` 中，将我的衣橱列表的图片地址、展示名称、分类标签和稳定 key 统一提前归一化。
- 在 `uni/src/pages/myCloth/index.vue` 中，将我的衣橱裁剪弹窗的舞台、裁剪框和隐藏 canvas 尺寸样式收敛到只读展示对象。
- 在 `uni/src/pages/orderInfo/orderInfo.vue` 中，将订单确认信息页商品行的图片、名称、描述、规格、价格和稳定 key 收敛到缓存型展示列表。
- 在 `uni/src/pages/tabBar/shop/components/shop-goods-list.vue` 中，将购物车商品行的图片、名称、规格/描述、单价和稳定 key 统一提前归一化。
- 在 `uni/src/pages/tabBar/shop/components/shop-goods-list.vue` 中，将购物车商品选择入口改为直接切换当前商品项，模板不再依赖循环索引。
- 在 `uni/src/pages/pay/index.vue` 中，将支付页二维码名称/图片地址、首选支付方式图片地址和稳定 key 统一提前归一化。
- 在 `uni/src/pages/pay/index.vue` 中，为支付页二维码列表补充原始索引字段，二维码切换入口不再依赖循环索引。
- 在 `uni/src/pages/pay/index.vue` 中，将支付页付款提示文本和提示样式收敛到只读展示对象。
- 在 `uni/src/pages/tabBar/components/categories.vue` 中，为首页分类 tab 展示列表补充稳定 key，避免继续使用索引 key。
- 在 `uni/src/pages/tabBar/components/categories.vue` 中，为首页分类 tab 展示列表补充原始索引字段，模板选中和滚动不再依赖循环索引。
- 在 `uni/src/pages/evaluate/addEvaluate.vue` 中，将评价图片草稿/查看图片的展示 URL、预览 URL 和稳定 key 收敛到缓存型图片视图。
- 在 `uni/src/pages/orderDetail/orderDetail.vue` 中，将订单创建时间、支付时间、收货时间收敛到只读展示对象，减少信息卡片重复格式化时间。
- 在 `uni/src/pages/goodsDetails/goodsDetails.vue` 中，将商品详情优惠券弹层的名称、门槛文案、有效期文案和稳定 key 收敛到缓存型展示列表。
- 在 `uni/src/pages/goodsDetails/goodsDetails.vue` 中，将商品属性列表的属性名、属性值和稳定 key 收敛到 `parsedAttrs` 派生字段。
- 在 `uni/src/pages/user/profile.vue` 中，将个人资料页注册时间收敛到只读展示字段，避免模板直接格式化时间。
- 在 `uni/src/pages/evaluate/addEvaluate.vue` 中，将查看评价时的商家回复时间收敛到只读展示字段，避免模板直接格式化时间。
- 在 `uni/src/pages/order/order.vue` 中，为订单状态 tab 补充稳定 key，避免继续使用索引 key。
- 在 `uni/src/pages/orderInfo/orderInfo.vue` 中，将订单确认页优惠券弹层的金额、门槛文案、名称、有效期和稳定 key 收敛到缓存型展示列表。
- 在 `uni/src/pages/player/index.vue` 中，将播放页标题、简介、剧场当前集标题和简介收敛到缓存型展示对象。
- 在 `uni/src/pages/player/index.vue` 中，将播放页剧场快选和主选集列表的集数文案、稳定 key、原始索引收敛到缓存型展示列表。
- 在 `uni/src/pages/goodsDetails/goodsDetails.vue` 中，将商品详情头部标题、描述和预售时间文案收敛到缓存型展示对象。
- 在 `uni/src/pages/collect/collect.vue` 中，将收藏商品卡片标题收敛到归一化展示字段，避免模板直接解析多语言标题。
- 在 `uni/src/pages/presale/list.vue` 中，将预售商品卡片标题收敛到归一化展示字段，避免模板直接解析多语言标题。
- 在 `uni/src/pages/browseHistory/index.vue` 中，将浏览历史卡片标题收敛到归一化展示字段，避免模板直接解析多语言标题。
- 在 `uni/src/pages/maintenance/index.vue` 中，将维护页弹窗标题和内容收敛到只读展示对象，避免模板直接解析多语言配置文案。
- 在 `uni/src/pages/evaluate/addEvaluate.vue` 中，将评价页导航标题、评分文案、图片标题、占位文案和提交按钮文案收敛到只读展示对象。
- 在 `uni/src/pages/evaluate/addEvaluate.vue` 中，为评价图片展示视图补充原始索引字段，模板预览和删除不再依赖循环索引。
- 在 `uni/src/pages/evaluate/orderEvaluate.vue` 中，将订单评价项商品图、商品名、评分文案、已上传图片 URL 和列表 key 收敛到只读展示字段。
- 在 `uni/src/pages/address/addAddress.vue` 中，将新增地址页区号弹窗国家名和列表 key 收敛到只读展示字段。
- 在 `uni/src/pages/address/editAddress.vue` 中，将编辑地址页区号弹窗国家名和列表 key 收敛到只读展示字段。
- 在 `uni/src/pages/user/login.vue` 中，将登录页手机号区号弹窗国家名和列表 key 收敛到只读展示字段。
- 在 `uni/src/pages/user/register.vue` 中，将注册页手机号区号弹窗国家名和列表 key 收敛到只读展示字段。
- 在 `uni/src/pages/tryon/index.vue` 中，将当前试衣任务图片和任务列表时间收敛到只读展示字段。
- 在 `uni/src/pages/tryon/generate.vue` 中，将结果对比层宽度、分割线位置和引导手势位置样式收敛到只读展示对象。
- 在 `uni/src/pages/pay/index.vue` 中，将支付步骤序号、文案和列表 key 收敛到只读展示字段。
- 在 `uni/src/pages/tryon/history.vue` 中，将试衣历史普通预览和对比预览图片 URL 收敛到只读展示字段。
- 在 `uni/src/pages/tryon/history.vue` 中，将试衣历史对比层宽度、分割线位置和引导手势位置样式收敛到只读展示对象。
- 在 `uni/src/pages/tryon/history.vue` 中，将试衣历史 H5 自定义预览图片、序号文案和列表 key 收敛到只读展示字段。
- 在 `uni/src/pages/tabBar/my/index.vue` 中，将充值套餐弹窗标题、价格和列表 key 收敛到只读展示字段。
- 在 `uni/src/pages/tabBar/index.vue` 中，将试衣间教程弹窗的序号文案和列表 key 收敛到只读展示字段。
- 在 `uni/src/pages/tabBar/index.vue` 中，将试衣间首页上传抽屉静态示例卡片的预览 URL 和稳定 key 收敛到只读展示字段。
- 在 `uni/src/pages/tabBar/index.vue` 中，将试衣间首页上传抽屉我的模特/衣橱条目的预览 URL 和稳定 key 收敛到只读展示字段。
- 在 `uni/src/pages/tabBar/index.vue` 中，将试衣间首页上传图片大小遮罩文案收敛到只读展示字段。
- 在 `uni/src/pages/tabBar/index.vue` 中，将试衣间首页上传抽屉裁剪舞台、裁剪框和隐藏 canvas 尺寸样式收敛到只读展示对象。
- 在 `uni/src/pages/goodsDetails/components/goods-sku.vue` 中，将 SKU 规格组标题、规格值文案和列表 key 收敛到只读展示字段。
- 在 `uni/src/pages/tabBar/clothes/index.vue` 中，将衣橱商城瀑布流卡片的图片、名称、卡片样式和列表 key 收敛到只读展示字段。
- 在 `uni/src/pages/tabBar/clothes/index.vue` 中，将衣橱商城分类 tab 的名称、DOM id、分类 id 和列表 key 收敛到只读展示字段。
- 在 `uni/src/pages/tabBar/shop/shop.vue` 中，将鞋靴试穿上传抽屉示例卡片的预览 URL 和稳定 key 收敛到只读展示字段。
- 在 `uni/src/pages/tabBar/shop/shop.vue` 中，将鞋靴试穿上传抽屉我的模特/鞋靴条目的预览 URL 和稳定 key 收敛到只读展示字段。
- 在 `uni/src/pages/tabBar/shop/shop.vue` 中，将鞋靴试穿上传图片大小遮罩文案收敛到只读展示字段。
- 列表模板改用归一化后的派生字段，减少滚动渲染时重复解析日期和拼接图片 URL。
- 为本次新增的归一化、时间处理和 CSS 截断兼容逻辑补充注释，说明用途、影响边界和关联点。
- 为商品标题补充标准 `line-clamp` 属性，保留原 `-webkit-line-clamp`，提升跨端 CSS 兼容性。

有什么用：
- 降低浏览历史页在滚动和分页展示时的重复计算成本。
- 降低试衣历史页在滚动筛选时重复解析缩略图、时间和记录类型的成本。
- 降低收藏页列表渲染时重复解析图片 URL 的成本。
- 降低积分记录页在列表渲染时重复执行原因翻译和时间格式化的成本。
- 降低邀请页好友列表重复格式化头像、名称和注册时间的成本。
- 降低优惠券页重复解析多语言文案、背景图和金额门槛文案的成本。
- 降低客服列表页重复解析头像、状态、多语言名称和头像兜底样式的成本。
- 降低客服聊天页多条客服消息重复解析客服头像 URL 的成本，同时消息列表、发送状态、图片预览和输入发送流程仍读取原始状态。
- 降低地址列表页语言切换后重复拼接电话、地区与街道文案的成本。
- 降低试衣币充值记录页重复格式化状态、金额、时间和支付方式的成本。
- 降低预售列表页重复解析图片、价格和进度样式的成本，同时保持倒计时实时计算。
- 降低订单列表页重复解析商品缩略图、金额、状态和订单时间的成本。
- 降低评价列表页重复解析头像、评价图片 URL 和日期格式化的成本。
- 降低评价列表页图片预览继续依赖模板循环索引的节点复用风险，同时预览图片集合仍读取归一化后的原始 URL 列表。
- 降低订单详情页重复解析商品图、规格和多语言商品文案的成本。
- 降低订单详情页商品行重复计算本地化单价文案的成本，同时订单金额汇总、支付金额和订单快照币种逻辑仍保持不变。
- 降低订单详情页状态和信息卡支付方式重复生成展示文案/class 的成本，同时支付、退款、倒计时和客服跳转仍读取原始订单状态与支付方式。
- 降低评价新增/查看页商品头部重复解析 SKU 图片、多语言名称和规格文案的成本，同时不影响用户正在编辑的评价草稿。
- 降低商品轮播图点击预览时重复过滤和映射图片 URL 的成本，视频签名完成后仍会同步刷新预览缓存。
- 降低商品轮播继续依赖模板循环索引的节点复用风险，同时图片预览、视频签名和预览 URL 缓存仍读取归一化后的轮播条目。
- 降低商品轮播文字叠加层在渲染时重复创建样式对象的成本，同时视频签名、图片预览和轮播原始数据仍保持不变。
- 降低首页轮播重复解析图片地址、遮罩多语言文案和内联样式对象的成本。
- 降低首页热卖横滑列表重复解析图片地址、多语言标题和本地化价格的成本。
- 降低分类详情页重复解析分类图标、商品图片、多语言标题、标签样式和本地化价格的成本。
- 降低分类详情页横栏继续依赖模板循环索引的节点复用风险，同时分类选中、分类 ID 切换和商品分页加载仍使用原始分类数据。
- 降低主分类页左右栏切换和商品网格滚动时重复解析多语言分类、商品文案、图片 URL 和价格的成本。
- 降低主分类页左侧分类继续依赖模板循环索引的节点复用风险，同时分类选中、当前分类同步和右侧商品加载仍使用原始分类数据。
- 降低签到记录列表重复格式化签到日期的成本，并让日历单元格/记录列表 key 更稳定。
- 降低订单确认页 mock 商品行因索引 key 带来的节点复用风险。
- 降低评价图片组件缩略图渲染和预览时重复解析图片地址的成本。
- 降低评价图片组件继续依赖模板循环索引的节点复用风险，同时预览顺序和预览 URL 列表仍来自原始图片列表派生结果。
- 降低播放页多个封面节点重复调用 `getUrl(data.imageUrl)` 的成本。
- 降低我的模特列表重复解析图片地址和兜底名称的成本，预览复用同一展示 URL。
- 降低我的模特裁剪弹窗重复创建尺寸样式对象的成本，同时上传、裁剪坐标、导出 canvas 和删除/选择流程仍读取原始裁剪状态。
- 降低我的衣橱列表重复解析图片地址、分类文案和兜底名称的成本，预览复用同一展示 URL。
- 降低我的衣橱裁剪弹窗重复创建尺寸样式对象的成本，同时上传、裁剪坐标、导出 canvas 和删除流程仍读取原始裁剪状态。
- 降低订单确认信息页重复解析商品图、多语言商品文案、规格和本地化价格的成本，同时保留下单用原始商品明细。
- 降低购物车列表重复解析商品图、多语言名称、规格/描述和本地化单价的成本，同时保留数量变更、选中和结算用原始字段。
- 降低购物车商品行继续依赖模板循环索引的节点复用风险，同时选择、全选、删除和结算仍操作当前购物车项与原始商品字段。
- 降低支付页二维码切换、预览/保存和首选支付方式展示时重复解析多语言名称与图片地址的成本。
- 降低支付页二维码 tab 继续依赖模板循环索引的节点复用风险，同时当前二维码、预览、保存和支付确认仍使用原 `currentQrIndex` 与归一化二维码列表。
- 降低支付页付款提示区域重复创建样式对象的成本，同时后台提示配置读取、多语言解析和支付流程仍保持不变。
- 降低首页分类 tab 顺序变化时使用索引 key 带来的节点复用风险。
- 降低首页分类 tab 继续依赖模板循环索引的节点复用风险，同时分类选中、`v-model` 同步和滚动居中仍按原索引执行。
- 降低评价新增/查看页图片网格渲染和预览时重复解析图片地址的成本，同时保留原图片草稿数组。
- 降低订单详情页信息卡片重复执行时间格式化的成本，同时保留原始订单时间字段供业务判断使用。
- 降低商品详情优惠券弹层重复解析多语言名称、门槛文案和有效期文案的成本，同时领取/选用仍回传原优惠券对象。
- 降低商品详情属性列表重复解析多语言属性名/属性值的成本，并避免属性顺序变化时继续使用索引 key。
- 降低个人资料页重复格式化注册日期的成本，同时保留原始用户时间字段供资料同步或后续接口扩展使用。
- 降低评价查看页商家回复区域重复格式化回复时间的成本，同时保留原始回复时间字段供详情重拉使用。
- 降低订单状态 tab 顺序或文案变化时使用索引 key 带来的节点复用风险，同时状态筛选仍使用原 id 字段。
- 降低订单确认页优惠券弹层重复格式化金额、门槛文案、名称和有效期的成本，同时领取、取消选择和提交订单仍使用原优惠券对象。
- 降低播放页标题、简介和剧场当前集文案重复执行多语言解析的成本，同时播放地址、选集、进度和收藏仍读取原始数据。
- 降低播放页两处选集列表重复生成集数文案的成本，并用稳定 key 减少选集节点复用风险；播放切换仍按原始索引执行。
- 降低商品详情头部重复解析多语言标题/描述和重复格式化预售时间的成本，同时购买、SKU、收藏和浏览记录仍读取原始商品数据。
- 降低收藏列表卡片重复解析多语言标题的成本，同时取消收藏、分页追加和跳转详情仍读取原始收藏商品字段。
- 降低预售列表卡片重复解析多语言标题的成本，同时倒计时、分页追加和跳转详情仍读取原始预售商品字段。
- 降低浏览历史卡片重复解析多语言标题的成本，同时历史加载、分页显示、清空和跳转详情仍读取原始历史记录字段。
- 降低维护页弹窗重复解析多语言配置文案的成本，同时维护开关、背景图和按钮跳转仍读取原始配置字段。
- 降低评价新增/查看页固定文案重复取翻译和拼接上传标题的成本，同时评分、提交、图片上传/删除/预览和查看模式判断仍读取原始响应式状态。
- 降低评价新增/查看页图片网格继续依赖模板循环索引的节点复用风险，同时图片预览、删除、上传和提交仍按原 `pics` 数组执行。
- 降低订单评价页多商品评价时重复解析商品图、商品名、评分文案和图片 URL 的成本，同时评分、评论输入、图片增删、预览和提交仍按原 `orderItems` 索引执行。
- 降低新增地址页区号弹窗重复解析多语言国家名的成本，同时选择区号和保存地址仍使用原始区号对象。
- 降低编辑地址页区号弹窗重复解析多语言国家名的成本，同时区号回显、选择区号和保存地址仍使用原始表单字段。
- 降低登录页手机号区号弹窗重复解析多语言国家名的成本，同时登录提交仍使用原始 `form.areaCode` 字段。
- 降低注册页手机号区号弹窗重复解析多语言国家名的成本，同时注册提交仍使用原始 `form.areaCode` 字段。
- 降低试衣页当前任务图片重复解析 URL 和任务列表重复格式化创建时间的成本，同时预览、重试、轮询刷新和选中任务仍读取原始任务字段。
- 降低试衣生成页结果对比层随滑块变化时重复创建宽度/位置样式对象的成本，同时对比预览、拖拽、缩放和滑块取值仍读取原始对比状态。
- 降低支付页步骤列表使用索引 key 带来的节点复用风险，同时支付方式选择、二维码保存、复制草稿和联系客服动作仍读取原始支付状态。
- 降低试衣历史预览弹窗和对比层重复解析图片 URL 的成本，同时下载、H5 预览、索引切换和对比拖拽缩放仍读取原始预览数组。
- 降低试衣历史对比层随滑块变化时重复创建宽度/位置样式对象的成本，同时下载、H5 预览、索引切换和对比拖拽缩放仍读取原始预览状态。
- 降低试衣历史 H5 自定义预览重复拼接序号文案和索引 key 的成本，同时当前索引、长按下载和当前图下载仍使用原始预览数组。
- 降低我的页充值套餐弹窗重复解析套餐点数、币种文案和本地化价格的成本，同时创建充值订单和支付方式选择仍使用原始套餐对象。
- 降低试衣间教程弹窗列表使用索引 key 带来的节点复用风险，同时教程弹窗开关、试衣上传和任务提交流程仍读取原始状态。
- 降低试衣间首页上传抽屉静态示例卡片在渲染时重复组装候选图片域名的成本，同时示例应用、图片失败重试、我的模特/衣橱选择和上传流程仍使用原始 URL。
- 降低试衣间首页上传抽屉我的模特/衣橱列表在渲染时重复组装候选图片域名的成本，同时选择模特、选择衣橱和上传目标赋值仍使用原始条目 URL。
- 降低试衣间首页三处上传图片大小遮罩重复格式化字节数的成本，同时上传、清空和任务提交仍读取原始 size bytes。
- 降低试衣间首页上传抽屉裁剪区域重复创建尺寸样式对象的成本，同时本地/远程上传、裁剪坐标、导出 canvas 和任务提交仍读取原始裁剪状态。
- 降低商品 SKU 弹层重复解析规格组标题和规格值多语言文案的成本，并避免规格组/规格值继续使用索引 key；选规格、禁用状态、加入购物车、立即购买和试穿入口仍使用原始 `specGroups`。
- 降低衣橱商城瀑布流卡片重复解析商品图片、多语言名称和左右列样式的成本，同时试穿、购买跳转和图片预览仍使用原始商品对象。
- 降低衣橱商城分类 tab 重复解析多语言名称和 DOM id 的成本，同时分类切换、滚动居中和分页重载仍使用原始分类对象。
- 降低鞋靴试穿上传抽屉示例卡片在渲染时重复组装候选图片域名的成本，同时示例应用、图片失败重试和本地/远程上传流程仍使用原始示例 URL。
- 降低鞋靴试穿上传抽屉我的模特/鞋靴列表在渲染时重复组装候选图片域名的成本，同时选择模特、选择鞋靴和上传目标赋值仍使用原始条目 URL。
- 降低鞋靴试穿三处上传图片大小遮罩重复格式化字节数的成本，同时上传、清空和任务提交仍读取原始 size bytes。
- 稳定列表 key，减少 Vue diff 时因索引 key 带来的不必要节点复用风险。
- 后续维护者新增时间字段、图片字段或调整展示格式时，可以从归一化函数集中修改。

能防止什么问题：
- 大列表滚动时频繁执行 `new Date()`、URL 解析和日期格式化导致的轻微卡顿累积。
- 使用数组索引作为 key 时，列表顺序变化可能造成的渲染复用异常。
- 首屏列表和追加分页数据字段形态不一致，导致后续维护时漏改某一路径。
- 只写浏览器前缀截断属性导致的兼容性诊断噪音。

验证：
- `cd uni && npx eslint src/pages/browseHistory/index.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/tryon/history.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/collect/collect.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/integral/integral.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/invite/index.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/coupon/index.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/kefu/index.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/kefu/chat.vue`：退出码 `0`（客服聊天头像 URL 展示缓存）。
- `cd uni && npx eslint src/pages/address/address.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/rechargeRecord/index.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/presale/list.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/order/order.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/evaluate/evaluate.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/evaluate/evaluate.vue`：退出码 `0`（评价列表图片原始索引字段缓存）。
- `cd uni && npx eslint src/pages/orderDetail/orderDetail.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/orderDetail/orderDetail.vue`：退出码 `0`（订单详情商品行本地化单价展示字段缓存）。
- `cd uni && npx eslint src/pages/orderDetail/orderDetail.vue`：退出码 `0`（订单详情状态/支付方式展示对象缓存）。
- `cd uni && npx eslint src/pages/evaluate/addEvaluate.vue src/pages/goodsDetails/components/goods-swiper.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/goodsDetails/components/goods-swiper.vue`：退出码 `0`（商品轮播原始索引字段缓存）。
- `cd uni && npx eslint src/pages/goodsDetails/components/goods-swiper.vue`：退出码 `0`（商品轮播文字样式展示字段缓存）。
- `cd uni && npx eslint src/pages/tabBar/components/swiper.vue src/pages/tabBar/components/seckilling.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/tabBar/components/category-page.vue src/pages/tabBar/category/category.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/tabBar/components/category-page.vue`：退出码 `0`（分类详情横栏原始索引字段缓存）。
- `cd uni && npx eslint src/pages/tabBar/category/category.vue`：退出码 `0`（主分类左栏原始索引字段缓存）。
- `cd uni && npx eslint src/pages/signIn/signIn.vue src/pages/order/order-confirm.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/evaluate/evaluate-img.vue src/pages/player/index.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/evaluate/evaluate-img.vue`：退出码 `0`（评价图片组件原始索引字段缓存）。
- `cd uni && npx eslint src/pages/myModel/index.vue src/pages/myCloth/index.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/myModel/index.vue`：退出码 `0`（我的模特裁剪尺寸样式缓存）。
- `cd uni && npx eslint src/pages/myCloth/index.vue`：退出码 `0`（我的衣橱裁剪尺寸样式缓存）。
- `cd uni && npx eslint src/pages/orderInfo/orderInfo.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/tabBar/shop/components/shop-goods-list.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/tabBar/shop/components/shop-goods-list.vue`：退出码 `0`（购物车商品行选择入口去循环索引）。
- `cd uni && npx eslint src/pages/pay/index.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/pay/index.vue`：退出码 `0`（支付二维码原始索引字段缓存）。
- `cd uni && npx eslint src/pages/pay/index.vue`：退出码 `0`（支付提示展示对象缓存）。
- `cd uni && npx eslint src/pages/tabBar/components/categories.vue src/pages/evaluate/addEvaluate.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/tabBar/components/categories.vue`：退出码 `0`（首页分类 tab 原始索引字段缓存）。
- `cd uni && npx eslint src/pages/orderDetail/orderDetail.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/goodsDetails/goodsDetails.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/goodsDetails/goodsDetails.vue`：退出码 `0`（商品属性列表展示字段缓存）。
- `cd uni && npx eslint src/pages/user/profile.vue`：退出码 `0`。
- `cd uni && npx eslint src/pages/evaluate/addEvaluate.vue`：退出码 `0`（商家回复时间展示字段缓存）。
- `cd uni && npx eslint src/pages/order/order.vue`：退出码 `0`（订单状态 tab 稳定 key）。
- `cd uni && npx eslint src/pages/orderInfo/orderInfo.vue`：退出码 `0`（订单确认优惠券弹层展示字段缓存）。
- `cd uni && npx eslint src/pages/player/index.vue`：退出码 `0`（播放页标题/简介展示字段缓存）。
- `cd uni && npx eslint src/pages/player/index.vue`：退出码 `0`（播放页选集展示列表缓存）。
- `cd uni && npx eslint src/pages/goodsDetails/goodsDetails.vue`：退出码 `0`（商品详情头部展示字段缓存）。
- `cd uni && npx eslint src/pages/collect/collect.vue`：退出码 `0`（收藏卡片标题展示字段缓存）。
- `cd uni && npx eslint src/pages/presale/list.vue`：退出码 `0`（预售卡片标题展示字段缓存）。
- `cd uni && npx eslint src/pages/browseHistory/index.vue`：退出码 `0`（浏览历史卡片标题展示字段缓存）。
- `cd uni && npx eslint src/pages/maintenance/index.vue`：退出码 `0`（维护页弹窗文案展示字段缓存）。
- `cd uni && npx eslint src/pages/evaluate/addEvaluate.vue`：退出码 `0`（评价页固定文案展示字段缓存）。
- `cd uni && npx eslint src/pages/evaluate/addEvaluate.vue`：退出码 `0`（评价图片展示原始索引字段缓存）。
- `cd uni && npx eslint src/pages/evaluate/orderEvaluate.vue`：退出码 `0`（订单评价项展示字段缓存）。
- `cd uni && npx eslint src/pages/address/addAddress.vue`：退出码 `0`（新增地址页区号展示字段缓存）。
- `cd uni && npx eslint src/pages/address/editAddress.vue`：退出码 `0`（编辑地址页区号展示字段缓存）。
- `cd uni && npx eslint src/pages/user/login.vue`：退出码 `0`（登录页区号展示字段缓存）。
- `cd uni && npx eslint src/pages/user/register.vue`：退出码 `0`（注册页区号展示字段缓存）。
- `cd uni && npx eslint src/pages/tryon/index.vue`：退出码 `0`（试衣页任务图片和时间展示字段缓存）。
- `cd uni && npx eslint src/pages/tryon/generate.vue`：退出码 `0`（试衣生成页对比百分比样式缓存）。
- `cd uni && npx eslint src/pages/pay/index.vue`：退出码 `0`（支付步骤展示字段稳定 key）。
- `cd uni && npx eslint src/pages/tryon/history.vue`：退出码 `0`（试衣历史预览图片 URL 展示字段缓存）。
- `cd uni && npx eslint src/pages/tryon/history.vue`：退出码 `0`（试衣历史对比百分比样式缓存）。
- `cd uni && npx eslint src/pages/tryon/history.vue`：退出码 `0`（试衣历史 H5 自定义预览展示字段缓存）。
- `cd uni && npx eslint src/pages/tabBar/my/index.vue`：退出码 `0`（我的页充值套餐展示字段缓存）。
- `cd uni && npx eslint src/pages/tabBar/index.vue`：退出码 `0`（试衣间教程弹窗展示字段缓存）。
- `cd uni && npx eslint src/pages/tabBar/index.vue`：退出码 `0`（试衣间首页示例卡片预览 URL 展示字段缓存）。
- `cd uni && npx eslint src/pages/tabBar/index.vue`：退出码 `0`（试衣间首页我的模特/衣橱预览 URL 展示字段缓存）。
- `cd uni && npx eslint src/pages/tabBar/index.vue src/pages/tabBar/shop/shop.vue`：退出码 `0`（试衣间/鞋靴试穿图片大小遮罩展示字段缓存）。
- `cd uni && npx eslint src/pages/tabBar/index.vue`：退出码 `0`（试衣间首页裁剪尺寸样式缓存）。
- `cd uni && npx eslint src/pages/goodsDetails/components/goods-sku.vue`：退出码 `0`（商品 SKU 规格展示字段缓存）。
- `cd uni && npx eslint src/pages/tabBar/clothes/index.vue`：退出码 `0`（衣橱商城瀑布流卡片展示字段缓存）。
- `cd uni && npx eslint src/pages/tabBar/clothes/index.vue`：退出码 `0`（衣橱商城分类 tab 展示字段缓存）。
- `cd uni && npx eslint src/pages/tabBar/shop/shop.vue`：退出码 `0`（鞋靴试穿示例卡片预览 URL 展示字段缓存）。
- `cd uni && npx eslint src/pages/tabBar/shop/shop.vue`：退出码 `0`（鞋靴试穿我的模特/鞋靴预览 URL 展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/browseHistory/index.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/tryon/history.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/collect/collect.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/integral/integral.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/invite/index.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/coupon/index.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/kefu/index.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/kefu/chat.vue` 无错误（客服聊天头像 URL 展示缓存）。
- VS Code 目标文件诊断：`uni/src/pages/address/address.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/rechargeRecord/index.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/presale/list.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/order/order.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/evaluate/evaluate.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/evaluate/evaluate.vue` 无错误（评价列表图片原始索引字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/orderDetail/orderDetail.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/orderDetail/orderDetail.vue` 无错误（订单详情商品行本地化单价展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/orderDetail/orderDetail.vue` 无错误（订单详情状态/支付方式展示对象缓存）。
- VS Code 目标文件诊断：`uni/src/pages/evaluate/addEvaluate.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/goodsDetails/components/goods-swiper.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/goodsDetails/components/goods-swiper.vue` 无错误（商品轮播原始索引字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/goodsDetails/components/goods-swiper.vue` 无错误（商品轮播文字样式展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/components/swiper.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/components/seckilling.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/components/category-page.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/components/category-page.vue` 无错误（分类详情横栏原始索引字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/category/category.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/category/category.vue` 无错误（主分类左栏原始索引字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/signIn/signIn.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/order/order-confirm.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/evaluate/evaluate-img.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/evaluate/evaluate-img.vue` 无错误（评价图片组件原始索引字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/player/index.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/myModel/index.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/myModel/index.vue` 无错误（我的模特裁剪尺寸样式缓存）。
- VS Code 目标文件诊断：`uni/src/pages/myCloth/index.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/myCloth/index.vue` 无错误（我的衣橱裁剪尺寸样式缓存）。
- VS Code 目标文件诊断：`uni/src/pages/orderInfo/orderInfo.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/shop/components/shop-goods-list.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/shop/components/shop-goods-list.vue` 无错误（购物车商品行选择入口去循环索引）。
- VS Code 目标文件诊断：`uni/src/pages/pay/index.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/pay/index.vue` 无错误（支付二维码原始索引字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/pay/index.vue` 无错误（支付提示展示对象缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/components/categories.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/components/categories.vue` 无错误（首页分类 tab 原始索引字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/evaluate/addEvaluate.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/orderDetail/orderDetail.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/goodsDetails/goodsDetails.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/goodsDetails/goodsDetails.vue` 无错误（商品属性列表展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/user/profile.vue` 无错误。
- VS Code 目标文件诊断：`uni/src/pages/evaluate/addEvaluate.vue` 无错误（商家回复时间展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/order/order.vue` 无错误（订单状态 tab 稳定 key）。
- VS Code 目标文件诊断：`uni/src/pages/orderInfo/orderInfo.vue` 无错误（订单确认优惠券弹层展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/player/index.vue` 无错误（播放页标题/简介展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/player/index.vue` 无错误（播放页选集展示列表缓存）。
- VS Code 目标文件诊断：`uni/src/pages/goodsDetails/goodsDetails.vue` 无错误（商品详情头部展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/collect/collect.vue` 无错误（收藏卡片标题展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/presale/list.vue` 无错误（预售卡片标题展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/browseHistory/index.vue` 无错误（浏览历史卡片标题展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/maintenance/index.vue` 无错误（维护页弹窗文案展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/evaluate/addEvaluate.vue` 无错误（评价页固定文案展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/evaluate/addEvaluate.vue` 无错误（评价图片展示原始索引字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/evaluate/orderEvaluate.vue` 无错误（订单评价项展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/address/addAddress.vue` 无错误（新增地址页区号展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/address/editAddress.vue` 无错误（编辑地址页区号展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/user/login.vue` 无错误（登录页区号展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/user/register.vue` 无错误（注册页区号展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tryon/index.vue` 无错误（试衣页任务图片和时间展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tryon/generate.vue` 无错误（试衣生成页对比百分比样式缓存）。
- VS Code 目标文件诊断：`uni/src/pages/pay/index.vue` 无错误（支付步骤展示字段稳定 key）。
- VS Code 目标文件诊断：`uni/src/pages/tryon/history.vue` 无错误（试衣历史预览图片 URL 展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tryon/history.vue` 无错误（试衣历史对比百分比样式缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tryon/history.vue` 无错误（试衣历史 H5 自定义预览展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/my/index.vue` 无错误（我的页充值套餐展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/index.vue` 无错误（试衣间教程弹窗展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/index.vue` 无错误（试衣间首页示例卡片预览 URL 展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/index.vue` 无错误（试衣间首页我的模特/衣橱预览 URL 展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/index.vue` 无错误（试衣间首页图片大小遮罩展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/index.vue` 无错误（试衣间首页裁剪尺寸样式缓存）。
- VS Code 目标文件诊断：`uni/src/pages/goodsDetails/components/goods-sku.vue` 无错误（商品 SKU 规格展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/clothes/index.vue` 无错误（衣橱商城瀑布流卡片展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/clothes/index.vue` 无错误（衣橱商城分类 tab 展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/shop/shop.vue` 无错误（鞋靴试穿示例卡片预览 URL 展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/shop/shop.vue` 无错误（鞋靴试穿我的模特/鞋靴预览 URL 展示字段缓存）。
- VS Code 目标文件诊断：`uni/src/pages/tabBar/shop/shop.vue` 无错误（鞋靴试穿图片大小遮罩展示字段缓存）。

说明：
- 本次不改接口调用、分页数量、清空/删除历史、浏览历史加载/分页显示/跳转详情、取消收藏、跳转详情、邀请分享、积分统计、签到提交/分页/日历月份切换、优惠券领取/使用、客服联系方式选择、客服聊天连接/收发消息/图片上传预览/评价、地址选择/编辑/删除/新增保存、地址编辑回显与保存、登录/注册配置、验证码、账号/手机号登录注册提交、试衣上传/提交/轮询/重试/预览、衣橱商城搜索/分类筛选/分页/试穿/购买/预览、试衣间教程弹窗开关、试衣历史下载/H5 预览/对比拖拽缩放、充值支付/取消/退款/物流/评价、我的页充值下单/支付方式选择、订单状态筛选/支付/取消/退款/物流/评价、订单详情支付/退款/状态判断、订单确认信息页下单参数/提交跳转/优惠券领取与选择、订单评价提交/图片上传删除预览、支付方式选择/订单支付方式同步/确认已付款、支付页二维码保存/复制草稿/联系客服、维护模式开关/背景图/客服与回首页跳转、评价图片上传/删除/预览入口、评价提交/查看详情重拉、播放器选集/收藏/播放控制/进度记录、我的模特/衣橱上传裁剪/删除/选择流程、商品轮播视频签名、首页轮播/热卖/分类跳转、预售列表分页/倒计时/跳转详情、商品详情 SKU 选择/购买/加购/试穿入口、商品详情收藏/浏览记录/优惠券领取选用、个人资料昵称/邮箱/手机号修改与登出流程、多语言展示和价格格式化逻辑。
- 派生字段仅用于当前页面渲染，不写回 Pinia store 或后端返回数据。
