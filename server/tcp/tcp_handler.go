package tcp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

// set  S<klen><SP><vlen><SP><key><value>
// 客户端发送的command字节流以一个大写的“S”开始，
// 后面跟了一个数字klen表示key的长度，
// 然后是一个空格<SP>作为分隔符，
// 然后是另一个数字vlen表示value的长度，
// 然后又是一个空格，最后是key的内容和value的内容 。
func (s *Server) set(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, v, e := s.readKeyAndValue(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	go func() {
		c <- &result{nil, s.Set(k, v)}
	}()
}

// get  G<klen><SP><key>
// 客户端发送的command以一个大写的“G"开始，
// 后面跟了一个数字klen表示key的长度，
// 然后是一个空格作为分隔符， 最后是key的内容
func (s *Server) get(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, e := s.readKey(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	go func() {
		v, e := s.Get(k)
		c <- &result{v, e}
	}()
}

// del  D<klen><SP><key>
// 客户端发送的command以一个大写的“D”开始，
// 后面跟了一个数字klen表示key的长度，
// 然后是一个空格作为分隔符，最后是key的内容
func (s *Server) del(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, e := s.readKey(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	go func() {
		c <- &result{nil, s.Del(k)}
	}()
}

// -<len><SP><content>|<len><SP><content>
// response规则用来表示服务端发送给客户端的响应，由一个error或者一个bytes-array组成
// error由一个“一”(负号)和一个bytes-array组成，表示错误
func sendResponse(value []byte, err error, conn net.Conn) error {
	if err != nil {
		errString := err.Error()
		tmp := fmt.Sprintf("-%d ", len(errString)) + errString
		_, e := conn.Write([]byte(tmp))
		return e
	}
	vlen := fmt.Sprintf("%d ", len(value))
	_, e := conn.Write(append([]byte(vlen), value...))
	return e
}

type result struct {
	v []byte
	e error
}

func reply(conn net.Conn, resultCh chan chan *result) {
	defer conn.Close()
	for {
		c, open := <-resultCh
		if !open {
			return
		}
		r := <-c
		e := sendResponse(r.v, r.e, conn)
		if e != nil {
			log.Println("close connection due to error:", e)
			return
		}
	}
}

func (s *Server) process(conn net.Conn) {
	//在conn上套了一层bufio.Reader结构体，用来对客户端连接进行一个缓冲读取 。
	//这是很有必要的，因为来自网络的数据不稳定，在我们进行读取时， 客户端的数据可能只传输了一半 ，
	//我们希望可以阻塞等待 ，直到我们需要的数据全 部就位以后一次性返回给我们 。
	//所以这里我们用bufio.NewReader创建了一个bufio.Reader结构体 。
	//它提供了一些特殊的read功能，如Read.Byte和 ReadString等方法。
	//当我们从 ufio.Reader中读取数据时，实际的数据读取自客户端连接conn，
	//如果现有数据不能满足我们的要求，bufio.Reader会进行阻塞等待，直到数据满足要求了才返回。
	r := bufio.NewReader(conn)

	//管道中的数据是先进先出的
	resultCh := make(chan chan *result, 5000)
	defer close(resultCh)
	//开启响应线程
	go reply(conn, resultCh)
	for {
		op, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Println("close connection due to error:", e)
			}
			return
		}

		switch op {
		case 'S':
			s.set(resultCh, r)
		case 'G':
			s.get(resultCh, r)
		case 'D':
			s.del(resultCh, r)
		default:
			log.Println("close connection due to invalid operation:", op)
			return
		}
	}
}
