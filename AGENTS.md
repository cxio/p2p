# AGENTS.md

## 项目状态

这是一个**极早期设计阶段**的 Go 项目。`base.go` 是唯一的源文件，大量函数体为占位注释，无测试文件，无外部依赖，无 CI/CD。

## 模块信息

- 模块名：`github.com/cxio/p2p`
- Go 版本：1.25.6
- 无 vendor 目录，无 `go.sum`（当前无外部依赖）

## 开发命令

无 Makefile/Taskfile/CI，直接使用标准 Go 工具链：

```bash
go build ./...        # 构建
go test ./...         # 全量测试
go test -run TestXxx ./...  # 单个测试
go fmt ./...          # 格式化
go vet ./...          # 静态分析
```

## 架构要点

**两种网络实例：**
- `Open(owner Answerer) *Network` — 创建基网（固定名识 `BaseName`，即 `sha1("Hello, P2P..")` 的十六进制）
- `New(name string, answerer Answerer) *Network` — 创建应用子网

**核心接口：**
- `Finder` — 广播查找节点
- `Answerer` — 回应查找请求

**协议栈：**
- 主协议：QUIC/HTTP3，端口 7788
- 降级：HTTP/2，端口 443
- TLS 1.3，后量子加密（X25519MLKEM768），ECH，Protobuf 消息序列化

## 关键约定

- **NAT 类型分级**：`Public > FullC > RC > P_RC > Sym`，只有 `Public` 和 `FullC` 节点才进入候选池（shortlist）并对外分享
- **搜寻跳数**：硬上限 15（4 bit），防止洪泛
- **毒IP策略**：节点分享时混入大平台 IP（默认比例 30%），目的是抗审查
- **工作量证明**：在 TLS ClientHello 中内嵌，依赖 `github.com/cxio/wtls`，哈希算法为 HashX（抗 ASIC）
- **组网池轮换**：每 2~3 分钟随机轮换一个连接，连接存活超过 10 分钟才参与轮换

## 配置文件

`config.jsonc` 为 JSONC 格式（JSON with comments），包含端口、节点池参数、搜寻参数、毒IP策略等运行时配置。
