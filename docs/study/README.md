# mini-nacos 学习文档

这套文档用于把 `Nacos` 的核心原理，映射到当前 `Kratos` 版 `mini-nacos` 实现中。

## 文档体系

- `docs/playbook/ai-open-source-learning-playbook.md`
  - 通用学习协作手册
- `docs/study/00-roadmap.md`
  - 学习路线、阶段边界、当前顺序
- `docs/study/01-config.md`
  - `Config` 总览、关键抽象、当前实现边界
- `docs/study/01-config-publish.md`
  - `publish` 主链精读与 Kratos 映射
- `docs/study/01-config-get-listen.md`
  - `get / listen` 主链精读、长轮询语义、当前映射
- `docs/reviews/01-config-get-review.md`
  - `get` 学习过程答题与点评记录
- `docs/design/01-config-round-1.md`
  - Kratos 版 `Config` round-1 设计
- `docs/design/01-config-round-2.md`
  - Kratos 版 `Config get / listen` round-2 设计
- `docs/homework/01-config-round-1.md`
  - 当前轮作业与验收清单
- `docs/homework/01-config-round-2.md`
  - `get / listen` 当前轮作业与验收清单

## 当前原则

- `study` 负责讲清源码链路和设计动机
- `reviews` 负责保留学习过程，不把长问答直接堆进 `study`
- `design` 负责把上游语义映射成当前仓库结构
- `homework` 负责把设计拆成可执行任务
- 三类文档不要重复抄同一段长解释

## 当前进度

- 当前主题：`01-config`
- 当前已完成：`publish-first` 与 `get/listen` 的 Kratos 主链接通、文档与测试
- 当前阶段结论：`01-config` 的 MVP 学习闭环已经完成
- 当前下一步：进入 `02-naming round-1`，先学 `注册 -> 心跳 -> 健康检查 -> 服务发现`
