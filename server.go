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
	acceptTimeout = 3 * time.Second
)

// ----- Server Public API ---------------------------------------------------

type Server interface {
	Addr() string
	Start() error
	Stop()
	Fd() int
}

func NewServer(port int) Server {
	return &server{
		addr:          fmt.Sprintf(":%v", port),
		handlers:      new(sync.WaitGroup),
		stopAccepting: make(chan int),
	}
}

func (s *server) Addr() string {
	return s.addr
}

func (s *server) Start() error {
	if err := s.listen(); err != nil {
		return err
	}
	return s.serve()
}

func (s *server) Stop() {
	s.stopAccepting <- 1
}

func (s *server) Fd() int {
	return int(s.listenerFile.Fd())
}

// ----- Server Internal API -------------------------------------------------

type server struct {
	addr          string
	listener      net.TCPListener
	listenerFile  *os.File
	handlers      *sync.WaitGroup
	stopAccepting chan int
}

func (s *server) listen() error {
	var addr *net.TCPAddr
	var err error
	var listener *net.TCPListener

	if addr, err = net.ResolveTCPAddr("tcp", s.addr); err != nil {
		return err
	}
	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		return err
	}
	s.listener = *listener
	if s.listenerFile, err = listener.File(); err != nil {
		return err
	}
	if e := syscall.SetNonblock(s.Fd(), true); e != nil {
		return e
	}
	return nil
}

func (s *server) serve() (e error) {
	var accept = true

	for accept {
		var c net.Conn

		s.listener.SetDeadline(time.Now().Add(acceptTimeout))
		if c, e = s.listener.Accept(); e != nil {
			s.logAcceptError(e)
			continue
		}
		go s.handle(c)
		select {
		case <-s.stopAccepting:
			accept = false
		default:
		}
	}
	fmt.Printf("shutting down; waiting for handlers\n")
	s.handlers.Wait()
	return nil
}

func (s *server) handle(c net.Conn) {
	defer s.connTerminated(c)

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

func (s *server) connTerminated(c net.Conn) {
	c.Close()
	s.handlers.Done()
}

func (s *server) logAcceptError(e error) {
	if ope, ok := e.(*net.OpError); ok {
		if !(ope.Timeout() && ope.Temporary()) {
			fmt.Printf("error during accept: %v\n", ope)
		}
	} else {
		fmt.Printf("error during accept: %v\n", e)
	}
}
