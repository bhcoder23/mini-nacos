# 第一部分：Config 总览

## 1. 这一章到底在学什么

这一章不是学“怎么存一段字符串”，而是在学 `Nacos Config` 最核心的三条链路：

- `publish`
- `get/query`
- `listen`

真正要抓住的不是接口数量，而是：

- 事件驱动
- 职责解耦
- 最新快照
- 轻量通知
- 查询与监听分离

## 2. `Config` 的三个核心抽象

### 2.1 配置主键

`(namespace, group, dataId)` 唯一标识一份配置。

### 2.2 最新快照

服务端真正维护的是“当前最新配置状态”，至少包含：

- 正文
- `md5`
- 创建时间
- 更新时间

### 2.3 变更事实

`publish` 后传播的不应该是完整内容，而应该是轻量的“哪份配置变了”：

- 哪个 key
- 新 `md5`
- 变更时间

## 3. 当前 Kratos 仓库里的映射

### `api`

- `api/configcenter/v1/config_center.proto`
  - 对外定义 `PublishConfig`

### `service`

- `internal/service/configcenter.go`
  - 薄 service
  - 负责 proto request -> biz key 的转换

### `biz`

- `internal/biz/config.go`
  - 定义 `ConfigKey / ConfigItem / ConfigChange`
  - 定义 `ConfigRepo / ConfigWatchHub`
  - 实现 `Publish`

### `data`

- `internal/data/configcenter.go`
  - 内存版最新快照存储
  - 最小版通知 hub

### `server + wire`

- `internal/server/http.go`
- `internal/server/grpc.go`
- `cmd/mini-nacos/wire.go`

它们负责把 `ConfigCenterService` 真正暴露出去。

## 4. 当前轮边界

### 当前做什么

- 单机版
- 单进程
- 基于 Kratos proto contract 暴露 `publish`
- 保留轻量变更通知语义

### 当前不做什么

- 集群同步
- beta / tag / gray
- 历史版本
- 完整长轮询等待集合

## 5. 当前结论

到这一阶段，你应该已经能讲清：

- 为什么 `publish` 不是简单 CRUD
- 为什么服务端传播的是变更事实，不是完整内容
- 为什么 `service` 应该薄，核心规则要留在 `biz`

下一步建议：

- 复习 `publish` 链路
- 再进入 `get/listen` 链路
