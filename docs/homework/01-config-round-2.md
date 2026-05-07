# Config Round 2 作业（Kratos 版）

## 1. 学习来源

- 总览：`docs/study/01-config.md`
- `get/listen` 笔记：`docs/study/01-config-get-listen.md`
- 设计文档：`docs/design/01-config-round-2.md`

当前作业覆盖 `get + listen`，不要求你现在同时做历史版本、灰度、集群同步或 `gRPC push`。

## 2. 本轮目标

把你对 `Nacos Config get / listen` 主链的理解，落成当前 `Kratos` 仓库里的第二个可验证闭环。

要求保留这些语义：

- `get` 返回完整最新快照
- `listen` 先比较 `md5`
- `listen` 无变化才等待
- 配置变化时只传播轻量变更信号
- `listen` 不直接返回完整内容
- 客户端收到变化后再调 `get`

## 3. 当前仓库里的落点

- proto contract：`api/configcenter/v1/config_center.proto`
- 薄 service：`internal/service/configcenter.go`
- 核心业务：`internal/biz/config.go`
- 内存实现：`internal/data/configcenter.go`
- server 注册：`internal/server/http.go`、`internal/server/grpc.go`
- 注入：`cmd/mini-nacos/wire.go`

## 4. 子任务清单

### Task 0：源码链路复述

完成标准：

- 你能独立讲清：
  - `ConfigOpenApiController`
  - `ConfigQueryChainService`
  - `ConfigChainEntryHandler`
  - `ConfigServletInner`
  - `LongPollingService`
  - `ConfigCacheService`
  - `DumpService`
  - `RpcConfigChangeNotifier`

### Task 1：补齐 proto contract

文件：

- `api/configcenter/v1/config_center.proto`

完成标准：

- 新增 `GetConfig`
- 新增 `ListenConfig`
- `GetConfig` 走 `GET /v1/configs`
- `ListenConfig` 走 `POST /v1/configs/listen`
- `ListenConfigResponse` 不返回完整配置内容

### Task 2：补齐 biz 模型和接口

文件：

- `internal/biz/config.go`

完成标准：

- 定义 `ListenResult`
- 给 `ConfigWatchHub` 增加 `Wait(...)`
- 保留 `Notify(...)`
- 不把等待集合实现细节写进 `biz`

### Task 3：实现 data 层等待集合

文件：

- `internal/data/configcenter.go`

完成标准：

- 按 `ConfigKey` 维护等待中的监听请求
- 支持超时返回
- 支持命中 key 时唤醒等待者
- 支持返回后清理等待集合
- 不在 `data` 层重复做 `md5` 业务判断

### Task 4：实现 `Get`

文件：

- `internal/biz/config.go`
- `internal/biz/config_test.go`

完成标准：

- `Get` 读取最新快照
- 覆盖“命中 / 不存在”
- `service` 不直接碰 repo

### Task 5：实现 `Listen`

文件：

- `internal/biz/config.go`
- `internal/biz/config_test.go`

完成标准：

- 先查当前快照
- 先比较客户端 `md5`
- 不同则立即返回 `changed = true`
- 相同则进入等待集合
- 超时返回 `changed = false`
- 被通知后返回 `changed = true`

### Task 6：接入 `GetConfig / ListenConfig`

文件：

- `internal/service/configcenter.go`
- `internal/server/http.go`
- `internal/server/grpc.go`
- `cmd/mini-nacos/wire.go`

完成标准：

- `service` 只做 request/response 转换
- 不在 `service` 里算 `md5`
- 不在 `service` 里维护等待集合
- 对外暴露 `GET /v1/configs`
- 对外暴露 `POST /v1/configs/listen`

### Task 7：补齐测试

文件：

- `internal/service/configcenter_test.go`
- `internal/server/http_test.go`

完成标准：

- 覆盖 `GetConfig` 映射
- 覆盖 `ListenConfig` 映射
- 验证真实 HTTP 路由已接通
- `listen` 的 transport 测试只断言“变更信号”，不把完整业务逻辑搬回 `server`

## 5. 验收标准

- 能查询一份配置的完整最新快照
- 能返回当前 `md5`
- 客户端 `md5` 过期时，`listen` 能立即返回
- 客户端 `md5` 一致时，`listen` 能进入等待
- `publish` 后命中的等待者能被唤醒
- 超时后等待者能被清理
- `listen` 不直接返回配置正文
- `service` 保持薄层
- `biz` 保留核心规则

## 6. 提交前自查

- 是否先看了 `docs/study/01-config-get-listen.md`
- 是否能口头讲清 `get / listen` 主链
- 是否区分了“完整快照”和“轻量变更信号”
- 是否先比较 `md5`，再决定要不要等待
- 是否把等待集合实现放在 `data`
- 是否让 `service` 只做 request/response 转换
