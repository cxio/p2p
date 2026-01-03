# 设计优化建议

> **核心目标不变：** 百万级规模的完全去中心化P2P节点发现网络，彻底消除硬植入初始节点信息的中心化弊端。

## 架构核心理念

### 统一的服务发现机制

**关键认识：所有功能都是应用子网，所有子网都通过在基网中搜寻获得初始连接信息。所有子网都集成对基网的连接。**

```
节点自由市场
    ├─ 文件分享应用子网
    ├─ 即时通讯应用子网
    ├─ STUN服务子网 ─────┐
    ├─ 视频直播应用子网    │
    ├─ 区块链应用子网      │
    └─ ...               │
                          │
        ┌─────────────┘
        │ 需要NAT穿透？
        │ 搜寻 "stun-service"
        │ 连入STUN子网
        └─ 获得NAT探测和打洞协助
```

**设计优势：**

1. **完全对等**：STUN不是特殊设计，而是普通的应用子网（服务类，仅用途上划分，不特殊）。
2. **按需连接**：只有需要NAT服务的应用才搜寻并连入STUN子网。
3. **服务链**：应用子网间可以互为服务关系。
   - STUN子网为其它子网提供NAT服务。
   - CDN子网为其它子网提供加速服务。
   - 存储子网为其它子网提供持久化服务。
4. **自然扩展**：任何新服务只需定义名识，就能被发现和使用。
5. **去中心化纯粹**：基网就是一个宽泛的基础设施，所有应用都是子网，所有子网都内嵌连入基网。


### 应用子网间的服务关系

子网为第三方自由独立开发，虽然连入基网获取节点发现功能，但逻辑上是完全独立的P2P网络。

子网中对基网的连接，如果它们不断开连接的话，实际上也就支持并扩展了基网的规模。这是百万级节点连接规模的基础性支撑。

各个服务子网或应用子网间可以互为服务关系，但依然各自独立。依赖是业务上的，不是物理连接上的。每个子网都可以选择性地连入其它子网以获得所需的服务。

也即：所有的子网在连接上都是平等的（各自独立），只是在用途上有层级关系。



## 总体评估

在**百万级规模**和**抗审查**的背景下，当前设计高度合理：

- ✅ 暴力发现：百万节点下期望发现时间 < 1分钟（并发1000）。
- ✅ 广播搜寻：去中心化、流量随机、负载可控。
- ✅ 毒IP策略：提高审查成本的有效手段。
- ✅ 双层网络：基网+任意应用子网，规模效应的正反馈循环。
- ✅ 统一服务发现：所有功能都通过同一机制发现并连接（注：子网内有自己的节点分享扩散能力）。

以下优化建议旨在**提升效率**而非改变核心架构。



## 1. 暴力发现优化

### 1.1 智能扫描策略

```go
// 配置项建议
discovery: {
    strategy: "smart",          // smart | random | sequential
    concurrent: 1000,           // 并发探测数
    timeout: 200,               // 单次探测超时(ms)

    // 智能策略参数
    prefer_cloud: true,         // 优先扫描云服务商IP段
    prefer_history: true,       // 优先扫描历史成功IP附近
    history_radius: 256,        // 历史IP扫描半径(/24网段)

    // IP段优先级配置
    priority_ranges: [
        "1.0.0.0/8",            // 中国电信
        "58.0.0.0/8",           // 中国移动
        "202.0.0.0/8",          // 亚太地区
        // ... 更多云服务商和主要ISP
    ],

    // 避免扫描的IP段
    exclude_ranges: [
        "10.0.0.0/8",           // 私有网络
        "127.0.0.0/8",          // 本地回环
        "224.0.0.0/4",          // 多播
        // ... 更多保留地址
    ]
}
```

**效率提升：**
- 云服务IP集中在特定段，命中率可提升 5-10倍。
- 历史成功IP附近通常有更多节点（同ISP、同地区）。
- 避免无意义的保留地址扫描。


### 1.2 渐进式发现

```go
// 阶段式发现策略
Phase 1: 快速发现（前30秒）
    - 并发: 1000
    - 目标: 发现第1个节点
    - 策略: 高优先级IP段 + 历史记录

Phase 2: 连接建立（30秒-2分钟）
    - 尝试连接发现的节点
    - 如失败，继续Phase 1

Phase 3: 节点获取（2-5分钟）
    - 从成功连接的节点获取上线协助（可连接节点分享）
    - 停止暴力发现

Fallback: 如5分钟仍未成功
    - 扩大到全IPv4扫描
    - 降低并发（避免ISP限制）
    - 或提示用户手动导入节点
```

