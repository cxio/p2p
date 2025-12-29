package p2p

import (
	"context"
	"errors"
)

// BaseName 基网名识
// 注：匹配测试时也会兼容空字符串。
// sha1("Hello, P2P..")
const BaseName = "dec93aa06b30069759af0920b0d0ca31571d75c5"

// ErrNameMismatch 基网名识不匹配错误。
var ErrNameMismatch = errors.New("base-net's name mismatch")

// NatType NAT类型枚举。
type NatType int

// NAT类型常量。
const (
	Public NatType = iota // 公网（Open Internet）
	FullC                 // 完全锥型（Full Cone）
	RC                    // 受限锥型（Restricted Cone）
	P_RC                  // 端口受限锥型（Port Restricted Cone）
	Sym                   // 对称型（Symmetric）
)

// Peer 节点对象。
type Peer struct {
	Version int      // 版本号（初始版为1）
	Name    string   // 名识
	Proto   string   // 协议（quic | http2）
	Addrs   []string // 地址列表（IP:Port 格式）
	SPKI    []byte   // SPKI 指纹（sha256(SPKI)）
	NATT    NatType  // NAT 类型，可选（仅UDP需要）
	Extra   []byte   // 额外信息，可选
}

// Bytes 序列化节点信息为字节切片。
func (p *Peer) Bytes() []byte {
	// ... 序列化实现
	return nil
}

// Finder 通用广播查找。
// 既用于基网的节点搜寻，也用于子网的特定需求实现。
// 查询依据（msg）的含义由应用自行定义。
// 注记：
// 回应消息会适当加密以保护隐私。但这对用户是透明的。
type Finder interface {
	// Find 搜寻指定键值的消息。
	// @param ctx 上下文控制。
	// @param msg 查找依据。
	// @return 信息提取通道。
	Find(ctx context.Context, msg []byte) <-chan []byte
}

// Answerer 对 Finder 的个体回应。
type Answerer interface {
	// Answer 回答查找请求。
	// @param msg 查找依据。
	// @return 回应数据。
	Answer(msg []byte) ([]byte, error)
}

// Network 网络基础实现
// 包括P2P网络的基本功能：节点管理、分享、连接更新等。
// 另外包含一个通用查找功能。
type Network struct {
	named    string   // 名识，连接时检查，需相同
	answerer Answerer // 对查找的回应

	// 其它字段
}

// New 创建P2P网络实例。
// 主要用于子网嵌入复用，子网名识即为子网应用的名识，
// 在节点连接时，对端必须提供同样的名识才行。
// 注：基网由另外的方法（Open）创建。
// @param name 子网名识。
// @param answerer 查找回应接口实现。
// @return P2P网络实例。
func New(name string, answerer Answerer) *Network {
	return &Network{
		named:    name,
		answerer: answerer,

		// 其它字段初始化
	}
}

// Find 通用查找实现。
// 向连接的节点发送查找请求。
// 从返回的通道中逐个提取回应。
func (nt *Network) Find(ctx context.Context, msg []byte) <-chan []byte {
	ch := make(chan []byte)
	// ... 查找接口实现

	return ch
}

// Peer 随机获取一个联系节点。
// 从当前连接池或候选池中提取（无阻塞）。
// 网络的一个基本方法，主要用于外部取得节点创建连接，获取服务等。
// 注意：返回的节点可能不可用。
func (nt *Network) Peer() *Peer {
	// ... 具体实现
	return nil
}

// Peers 获取多个联系节点。
// 从当前连接池或候选池中提取（无阻塞）。
// @param n 需要的节点数量。
// @return 节点对象切片。
// 注意：返回的节点集未必全部可用。
func (nt *Network) Peers(n int) []*Peer {
	// ... 具体实现
	return nil
}

// Open 打开基网。
// 通用查找实现为自我识别（名识判断）。
// 参考：
// - 基网纯节点仅完成自我识别即可，没有其它功能。
// - 应用节点传递自己的识别实现，即连入基网成为基网兼节点。
// @param owner 自我识别实现。
// @return 基网对象。
func Open(owner Answerer) *Network {
	nt := Network{
		named:    BaseName,
		answerer: owner,

		// 其它内部字段初始化
	}
	return &nt
}
