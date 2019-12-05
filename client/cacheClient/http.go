package cacheclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type httpClient struct {
	*http.Client
	server string
	port   int
}

func (c *httpClient) getUrl(key string) string {
	server, _ := shouldProcess(key)
	println("当前调用服务：" + server)
	return fmt.Sprintf("http://%s:%d/cache/", server, c.port)
}

func (c *httpClient) get(key string) (string, bool) {
	resp, e := c.Get(c.getUrl(key) + key)
	if e != nil {
		log.Println(key)
		panic(e)
	}
	if resp.StatusCode == http.StatusNotFound {
		return "", false
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		panic(e)
	}

	result := &Result{}
	json.Unmarshal(b, result)

	if result.Redirect != nil {
		setNodes(result.Redirect)
		return "", true
	}

	return result.Data, false
}

func (c *httpClient) set(key, value string) bool {
	req, e := http.NewRequest(http.MethodPut,
		c.getUrl(key)+key, strings.NewReader(value))
	if e != nil {
		log.Println(key)
		panic(e)
	}
	resp, e := c.Do(req)
	if e != nil {
		log.Println(key)
		panic(e)
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		panic(e)
	}

	result := &Result{}
	json.Unmarshal(b, result)
	if result.Redirect != nil {
		setNodes(result.Redirect)
		return true
	}

	return false
}

func (c *httpClient) Run(cmd *Cmd) {
	retry := false
	if cmd.Name == "get" {
		cmd.Value, retry = c.get(cmd.Key)
	}
	if cmd.Name == "set" {
		retry = c.set(cmd.Key, cmd.Value)
	}

	if retry {
		c.Run(cmd)
	}
}

func newHTTPClient(server string, port int) *httpClient {
	client := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1}}
	return &httpClient{client, server, port}
}

func (c *httpClient) PipelinedRun([]*Cmd) {
	panic("httpClient pipelined run not implement")
}
