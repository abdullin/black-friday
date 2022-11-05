package pipe

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

// from: https://www.sobyte.net/post/2022-03/grpc-in-memory/
var ErrPipeListenerClosed = errors.New(`pipe listener already closed`)

type PipeListener struct {
	ch    chan net.Conn
	close chan struct{}
	done  uint32
	m     sync.Mutex
}

func ListenPipe() *PipeListener {
	return &PipeListener{
		ch:    make(chan net.Conn),
		close: make(chan struct{}),
	}
}

// Accept 等待客户端连接
func (l *PipeListener) Accept() (c net.Conn, e error) {
	select {
	case c = <-l.ch:
		fmt.Printf("accepted connection: '%s'\n", c.RemoteAddr().String())
	case <-l.close:
		e = ErrPipeListenerClosed
	}
	return
}

// Close 关闭 listener.
func (l *PipeListener) Close() (e error) {
	if atomic.LoadUint32(&l.done) == 0 {
		l.m.Lock()
		defer l.m.Unlock()
		if l.done == 0 {
			defer atomic.StoreUint32(&l.done, 1)
			close(l.close)
			return
		}
	}
	e = ErrPipeListenerClosed
	return
}

// Addr 返回 listener 的地址
func (l *PipeListener) Addr() net.Addr {
	return pipeAddr(0)
}
func (l *PipeListener) Dial(network, addr string) (net.Conn, error) {
	return l.DialContext(context.Background(), network, addr)
}
func (l *PipeListener) DialContext(ctx context.Context, network, addr string) (conn net.Conn, e error) {
	// PipeListener是否已经关闭
	if atomic.LoadUint32(&l.done) != 0 {
		e = ErrPipeListenerClosed
		return
	}

	// 创建pipe
	c0, c1 := net.Pipe()
	// 等待连接传递到服务端接收
	select {
	case <-ctx.Done():
		e = ctx.Err()
	case l.ch <- c0:
		conn = c1
	case <-l.close:
		c0.Close()
		c1.Close()
		e = ErrPipeListenerClosed
	}
	return
}

type pipeAddr int

func (pipeAddr) Network() string {
	return `sim`
}
func (pipeAddr) String() string {
	return `sim`
}