> **提示：**
> 上线协助是一个专项功能，由在线节点提供*肯定有实际可连接*的节点的列表。
> 列表通常包含节点连接池里的 `Public/FullC` 成员（保证可连接），以及候选池中的部分节点（可能含毒IP）。


### 1.3 发现模式配置

```hjson
// config.hjson 新增配置
{
    // ... 现有配置

    // 节点发现配置
    discovery_mode: "auto",     // auto | aggressive | passive | manual
    discovery_timeout: 300,     // 发现超时时间(秒)，0表示无限制
    discovery_concurrent: 1000, // 并发探测数
    discovery_cloud_first: true,// 优先云服务商IP段
    discovery_history: true,    // 使用历史记录辅助

    // 历史节点缓存
    history_cache: true,        // 启用历史节点缓存
    history_file: "~/.ppnet/history.json",
    history_max: 1000,          // 最多缓存节点数
    history_ttl: 2592000,       // 历史记录有效期(秒，30天)
}
```



## 2. 广播搜寻优化

### 2.1 跳数机制详解

**跳数的工作原理：**

```
询问者设置：
- MaxHops: 最大跳数上限（如15）
- StartHops: 起跳值（可选，默认0）

询问包结构：
参见 docs/design.md 中的“搜寻请求”章节。

转播流程：
询问者(跳数=StartHops)
  → 节点A(跳数+1)
  → 节点B(跳数+1)
  → ...
  → 达到MaxHops停止
```

**关键技巧：起跳值控制广播范围**

```go
// 询问者通过起跳值控制广播范围
query := &Query{
    CurrentHops: 3,    // 从3开始
    MaxHops:     15,   // 硬编码上限（在App中固定）
    TargetID:    "my-app",
}

// 工作原理：
// - 每个转播者简单地对 CurrentHops +1
// - 当 CurrentHops >= MaxHops 时停止转播
// - 起跳值越大，广播范围越小（更快达到上限）
//
// 广播范围控制：
// - StartHops = 0  → 可转播15次，广播范围最大
// - StartHops = 5  → 可转播10次，广播范围中等
// - StartHops = 10 → 可转播5次，广播范围较小
// - StartHops = 14 → 可转播1次，仅直连节点
//
// 用途：
// - 小起跳值：需要大范围搜寻（找稀有应用）
// - 大起跳值：仅需局部搜寻（常见应用或测试）
//
// 注意：
// - 转播者无法验证起跳值的真实性
// - 转播者不知道当前跳数是否真实
// - 转播者只需简单 +1 和检查上限
```

**自适应策略：**

```go
// 询问者角度：根据网络规模调整起跳值
// 实际转播次数 = 15 - 返回值
func (n *Node) SelectStartHops(requirement string) int {
    const MaxHops = 15  // 硬编码上限

    switch requirement {
    case "wide":
        return 0  // 最大广播范围（可转播15次）
    case "medium":
        return 5  // 中等范围（可转播10次）
    case "narrow":
        return 10 // 小范围（可转播5次）
    case "direct":
        return 14 // 仅直连（可转播1次）
    default:
        return 0  // 默认最大范围
    }
}

// 转播者角度：根据自身状态决定是否转播
func (n *Node) ShouldForward(query *Query) bool {
    // 1. 检查跳数是否达到硬编码上限（如15）
    // 这是唯一需要检查跳数的地方
    if query.CurrentHops >= MaxHopsHardLimit {  // 硬编码，如15
        return false
    }

    // 2. 随机概率（避免确定性）
    if rand.Float64() < n.config.ForwardProbability {
        return true
    }

    return false
}

// 转播者角度：决定转播数量
// 注意：不依赖当前跳数！因为跳数可能被篡改或起跳值不为0
func (n *Node) CalculateForwardCount() int {
    // 完全基于自身状态，不考虑跳数
    baseCount := n.config.MaxForward  // 默认8

    num := len(n.GetConnectedPeers())
    if num > baseCount {
        num = baseCount
    }
    return num
}
```

