
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

## 后端
所在文件位置 `/server`

```shell
# 安装依赖
go mod tidy

#启动项目
go run
```

## 使用展示



