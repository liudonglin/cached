package cacheclient

type Cmd struct {
	Name  string
	Key   string
	Value string
	Error error
}

type Client interface {
	Run(*Cmd)
	PipelinedRun([]*Cmd)
}

func New(typ, server string, port int) Client {
	if typ == "redis" {
		return newRedisClient(server, port)
	}
	if typ == "http" {
		return newHTTPClient(server, port)
	}
	if typ == "tcp" {
		return newTCPClient(server, port)
	}
	panic("unknown client type " + typ)
}
