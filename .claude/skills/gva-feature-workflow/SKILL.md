---
name: gva-feature-workflow
description: "GVA 功能开发标准工作流。用于新建/扩展插件或业务模块时，按 Model->Service->API->Router->Initialize->Plugin->Web/Uni 顺序生成可上线代码，并强制检查 Swagger、Casbin、菜单、数据类型一致性。"
---

# GVA 功能开发标准工作流

## 适用场景
- 新建业务模块、插件、后台页面、uni 页面
- 迁移老功能到可复用架构（试衣、电影等）
- 修复“接口已写但权限/菜单/初始化遗漏”类问题

## 必做顺序
1. 先设计数据模型与请求模型
- 后端先写 model 与 model/request
- 统一字段类型，避免 model/request/response 类型不一致
- 指针与非指针转换在 service 层显式处理

2. 再写 Service
- 只做业务与数据库，不接触 gin.Context
- 统一返回 (data, error) 或 error

3. 再写 API
- 参数绑定与校验在 API 层
- 用统一 response 包返回
- 每个对外接口补全 Swagger 注释

4. 再写 Router
- 按 public/private 分组
- 鉴权路由统一挂 JWTAuth + CasbinHandler
- 需要审计的写接口挂 OperationRecord

5. 初始化与注册
- initialize/api.go: 注册 sys_apis
- initialize/menu.go: 注册菜单并分配角色
- initialize/casbin.go 或 init 文件: 注册并收敛规则（必要时清理历史脏权限）
- initialize/gorm.go: AutoMigrate
- plugin.go: init 注册插件 + Register 挂载路由
- plugin/register.go: 匿名导入激活插件

6. 前端落地
- web/src/api 按模块封装接口
- web/src/plugin/<name>/view 增加页面
- uni 端文案必须使用 i18n key，禁止新增硬编码
- 上传规则（格式/大小）优先走后端配置源

## 提交前检查清单
- 新增 API 是否同步注册到 sys_apis
- 新增 API 是否有对应 Casbin 规则
- 是否存在 authority=8080 越权路径
- 写接口是否记录操作日志
- uni 与 web 是否跟随后端单一配置源
- 是否存在中文硬编码未进入 i18n
