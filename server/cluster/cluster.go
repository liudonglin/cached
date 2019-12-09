package cluster

import (
	"io/ioutil"
	"time"

	"cached-server/consistent"
	"github.com/hashicorp/memberlist"
)

// Node 接口
type Node interface {
	//用来告诉节点该key是否应该由自己处理
	ShouldProcess(key string) (string, bool)
	Members() []string
	Addr() string
}

// Node接口实现
type node struct {
	*consistent.Consistent
	addr string
}

func (n *node) Addr() string {
	return n.addr
}

func (n *node) ShouldProcess(key string) (string, bool) {
	addr, _ := n.Get(key)
	return addr, addr == n.addr
}

// New ...
func New(addr, cluster string) (Node, error) {
	//memberlist提供了3个类似的函数来生成默认配置项，
	//分别是用于局域网(Local Area Network, LAN)的DefaultLANConfig
	//用于广域网(Wide Area Network, WAN)的DefaultWANConfig
	//以及用于本地回环设备的DefaultLocalConfig
	conf := memberlist.DefaultLANConfig()
	conf.Name = addr
	conf.BindAddr = addr
	//ioutil.Discard 实现了Write方法，任何写入操作都会成功且内容会被直接丢弃，这使得我们的控制终端免于被memberlist的日志刷屏
	conf.LogOutput = ioutil.Discard
	l, e := memberlist.Create(conf)
	if e != nil {
		return nil, e
	}
	if cluster == "" {
		cluster = addr
	}
	clu := []string{cluster}
	_, e = l.Join(clu)
	if e != nil {
		return nil, e
	}

	circle := consistent.New()
	//每个节点的虚拟节点的数量，默认为20。当节点数较少时，20个虚拟节点还不能做到较好的负载均衡，所以将其改为256。
	circle.NumberOfReplicas = 256
	go func() {
		for {
			m := l.Members()
			nodes := make([]string, len(m))
			for i, n := range m {
				nodes[i] = n.Name
			}
			circle.Set(nodes)
			time.Sleep(time.Second)
		}
	}()

	return &node{circle, addr}, nil
}
