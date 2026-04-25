# AGENTS.md

This file provides guidance to AI Agent when working with code in this repository.

## 项目概述

`github.com/cxio/p2p` 是一个去中心化 P2P 节点发现与连接系统，为其他 P2P 应用提供"基网"（Base Network）服务。Go 模块，当前 Go 版本 1.26.2。

## 常用命令

```bash
# 构建
go build ./...

# 运行测试
go test ./...

# 运行单个测试
go test -run TestXxx ./path/to/package

# 代码格式化与检查
go fmt ./...
go vet ./...
```

> 注意：本项目目前处于初期实现阶段，暂无 Makefile、CI/CD 或外部依赖。

## 架构概览

### 核心分层

```
基网（Base Network）
├── 基网纯节点 —— 仅提供基网服务
├── 基网兼节点 —— 兼跑应用子网
└── 应用子网（Application Subnets）—— 各类 P2P 应用
```

### 主要类型（`base.go`）

| 类型 | 说明 |
|------|------|
| `NatType` | NAT 类型枚举（Public / FullCone / RestrictedCone / PortRestrictedCone / Symmetric） |
| `Peer` | 节点信息（版本、身份名、协议、地址列表、SPKI、ECH Key、NAT 类型、扩展数据） |
| `Finder` | 广播查找接口：`Find(ctx, msg) <-chan []byte` |
| `Answerer` | 响应处理接口：`Answer(msg) ([]byte, error)` |
| `Network` | 核心网络实现，管理连接池、节点分享、查找路由 |

### 节点发现机制（被动广播）

1. 客户端发出携带目标 ID、跳数、公钥的查询报文
2. 中继节点若 ID 不匹配则转发（跳数上限 15）
3. 目标节点以加密方式返回联系信息（IP:Port、SPKI、ECH Key）
4. 响应沿原路返回，使用原始 Query ID
5. 每节点最多返回 8 个包；上游每个下游节点取 2–3 个

### 连接生命周期

1. **明文上线协助**：PoW 明文 ClientHello → 服务端返回 3–10 个节点联系信息 + ECH Keys（一次性）
2. **一级认证**：ECH 加密 ClientHello + PoW → 建立持久连接
3. **拥塞控制**：服务端过载时触发二级 PoW

### 安全机制

- TLS 1.3 + 自签名证书，以 SPKI 指纹替代 CA 验证
- HashX PoW 写入 TLS ClientHello，防止 DDoS
- ECH（Encrypted Client Hello）流量混淆
- 毒化 IP 策略：节点分享中混入知名平台 IP（可能超过50%），提升审查成本

### 配置文件（`config.jsonc`）

关键配置字段：`server_uport`（QUIC/UDP，默认 7788）、`server_tport`（HTTP/2/TCP，默认 443）、连接池大小（shortlist 100、outgoing 10、incoming 50）、毒化 IP 比例（30%）。

## 文档结构

设计文档位于 `docs/` 目录，阅读顺序建议：

1. `docs/conception/design.md` —— 核心架构与协议设计
2. `docs/conception/guidance.md` —— 实现参考指南
3. `docs/conception/工作量认证.md` —— PoW 算法说明
4. `docs/proposal/` —— 初始提案
5. `docs/plan/` —— 实现计划
