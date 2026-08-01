# 文档结构

本文档在 `docs` 目录之下，因此子目录和文档文件相对于此目录，不再从项目根目录表达路径。


## 目录设计

项目文档分为四个层级：2个核心层级由用户设计，2个辅助层级由 AI 生成：

| 层级 | 目录 | 作者 | 说明 |
|------|------|------|------|
| Conception（构想层） | `conception/` | 人工编写（或人类主导） | 设计构想，作者对协议、系统和应用边界的原始设计。 |
| Decision（决策层） | `decision/` | AI 生成 + 人工审阅 | 架构决策，记录 Conception 尚未明确或澄清模糊的补充决策。 |
| Proposal（提案层） | `proposal/` | AI 生成 | 详细技术规格，追溯自 Conception + Decision。 |
| Plan（方案层） | `plan/` | AI 生成 | 按阶段的实施计划（TDD 任务、包边界、文件清单），追溯自 Proposal。 |

### 权重关系

**权威顺序**：`Conception` + `Decision` > `Proposal` > `Plan`。

如遇冲突以更上层为准，最终以 `Conception` 和 `Decision` 为准。`Decision` 可以是 `Conception` 的细化或模糊 => 明确，但两者不得冲突，否则必须人工介入裁决。

若发现 `Proposal` 或 `Plan` 与 `Conception` 或 `Decision` 不一致，先修改受影响的 `Conception` 或 `Decision` 文档，再重新生成对应的 `Proposal` 和 `Plan` 文件；不要只修改一个文件而不同步其下游内容。

### 排除目录

`plans/`（带 s）用于 AI Agent 工作过程中的临时实施计划，不作为正式文档的一部分；正式方案在 `plan/`。

另外，上级目录中的 `working/` 属于工作间目录，存放临时杂物，也不在正式文档序列里。


## 维护总则（Agent）

以上章节和本章节不可修改，这是文档结构的基础设计和关系逻辑。

下面四个章节记录了项目的设计构想、决策、技术提案和实施计划部分，需要根据实际情况更新。


## 设计构想（Conception）

位于 `conception/` 目录下，包含以下内容：

|    功能    | 设计构想文件 |  简要说明  |
|------------|--------------|------------|
| 核心设计 | [design.md](conception/design.md) | 基网架构、节点信息（分层/PoW/身份绑定）、DHT 结构化查询、自由搜寻微广播、应用池稀疏性策略、安全风险与项目边界 |
| 设计评估 | [evaluate.md](conception/evaluate.md) | 架构核心理念（统一服务发现、子网间服务关系）、百万级规模要点评估、愿景 |
| VRF 周期重洗 | [VRF周期重洗.md](conception/VRF周期重洗.md) | 存储位置 Epoch 重洗（drand 信标）、24 批平滑转储、VRF 资格抽签、时序与参数 |
| 工作量认证 | [工作量认证.md](conception/工作量认证.md) | TLS 连接层 Equi-X PoW（wTLS）、CH.random 首包认证、响应分级与首分享机制 |
| 查询请求调度 | [功能性封面查询调度](功能性封面查询调度.md) | 利用日常维护查询与生态探测使查询流多样化，降低审查系统对单一目标的实时聚焦能力 |


## 架构决策（Decision）

位于 `decision/` 目录。Decision 的作用是补充 Conception 尚未直接固定的规范化细节，例如字节编码、极端边界、字段宽度和实现路径等。

|    编号    |    文件    |  覆盖主题  |
|------------|------------|------------|
| （待更新） | （待更新） | （待更新） |

维护规则：

- 新增 Decision 前必须先检查 `conception/` 是否已经明确该规则。
- 若 Conception 已明确，直接引用 Conception，不新增 Decision。
- 若后续 Conception 修订吸收了某个 Decision，应删除或标记该 Decision 已被吸收。
- Decision 文件命名为 `DEC-<NNNN>-<short-description>.md`。其中 `<NNNN>` 为序号，可以有类别区分（如 `0201`，第2类第1项），也可没有类别（如 `0001`）。


## 技术提案（Proposal）

位于 `proposal/` 目录，由 Conception + Decision 重新生成的可实施技术规格。
每篇含「来源追溯」「规格正文」「边界与限制」「待决问题」「对 Plan 的约束」。

|    编号    |    文件    |  覆盖主题  |
|------------|------------|------------|
| （待更新） | （待更新） | （待更新） |

维护规则：

- Proposal 的权威性低于 Conception 与 Decision。
- 每篇「来源追溯」必须可回溯到具体 Conception 章节与 `DEC-<NNNN>`。
- 待决项严格限于全局待决集，相关规格须显式标注，不得默认选值固化。


## 实施方案（Plan）

位于 `plan/` 目录，由 Proposal 转化的阶段化实施计划。与 Proposal 章节、代码包、实施阶段对齐。
每篇含「来源提案」「包边界」「建议文件」「TDD Task」「阶段门禁/验收」。

|    文件    | 覆盖 Proposal | 对应包（层） |
|------------|---------------|--------------|
| （待更新） |  （待更新）   |  （待更新）  |

维护规则：

- 决策引用为 `DEC-<NNNN>` 且主题匹配。
- 待决项对应的 Task 显式标注阻塞/占位。
- 实现任何功能前，应先读对应 `plan/` 文件，再回溯 `proposal/`，如有疑问查 `decision/` 与 `conception/`。
