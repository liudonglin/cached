package main

import (
	"cached-server/cache"
	"cached-server/cluster"
	"cached-server/http"
	"cached-server/tcp"
	"flag"
	"log"
	"net"
)

var typ, node, clus string

func init() {

	ipv4, _ := LocalIPv4s()
	if ipv4 == "" {
		ipv4 = "127.0.0.1"
	}

	flag.StringVar(&typ, "type", "inmemory", "cache type")
	flag.StringVar(&node, "node", ipv4, "current node address")
	flag.StringVar(&clus, "cluster", "", "cluster address")
	flag.Parse()
	log.Println("type is", typ)
	log.Println("node is", node)
	log.Println("cluster is", clus)
}

// LocalIPv4s 获取本机ip地址
func LocalIPv4s() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	ipv4 := ""
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipv4 = ipnet.IP.String()
				break
			}
		}
	}
	return ipv4, nil
}

func main() {
	c := cache.New(typ)

	n, e := cluster.New(node, clus)
	if e != nil {
		panic(e)
	}
	// 启动tcp监听
	go tcp.New(c, n).Listen()
	// 启动http监听
	http.New(c, n).Listen()
}
