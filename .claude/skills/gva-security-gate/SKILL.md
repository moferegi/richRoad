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
- Uni 的模板文案、toast、modal、placeholder、JS 拼接文本必须全部走 i18n
- 后端给 Uni 展示的字段若是多语言对象/JSON 字符串，必须通过 localText(value, langStore.locale) 解析
- 禁止把多语言对象直接 String(...) 渲染到页面
- 系统配置文案（重点：试衣间 announcement_content）必须验证是否支持语言切换
- 新功能上线前检查 layout/client/language（或项目实际 i18n 文件）是否补齐词条，避免缺词回退硬编码
- 扫描并拦截 raw msg 直透模式：data.msg ||、res.msg ||、e?.message ||
- 扫描并拦截后端 err.Error() 面向用户直接返回（尤其 client API）
- 上传/支付/下单/客服聊天等高频链路错误提示必须统一走 resolveApiMessage
- 提交前核对新增 key 在所有启用语种中齐全，不能只补部分语种

## 输出格式
- 风险级别：严重 / 高 / 中
- 每条给文件与行号
- 给最小修复建议与回归验证点
