# mini-nacos 学习路线图

## 1. 项目目标

我们不是复刻完整 `Nacos`，而是要做一个：

- 用 `Go + Kratos` 实现的 `mini-nacos`
- 聚焦 `Nacos` 的核心机制
- 优先做出可讲清楚、可运行、可扩展的 `MVP`

核心原则：

- 学抽象，不背模块名
- 先做闭环，再做增强
- 先学内核，再碰边缘能力

## 2. 当前技术载体

当前仓库已经切到 `Kratos layout`，关键分层如下：

- `api`
  - 对外 contract，维护 proto、HTTP 注解、生成代码
- `internal/service`
  - 薄 service，负责 request/response 到 biz 的转换
- `internal/biz`
  - 核心业务规则、领域对象、repo 接口
- `internal/data`
  - repo / cache / in-memory 状态实现
- `internal/server`
  - HTTP / gRPC server 注册
- `cmd/mini-nacos`
  - 应用入口与 `wire` 注入

这意味着后续学习要把旧的 `controller/usecase/repo` 映射成：

- `controller` ≈ `service`
- `usecase` ≈ `biz`
- `repo implementation` ≈ `data`

## 3. MVP 范围

第一阶段只抓最重要的内核：

- `Config`
  - 配置发布、查询、监听、长轮询、变更通知
- `Naming`
  - 服务注册、发现、心跳、健康检查

暂不作为第一优先级：

- 控制台 UI
- 认证授权
- 插件系统
- 完整分布式一致性实现

## 4. 当前学习顺序

1. `Config`
2. `Naming`
3. `AP / CP`
4. `AI Resource Layer`

## 5. 当前阶段

### 阶段 A：publish-first

目标：

- 吃透 `Nacos Config publish` 主链
- 在 Kratos 版 `mini-nacos` 中落一个最小可运行闭环

当前已完成：

- `biz.Publish(...)`
- `service.PublishConfig(...)`
- `POST /v1/configs`
- 内存版最新快照和轻量通知路径

### 阶段 B：get + listen

下一步要补：

- 查询完整快照
- 监听请求和等待集合
- “变更信号”和“完整内容”分离

## 6. 每一章学习完成标准

- 你能不用术语堆砌，把原理讲清
- 你能说出关键对象和调用链
- 你能把链路映射到 `api/service/biz/data`
- 你能解释“为什么这样设计”
