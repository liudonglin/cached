package tcp

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

// readKey 方法用readLen和io.ReadFull
// 函数来解析客户端发来的command,从中获取key和value
func (s *Server) readKey(r *bufio.Reader) (string, error) {
	klen, e := readLen(r)
	if e != nil {
		return "", e
	}
	k := make([]byte, klen)
	_, e = io.ReadFull(r, k)
	if e != nil {
		return "", e
	}
	//判断是否应由当前节点处理
	key := string(k)
	addr, ok := s.ShouldProcess(key)
	if !ok {
		return "", errors.New("redirect:" + addr)
	}

	return key, nil
}

func (s *Server) readKeyAndValue(r *bufio.Reader) (string, []byte, error) {
	klen, e := readLen(r)
	if e != nil {
		return "", nil, e
	}
	vlen, e := readLen(r)
	if e != nil {
		return "", nil, e
	}
	k := make([]byte, klen)
	_, e = io.ReadFull(r, k)
	if e != nil {
		return "", nil, e
	}
	//判断是否应由当前节点处理
	key := string(k)
	addr, ok := s.ShouldProcess(key)
	if !ok {
		return "", nil, errors.New("redirect:" + addr)
	}

	v := make([]byte, vlen)
	_, e = io.ReadFull(r, v)
	if e != nil {
		return "", nil, e
	}
	return key, v, nil
}

// readLen 读取klen或者vlen
// 函数以空格为分隔符读取一个字符串并将之转化为一个整型。
func readLen(r *bufio.Reader) (int, error) {
	// ReadString读取直到输入中第一次出现delim,
	// 返回一个字符串，该字符串包含分隔符和分隔符之前的数据
	tmp, e := r.ReadString(' ')
	if e != nil {
		return 0, e
	}
	l, e := strconv.Atoi(strings.TrimSpace(tmp))
	if e != nil {
		return 0, e
	}
	return l, nil
}
