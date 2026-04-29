# Config Round 1 作业（Kratos 版）

## 1. 学习来源

- 设计文档：`docs/design/01-config-round-1.md`
- 总览：`docs/study/01-config.md`
- `publish` 笔记：`docs/study/01-config-publish.md`

当前作业先覆盖 `publish-first`，不要求你现在同时完整实现 `get + listen`。

## 2. 本轮目标

把你对 `Nacos Config publish` 主链的理解，先落成当前 `Kratos` 仓库里的第一个可验证闭环。

要求保留这些语义：

- 发布与监听唤醒职责分离
- 通过轻量变更元数据传播变化
- 先落成最新快照，再决定是否通知
- 同内容重复发布不应唤醒监听者

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
  - `ConfigControllerV3`
  - `ConfigOperationService`
  - `ConfigChangePublisher`
  - `DumpService`
  - `DumpProcessor`
  - `DumpConfigHandler`
  - `ConfigCacheService`
  - `LongPollingService`

### Task 1：核心模型

文件：

- `internal/biz/config.go`

完成标准：

- 定义 `ConfigKey`
- 定义 `ConfigItem`
- 定义 `ConfigChange`
- 定义 `ErrConfigNotFound`

### Task 2：repo 与 watch hub 接口

文件：

- `internal/biz/config.go`
- `internal/data/configcenter.go`

完成标准：

- `ConfigRepo.Save/Get`
- `ConfigWatchHub.Notify`
- 只传播 `ConfigChange`

### Task 3：实现 `Publish`

文件：

- `internal/biz/config.go`
- `internal/biz/config_test.go`

完成标准：

- 计算 `md5`
- 查询当前最新快照
- 保存新的 `ConfigItem`
- 只有在 `md5` 真变化时才调用 `Notify`
- 覆盖“首次发布 / 内容变化 / 内容不变”

### Task 4：接入 `PublishConfig`

文件：

- `api/configcenter/v1/config_center.proto`
- `internal/service/configcenter.go`
- `internal/server/http.go`
- `internal/server/grpc.go`
- `cmd/mini-nacos/wire.go`

完成标准：

- `POST /v1/configs`
- service 只负责 proto request/response 转换
- 不在 service 里算 `md5`
- 不在 service 里直接碰 repo / watch hub

### Task 5：transport / service 测试

文件：

- `internal/service/configcenter_test.go`

完成标准：

- 覆盖 `PublishConfig` 映射
- 确认 response 中包含最新快照关键字段

### Task 6：补齐 HTTP publish 验证

文件：

- `internal/server/http_test.go`

完成标准：

- 真实构造当前 Kratos HTTP server
- 发送 `POST /v1/configs`
- 断言路由已经由 generated HTTP binding 正常注册
- 断言 response 至少包含：
  - `namespace`
  - `group`
  - `dataId`
  - `content`
  - `md5`
- 这个 task 只验证 transport 接通，不把业务逻辑搬回 `server`

## 5. 当前状态

- `Task 3`：已完成
- `Task 4`：主链已接通
- `Task 5`：已完成
- `Task 6`：待开始

## 6. 验收标准

- 能唯一标识一份配置
- 能发布配置
- 能形成最新快照
- 能基于内容 `md5` 判断是否真正变化
- 能在 `md5` 变化时传播轻量 `ConfigChange`
- 同内容重复发布不会触发通知
- `service` 保持薄层
- `biz` 保留核心规则

## 7. 提交前自查

- 是否先看了 `docs/study/01-config-publish.md`
- 是否能口头讲清 `publish` 主链
- 是否把“完整内容”和“轻量变更事实”区分开
- 是否把业务逻辑留在 `biz`
- 是否让 `service` 只做 request/response 转换