**网络规模估算：**
```go
// 通过连接节点的连接池（组网池）大小估算
func EstimateNetworkSize(peers []Peer) int {
    avgPoolSize := calculateAverage(peers.PoolSize)
    avgConnections := calculateAverage(peers.ConnectionCount)

    // 简化估算：连接池大小 × 平均连接数 × 扩散系数
    estimate := avgPoolSize * avgConnections * 10
    return estimate
}
```


### 2.2 智能转发策略

```go
// 不是简单的全量转发，而是选择性转发
type ForwardStrategy struct {
    MaxForward      int     // 最多转发给几个节点(默认8)
    LoadThreshold   float64 // 负载阈值，超过则减少转发
    PreferNew       bool    // 优先转发给新连接的节点
    PreferQuality   bool    // 优先转发给高质量节点
}

// 转发决策
func (n *Node) ForwardQuery(query *Query) []Peer {
    candidates := n.GetConnectedPeers()

    // 排除源节点和已转发节点
    candidates = excludeVisited(candidates, query.Path)

    // 根据策略排序
    if n.config.PreferQuality {
        sortByQuality(candidates)
    } else if n.config.PreferNew {
        sortByConnectTime(candidates)
    } else {
        shuffle(candidates)  // 随机打乱
    }

    // 根据当前负载决定转发数量
    forwardCount := n.calculateForwardCount()
    return candidates[:min(forwardCount, len(candidates))]
}
```


### 2.3 回传优化

```go
// 当前策略：每级最多8个
// 优化：基于随机性和多样性的回传策略
//
// 重要：不依赖跳数！
// 原因：跳数可能被恶意节点篡改，用于物理分割封锁
// 策略：地理随机分布优于低延迟，因为只需一个有效节点即可连入应用子网

type ReplyPriority struct {
    NodeID      string   // 节点ID
    Quality     float64  // 节点质量评分（基本可靠性）
    Timestamp   time.Time
    IPPrefix    string   // IP段前缀（用于多样性检测）
}

func (n *Node) SelectReplies(replies []Reply) []Reply {
    maxReplies := n.config.MaxRepliesPerHop // 默认8

    // 如果回复较少，全部转发
    if len(replies) <= maxReplies {
        return replies
    }

    // 策略1：优先保证地理/网络多样性（60%）
    diverseReplies := selectDiverseReplies(replies, maxReplies*6/10)

    // 策略2：随机选择剩余名额（40%）
    // 完全随机，避免任何可预测的模式
    remainingCount := maxReplies - len(diverseReplies)
    randomReplies := selectRandomReplies(
        excludeSelected(replies, diverseReplies),
        remainingCount,
    )

    // 合并并随机打乱（避免顺序泄露选择策略）
    selected := append(diverseReplies, randomReplies...)
    shuffle(selected)

    return selected
}

// 选择多样性节点：不同IP段、不同地区
func selectDiverseReplies(replies []Reply, count int) []Reply {
    if count <= 0 || len(replies) == 0 {
        return nil
    }

    selected := make([]Reply, 0, count)
    usedPrefixes := make(map[string]bool) // /16网段去重

    // 随机打乱后遍历，选择不同IP段的节点
    shuffled := shuffle(replies)
    for _, reply := range shuffled {
        prefix := getIPPrefix(reply.NodeInfo.IP, 16) // /16前缀

        if !usedPrefixes[prefix] {
            selected = append(selected, reply)
            usedPrefixes[prefix] = true

            if len(selected) >= count {
                break
            }
        }
    }

    // 如果不同IP段不足，从剩余中随机补充
    if len(selected) < count {
        remaining := excludeSelected(shuffled, selected)
        needed := count - len(selected)
        selected = append(selected, remaining[:min(needed, len(remaining))]...)
    }

    return selected
}

// 完全随机选择
func selectRandomReplies(replies []Reply, count int) []Reply {
    if count <= 0 || len(replies) == 0 {
        return nil
    }

    shuffled := shuffle(replies)
    return shuffled[:min(count, len(shuffled))]
}

// 获取IP前缀（用于多样性检测）
func getIPPrefix(ip string, bits int) string {
    addr := netip.MustParseAddr(ip)
    prefix := netip.PrefixFrom(addr, bits)
    return prefix.Masked().String()
}
```

**设计原则：**

