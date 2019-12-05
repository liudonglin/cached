package cacheclient

import (
	"stathat.com/c/consistent"
)

var circle *consistent.Consistent

func init() {
	circle = consistent.New()
	//每个节点的虚拟节点的数量，默认为20。当节点数较少时，20个虚拟节点还不能做到较好的负载均衡，所以将其改为256。
	circle.NumberOfReplicas = 256
}

func setNodes(m []string) {
	nodes := make([]string, len(m))
	for i, n := range m {
		nodes[i] = n
	}
	circle.Set(nodes)
}

func shouldProcess(key string) (string, bool) {
	addr, _ := circle.Get(key)
	return addr, addr == addr
}
