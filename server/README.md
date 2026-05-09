
## 项目介绍

本项目为基于Gin-Vue-Admin 和 [GVA-uniapp开发专用基座](https://plugin.gin-vue-admin.com/#/layout/newPluginInfo?id=71)开发的UNIAPP 商城系统

项目购买后联系奇淼加入维护专用仓库
注：目前暂时无法加入仓库,代码中包含较为敏感的aksk等信息，等剔除完成后会开放仓库加入

## 主要板块为

### uniapp端

- 首页
- 登录
- 注册
- 轮播图
- 商品列表
- 商品详情
- 商品细分选择
- 购物车
- 购物车下单
- 单品下单
- 付款流程
- 评价体系
- 订单管理
- 物流查看
- 地址管理
- 个人中心

### 后端

- 客户端管理
- 客户端用户管理
- 商城管理
- 用户评价管理
- 订单管理
- 轮播图管理
- 商品分类管理
- 商品管理
- 商品规格
- 商品属性
- sku添加管理

# 安装教程

## 客户端
所在文件位置 `/uni`

```shell
# 安装依赖
npm install

```
需要使用hbuildx启动项目

## 管理端前端
所在文件位置 `/web`

```shell
# 安装依赖
npm install

# 启动项目
npm run dev
```

## 数据库文件

直接按照gva模式点击初始化即可初始所有数据

## 配置信息
需要手动填入 config.yaml内的微信支付配置和顺丰(非必须)配置
微信支付使用本插件 https://plugin.gin-vue-admin.com/#/layout/newPluginInfo?id=22
```yaml
wxpay:
    mch-id: ""  # 商户ID
    app-id: ""  # 绑定小程序的APPID
    secret: ""  # 绑定小程序的sk
    mch-certificate-serial-number: "" # 商户证书序列号
    mch-api-v3-key: "" # 商户APIv3密钥
    pem-path: ""  # 证书文件所在地址
    notify-url: ""  # 支付回调域名
sf:
  partner-id: ""
  check-code: ""
  sand-check-code: ""
  is-sandbox: ""
```

### 试衣与智能美肤（试衣模型列表）

智能美肤作为独立模型配置在 `tryon_models` 中，与试衣模型共用同一个列表，但按用途区分（`modelUsage`）。

管理端入口：`系统参数 -> 试衣模型列表`。

#### 关键参数说明

| 参数键 | 说明 | 备注 |
| --- | --- | --- |
| `tryon_models` | 试衣模型列表（JSON 数组） | 同时维护试衣模型和智能美肤模型 |

#### 模型字段约定

| 字段 | 说明 |
| --- | --- |
| `modelUsage` | 模型用途：`tryon`（试衣）或 `beautify`（智能美肤） |
| `beautifyModelKey` | 仅 `tryon` 模型使用：绑定独立美肤模型 key，留空则自动选第一个可用美肤模型 |
| `supportsBeautify` | 兼容旧配置字段；建议新配置使用独立 `beautify` 模型 |
| `beautifyExtraCost` | 每次智能美肤额外消耗试衣币，0 表示免费 |
| `beautifyModel` | 美肤模型/动作名（阿里可填 `RetouchSkin`，其他 provider 可填自定义模型名） |
| `beautifyRetouchDegree` | 磨皮强度，0-100 |
| `beautifyWhiteningDegree` | 美白强度，0-100 |
| `beautifyUrl` | 美肤接口地址，阿里场景可留空默认 `https://facebody.cn-shanghai.aliyuncs.com/` |
| `beautifyAccessKeyId` | 阿里场景 AccessKeyId（可选） |
| `beautifyAccessKeySecret` | 阿里场景 AccessKeySecret（可选） |
| `beautifySecurityToken` | 阿里 STS SecurityToken（可选） |
| `beautifyToken` | 美肤 token（阿里可用 AKSK 合并串，其他 provider 可作 Bearer Token） |

#### 智能美肤模型选择顺序

后端在调用 `applyTryonBeautify` 时按以下顺序选择模型：

1. 请求体显式传入的 `beautifyModelKey`
2. 当前试衣任务对应试衣模型上的 `beautifyModelKey`
3. `tryon_models` 中第一个可用的 `modelUsage=beautify` 模型（优先同场景）
4. 兼容旧配置：回退到当前试衣模型项内的美肤字段

#### 智能美肤 provider 调用规则

1. `mode=mock_success`：直接返回原试衣结果图（联调用）
2. 阿里美肤模型：走阿里 RPC 签名调用
3. 非阿里美肤模型：走通用 HTTP `POST` 调用（Bearer Token 可选）

#### 智能美肤接口说明

- 接口：`POST /tryonTask/applyTryonBeautify`
- 鉴权：需要登录（`ApiKeyAuth`）
- 业务约束：
  - 每个试衣任务仅可执行一次智能美肤
  - 仅任务状态为成功且已存在试衣结果图时可发起
  - 扣费取模型 `beautifyExtraCost`（失败会自动退回）

请求示例：

```json
{
  "taskID": 123,
  "retouchDegree": 70,
  "whiteningDegree": 30
}
```

响应体中的 `data.task` 会返回美肤状态字段：`beautifyStatus`、`beautifyTaskNo`、`beautifyResult`、`beautifyCost`、`beautifyRefund`、`beautifyError`、`beautifyAt`。

常见状态：

- `beautifyStatus=processing`：美肤处理中，可轮询 `GET /tryonTask/findTryonTask?ID=任务ID`
- `beautifyStatus=success`：美肤成功，`beautifyResult` 为结果图
- `beautifyStatus=failed`：美肤失败，`beautifyError` 为失败原因，若有扣费将写入退款

如需更新 Swagger 文档，请在 `server` 目录执行：

```shell
swag init
```

## 后端
所在文件位置 `/server`

```shell
# 安装依赖
go mod tidy

#启动项目
go run
```

## 使用展示