1. **不信任跳数**：跳数可被篡改，不作为选择依据（注：实际上跳数在回应包中已加密）。
2. **地理多样性优先**：不同IP段的节点，抗物理分割封锁。
3. **随机性**：避免可预测模式，防止针对性攻击。
4. **高延迟可接受**：只需一个有效节点即可，应用子网内会快速分享。
5. **顺序混淆**：选中后随机打乱，隐藏选择策略。

**为什么不在意延迟？**

```
搜寻目标：找到同类应用节点
连入流程：
  1. 收到第一个有效节点信息
  2. 连接该节点，加入应用子网
  3. 应用子网内节点相互分享
  4. 快速获得大量本地优质节点

结果：初始节点延迟高无妨，后续会优化
```


### 2.4 缓存机制

```go
// 搜寻结果缓存，避免短期内重复搜寻
type QueryCache struct {
    cache    map[string]*CacheEntry
    ttl      time.Duration
    maxSize  int
}

type CacheEntry struct {
    TargetID  string
    Results   []NodeInfo
    Timestamp time.Time
    Hits      int  // 命中次数，用于淘汰策略
}

// 配置项
query_cache: {
    enabled: true,
    ttl: 300,           // 缓存有效期(秒)
    max_entries: 100,   // 最多缓存条目
    min_results: 3,     // 至少3个结果才缓存
}
```



## 3. 毒IP策略优化

毒IP并不是不能被检测出来，如果实际去连接它们，肯定没有正确的反馈。因此毒IP的设计目标是**增加审查者的成本**——审查者必须逐个去尝试验证，而这通常会耗费不少的时间和资源。


### 3.1 分级毒IP池

```go
// 不同置信度的毒IP
type PoisonIPLevel int

const (
    PoisonConfirmed    PoisonIPLevel = 1  // 确认的大平台IP
    PoisonLikely       PoisonIPLevel = 2  // 可能的大平台IP
    PoisonExperimental PoisonIPLevel = 3  // 实验性毒IP
)

type PoisonPool struct {
    // 按等级分组
    confirmed    []netip.Addr  // 腾讯云、阿里云、AWS等
    likely       []netip.Addr  // 其它云服务商
    experimental []netip.Addr  // 动态生成的
}

// 混入比例可配置
poison_ip: {
    enabled: true,
    confirmed_ratio: 0.15,      // 15%确认的大平台IP
    likely_ratio: 0.10,         // 10%可能的大平台IP
    experimental_ratio: 0.05,   // 5%实验性IP
    total_ratio: 0.30,          // 总占比30%

    // 毒IP来源
    sources: [
        "builtin",              // 内置列表
        "cloud_providers",      // 从公开的云服务商IP列表获取
        "user_config",          // 用户自定义
    ],
}
```


### 3.2 毒IP合理性检测

```go
// 生成看起来合理的毒IP
func GenerateRealisticPoisonIP() netip.Addr {
    // 1. 从大型云服务商的IP段随机选择
    cloudRanges := []string{
        "47.0.0.0/8",      // 阿里云
        "119.0.0.0/8",     // 腾讯云
        "52.0.0.0/8",      // AWS
        "13.0.0.0/8",      // AWS
        // ...
    }

    // 2. 确保端口合理（常见服务端口）
    commonPorts := []int{443, 8080, 7788, 3000, 8443}

    // 3. 生成完整的"节点信息"
    return &NodeInfo{
        IP:       randomIPFromRanges(cloudRanges),
        Port:     commonPorts[rand.Intn(len(commonPorts))],
        Protocol: "quic",
        NATType:  NAT_Public,  // 毒IP肯定都是Public
        SPKI:     generateRandomSPKI(), // 伪装哈希（必要，不然就被识别排除了）
    }
}

// 避免过于随机的特征
func ValidatePoisonIP(ip netip.Addr) bool {
    // 避免：
    // - 保留地址
    // - 过于集中在某个小段
    // - 明显不真实的组合
    return !isReserved(ip) && !isOverConcentrated(ip)
}
```


### 3.3 动态调整策略

