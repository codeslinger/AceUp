// vim:set ts=2 noet ai ft=go:
// Copyright (c) Toby DiPasquale. See accompanying LICENSE file for
// detailed licensing information.
package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"syscall"
	"time"
)

const (
	acceptTimeout = 3e9
)

type Server struct {
	Addr          string
	listener      net.TCPListener
	listenerFile  *os.File
	handlers      *sync.WaitGroup
	stopAccepting chan int
}

func NewServer(port int) *Server {
	server := new(Server)
	server.Addr = fmt.Sprintf(":%v", port)
	server.handlers = new(sync.WaitGroup)
	server.stopAccepting = make(chan int)
	return server
}

func (server *Server) Start() error {
	if err := server.listen(); err != nil {
		return err
	}
	return server.serve()
}

func (server *Server) Stop() {
	server.stopAccepting <- 1
}

func (server *Server) Fd() int {
	return int(server.listenerFile.Fd())
}

func (server *Server) listen() error {
	var addr *net.TCPAddr
	var err error
	var listener *net.TCPListener

	if addr, err = net.ResolveTCPAddr("tcp", server.Addr); err != nil {
		return err
	}
	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		return err
	}
	server.listener = *listener
	if server.listenerFile, err = listener.File(); err != nil {
		return err
	}
	if e := syscall.SetNonblock(server.Fd(), true); e != nil {
		return e
	}
	return nil
}

func (server *Server) serve() (e error) {
	var accept = true

	for accept {
		var c net.Conn

		server.listener.SetDeadline(time.Now().Add(acceptTimeout))
		if c, e = server.listener.Accept(); e != nil {
			server.logAcceptError(e)
			continue
		}
		go server.handle(c)
		select {
		case <-server.stopAccepting:
			accept = false
		default:
		}
	}
	fmt.Printf("shutting down; waiting for handlers\n")
	server.handlers.Wait()
	return nil
}

func (server *Server) handle(c net.Conn) {
	defer server.connTerminated(c)

	var req *Request
	var err error

	for err == nil {
		if req, err = ReadRequest(c); err != nil {
			if err != io.ErrUnexpectedEOF {
				fmt.Printf("error reading request from %v: %v", c.RemoteAddr(), err)
			}
			continue
		}
		// TODO: do something with request here
		fmt.Printf("%v", req)
	}
}

func (server *Server) connTerminated(c net.Conn) {
	c.Close()
	server.handlers.Done()
}

func (server *Server) logAcceptError(e error) {
	if ope, ok := e.(*net.OpError); ok {
		if !(ope.Timeout() && ope.Temporary()) {
			fmt.Printf("error during accept: %v\n", ope)
		}
	} else {
		fmt.Printf("error during accept: %v\n", e)
	}
}

