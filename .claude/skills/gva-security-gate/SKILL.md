---
name: gva-security-gate
description: "GVA 安全闸门检查流程。用于发布前或功能合并前，集中扫描越权、Casbin 误放权、公开配置泄露、WebSocket token 传递风险、i18n 回退硬编码、上传策略不一致。"
---

# GVA 安全闸门检查

## 目标
在上线前快速判断：是否存在高风险越权和配置泄露问题。

## 检查项
1. 路由权限分层
- public 路由只保留真正公开能力
- private 路由统一 JWT + Casbin
- 对关键写接口增加 API 层二次鉴权（角色白名单）

2. Casbin 规则收敛
- 禁止“for all authorities”默认放权关键接口
- 对 admin-only 接口只给 888/8881 等管理角色
- 初始化时补充历史脏规则清理 SQL

3. 系统配置对外暴露
- 公共 getByKey 必须白名单化
- 禁止通过公开接口读取安全、认证、内部运维配置

4. WebSocket 鉴权
- 优先 Authorization 或 Sec-WebSocket-Protocol
- 生产默认禁用 query token 回退
- 保留 Origin 校验与连接限流

5. 上传安全
- 后端统一检查 MIME + 扩展名 + 大小
- 前端只做用户体验前置校验，规则来源必须同后端配置

6. 多语言一致性
- 页面与提示语禁止直接写中文
- 后端中文错误建议映射为错误码或 i18n key

## 输出格式
- 风险级别：严重 / 高 / 中
- 每条给文件与行号
- 给最小修复建议与回归验证点