```go
// 根据网络环境动态调整毒IP比例
type PoisonStrategy struct {
    BaseRatio      float64
    CurrentRatio   float64
    AdjustInterval time.Duration
}

func (s *PoisonStrategy) Adjust(metrics *NetworkMetrics) {
    // 如果封锁严重（连接成功率低），提高毒IP比例
    if metrics.ConnectionSuccessRate < 0.3 {
        s.CurrentRatio = min(s.BaseRatio * 1.5, 0.5)
    }

    // 如果网络良好，降低毒IP比例（减少开销）
    if metrics.ConnectionSuccessRate > 0.8 {
        s.CurrentRatio = max(s.BaseRatio * 0.7, 0.1)
    }

    log.Printf("调整毒IP比例: %.2f -> %.2f", s.BaseRatio, s.CurrentRatio)
}
```



## 4. 节点质量评估

**注意**：不同的子网应用的节点评估可能不同，因此评估函数应当设计为可配置和可扩展的。不过系统提供一个默认的多维度评分模型，供大多数应用使用。


### 4.1 多维度评分

```go
type NodeQuality struct {
    Latency         float64  // 延迟评分 (0-1)
    Bandwidth       float64  // 带宽评分 (0-1)
    Reliability     float64  // 可靠性评分 (0-1)
    ShareQuality    float64  // 分享质量评分 (0-1)
    ResponseRate    float64  // 响应率 (0-1)

    // 综合评分
    Overall         float64

    // 评估时间
    LastUpdate      time.Time
    SampleCount     int
}

func CalculateOverallQuality(q *NodeQuality) float64 {
    // 加权平均
    weights := map[string]float64{
        "uptime":       0.25,
        "latency":      0.20,
        "bandwidth":    0.15,
        "reliability":  0.20,
        "share":        0.10,
        "response":     0.10,
    }

    score := q.Latency * weights["latency"] +
             q.Bandwidth * weights["bandwidth"] +
             q.Reliability * weights["reliability"] +
             q.ShareQuality * weights["share"] +
             q.ResponseRate * weights["response"]

    return score
}
```


### 4.2 黑名单策略

```go
// 分级黑名单
type BanLevel int

const (
    BanNone      BanLevel = 0  // 无限制
    BanSoft      BanLevel = 1  // 软封禁（降低优先级）
    BanMedium    BanLevel = 2  // 中度封禁（限时阻止）
    BanHard      BanLevel = 3  // 硬封禁（长期阻止）
    BanPermanent BanLevel = 4  // 永久封禁（注：应尽量少使用）
)

type BanEntry struct {
    IP        netip.Addr
    Level     BanLevel
    Reason    string
    BanTime   time.Time
    Duration  time.Duration
    Violations int  // 违规次数
}

// 触发条件
ban_triggers: {
    // 软封禁
    soft: {
        consecutive_failures: 3,    // 连续失败3次
        duration: "1h",
    },

    // 中度封禁
    medium: {
        consecutive_failures: 5,
        timeout_rate: 0.5,          // 超时率>50%
        duration: "12h",
    },

    // 硬封禁
    hard: {
        malicious_response: true,   // 恶意响应
        fake_node_share: true,      // 分享假节点（纯粹无效的数据，非毒IP）
        duration: "3d",             // 考虑到IP变动，硬封禁时间不宜过长
    },

    // 永久封禁
    // App重启后会自然解除（状态不存储）
    // 注：用户可以在配置中手动添加永久封禁IP，每次启动时加载生效。
    permanent: {
        attack_detected: true,      // 检测到攻击行为
        spam_flooding: true,        // 垃圾信息洪泛
    },
}
```



## 5. 性能监控

### 5.1 关键指标

```go
type NetworkMetrics struct {
    // 网络规模
    TotalNodes          int
    ConnectedNodes      int
    CandidatePoolSize   int

    // 连接质量
    ConnectionSuccessRate   float64
    AverageLatency          time.Duration
    PacketLossRate          float64

    // 搜寻效率
    QuerySuccessRate        float64
    AverageQueryTime        time.Duration
    AverageHops             float64
    AverageRepliesReceived  int

    // 发现效率
    DiscoveryAttempts       int64
    DiscoverySuccesses      int64
    DiscoveryTime           time.Duration

    // 负载状态
    InboundConnections      int
    OutboundConnections     int
    MessagesPerSecond       float64
    BandwidthUsage          int64
}

// 性能报告
func (m *NetworkMetrics) GenerateReport() string {
    return fmt.Sprintf(`
