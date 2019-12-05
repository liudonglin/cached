package main

import "cached-client/cacheclient"

func main() {
	typ, server, port := "http", "10.0.0.27", 6800
	client := cacheclient.New(typ, server, port)

	client.Run(&cacheclient.Cmd{Name: "set", Key: "user", Value: "liudonglin"})

	result := &cacheclient.Cmd{Name: "get", Key: "user"}
	client.Run(result)
	println(result.Value)
}
