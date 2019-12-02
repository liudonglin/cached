package main

import (
	"flag"
	"log"

	"cached-server/cache"
	"cached-server/cluster"
	"cached-server/http"
	"cached-server/tcp"
)

var typ, node, clus string

func init() {
	flag.StringVar(&typ, "type", "inmemory", "cache type")
	flag.StringVar(&node, "node", "127.0.0.1", "current node address")
	flag.StringVar(&clus, "cluster", "", "cluster address")
	flag.Parse()
	log.Println("type is", typ)
	log.Println("node is", node)
	log.Println("cluster is", clus)
}

func main() {
	c := cache.New(typ)

	n, e := cluster.New(node, clus)
	if e != nil {
		panic(e)
	}

	go tcp.New(c, n).Listen()
	http.New(c, n).Listen()
}