P2P Network Metrics Report
==========================
网络规模: %d 节点 (已连接: %d)
连接成功率: %.2f%%
平均延迟: %v
搜寻成功率: %.2f%% (平均 %v, %.1f 跳)
平均收到回复数: %d
发现成功率: %.2f%% (耗时: %v)
当前负载: %.0f msg/s, %s 带宽
`,
        m.TotalNodes,
        m.ConnectedNodes,
        m.ConnectionSuccessRate * 100,
        m.AverageLatency,
        m.QuerySuccessRate * 100,
        m.AverageQueryTime,
        m.AverageHops,
        m.AverageRepliesReceived,
        float64(m.DiscoverySuccesses) / float64(m.DiscoveryAttempts) * 100,
        m.DiscoveryTime,
        m.MessagesPerSecond,
        formatBandwidth(m.BandwidthUsage),
    )
}
```

节点规模评估可以通过连接池大小和平均连接数（单位时间内）进行估算。

> **注意：**
> 连接池大小是动态变化的，仅作为粗略估算依据。


### 5.2 自适应调优

1. 调整连接池大小：
    - 连接成功率低，扩大连接池。
    - 连接成功率高，缩小连接池。
2. 调整转播策略：
    - 收到回复太多，说明网络很活跃或负载高，降低转播概率或减少转播数量。
    - 收到回复太少，说明网络稀疏，提高转播概率或增加转播数量。
3. 调整心跳间隔：
    - 延迟很低，可以降低心跳频率。
    - 丢包率高，提高心跳频率。



## 5. 服务生态示例

场景：用户A想和用户B进行视频通话。

```
1. 用户A启动视频通话应用
   └→ 连入基网（ppnet.Open）
   └→ 搜寻同类应用节点
   └→ 创建应用子网并连入

2. 用户A是NAT内网用户，需要NAT服务
   └→ 在基网搜寻 "stun-service"
   └→ 连入STUN子网
   └→ 探测自己的NAT类型：Sym

3. 用户B启动视频通话应用
   └→ 连入基网（ppnet.Open）
   └→ 搜寻同类应用节点
   └→ 创建应用子网并连入

4. 用户B也是NAT内网用户，需要NAT服务
   └→ 在基网搜寻 "stun-service"
   └→ 连入STUN子网
   └→ 探测自己的NAT类型：Sym

5. 用户A和用户B在视频子网内互相发现，但无法直连
   └→ 两者需要Relay服务
   └→ 用户A在基网搜寻 "relay-service"
   └→ 连入Relay子网
   └→ 向Relay节点提供一个自己的标识
   └→ 用户B也在基网搜寻 "relay-service"
   └→ 连入Relay子网
   └→ 在Relay子网中搜寻用户A（按其标识）

6. 用户B在Relay子网内找到用户A所连接的Relay节点
   └→ 请求该Relay节点协助中继
   └→ Relay节点同意协助
   └→ 用户A ←→ Relay节点 ←→ 用户B
   └→ 视频流通过中继传输
```

> **优化：**
> 用户A和用户B可以同时和多个Relay节点连接，最后选择延迟低的路径进行中继。


## 6. 总结与展望

当前设计将在**百万级规模节点连接**环境下运行。

**核心优势：**
1. ✅ **完全去中心化**：无硬编码节点，无单点故障。
2. ✅ **统一服务发现**：所有功能通过搜寻机制连接。
3. ✅ **规模效应**：应用越多，基网越大，发现越容易。
4. ✅ **抗审查设计**：暴力发现 + 毒IP + 流量随机化。
5. ✅ **服务生态**：应用子网间可互为服务关系。


### 愿景

如果达到预期规模，这将成为：

**真正去中心化的P2P基础设施**
- 无需任何中心化服务。
- 无需预置节点信息。
- 无单点故障风险。
- 抗审查能力强。

**P2P应用的公共市场**
- 任何P2P应用都能轻松接入。
- 应用间共享网络规模优势。
- 服务发现机制统一。
- 降低P2P应用开发门槛。

**繁荣的服务生态**
- 基础服务（STUN、Relay）。
- 平台服务（CDN、Storage、Message）。
- 应用服务（各种P2P应用）。
- 服务间可组合、可竞争。

**抗审查的网络基石**
- 暴力发现突破封锁。
- 毒IP提高审查成本。
- 百万级规模难以整体封锁。
- 为信息自由做点贡献。

---

**这不仅是一个P2P网络项目，更是一个关于去中心化未来的实验。**
