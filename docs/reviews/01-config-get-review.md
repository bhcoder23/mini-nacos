# Config Get 学习评审记录

## 1. 本轮目标

- 吃透 `Nacos Config get/query` 主链
- 理解 `HTTP / gRPC` 两个入口如何收敛到同一个查询内核
- 理解 `ConfigQueryChainService` 和默认责任链的职责划分
- 为后续 `listen` 学习做铺垫

## 2. 阅读范围

入口层：

- `config/src/main/java/com/alibaba/nacos/config/server/controller/v3/ConfigOpenApiController.java`
- `config/src/main/java/com/alibaba/nacos/config/server/remote/ConfigQueryRequestHandler.java`

查询链：

- `config/src/main/java/com/alibaba/nacos/config/server/service/query/ConfigQueryChainService.java`
- `config/src/main/java/com/alibaba/nacos/config/server/service/query/DefaultConfigQueryHandlerChainBuilder.java`
- `config/src/main/java/com/alibaba/nacos/config/server/service/query/ConfigQueryHandlerChain.java`

关键 handlers：

- `config/src/main/java/com/alibaba/nacos/config/server/service/query/handler/ConfigChainEntryHandler.java`
- `config/src/main/java/com/alibaba/nacos/config/server/service/query/handler/GrayRuleMatchHandler.java`
- `config/src/main/java/com/alibaba/nacos/config/server/service/query/handler/SpecialTagNotFoundHandler.java`
- `config/src/main/java/com/alibaba/nacos/config/server/service/query/handler/FormalHandler.java`
- `config/src/main/java/com/alibaba/nacos/config/server/service/query/handler/ConfigContentTypeHandler.java`

## 3. 我的答案

### 3.1 `get` 双入口的理解

问题：

1. 为什么 `get` 要有 `HTTP` 和 `gRPC` 两个入口，但最后还要收敛到同一个 `ConfigQueryChainService`？
2. 为什么入口层不直接返回磁盘内容，而是先转成 `ConfigQueryChainRequest`？
3. `ConfigOpenApiController` 里为什么要把 `sourceIp` 放进 `appLabels`？
4. 为什么 `v3 HTTP OpenAPI` 只提供 `get`，不提供 `listen`？

我的回答：

- 1. 多协议入口，单业务内核，保证逻辑一致、易维护、易扩展
- 2. 抹平协议差异，构建统一业务请求，支持责任链处理
- 3. 让灰度、权限、路由、限流都能使用客户端 IP
- 4. get 是一次性查询，listen 是推送订阅；HTTP 只做查询，gRPC 专做长连接推送

### 3.2 Query Chain 的理解

问题：

1. 为什么 `ConfigChainEntryHandler` 要先加读锁，再把 `CacheItem` 放进 `ThreadLocal`？
2. 为什么 `GrayRuleMatchHandler` 要放在 `FormalHandler` 前面？
3. `SpecialTagNotFoundHandler` 为什么不继续往下走 `FormalHandler`？
4. 为什么 `ConfigContentTypeHandler` 放在链前面，但实际是“后处理”？
5. 你怎么理解“元数据走 cache，正文走 disk”？它背后的设计动机是什么？

我的回答：

- 1.
- 2.
- 3.
- 4.
- 5.

## 4. AI 点评

### 4.1 `get` 双入口

- 待点评

### 4.2 Query Chain

- 待点评

## 5. 纠偏后的结论

- 待补充

## 6. 追问

- 待补充

## 7. 是否过关

- 当前状态：`进行中`

## 8. 下一轮入口

- `listen` 长轮询入口
- `publish -> get -> listen` 三条链路怎么闭环
