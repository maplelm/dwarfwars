package server

import (
	"context"
	"fmt"
	"net"
	"sync"
)

type ConnectionHandler interface {
	Welcome() error
	Update() error
}

type Server struct {
	Addr      string        // Address the server will listen on
	Port      string        // Listening Port
	Ln        net.Listener  // Used to listen for and accept connections
	Quitting  bool          // If true server will shutdown
	exit      chan struct{} // used to close down the server
	waitGroup sync.WaitGroup
}

func New(addr, port string) (s *Server, err error) {
	s = new(Server)
	s.Addr = addr
	s.Port = port
	s.Quitting = false
	s.exit = make(chan struct{})
	s.Ln, err = net.Listen("tcp", s.FullAddr())
	return
}

func (s *Server) FullAddr() string {
	return fmt.Sprintf("%s:%s", s.Addr, s.Port)
}

func (s *Server) Start() (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	go s.listen(ctx)
	<-s.exit
	cancel()
	return
}

func (s *Server) Stop() {
	s.exit <- struct{}{}
	s.waitGroup.Wait()
}

func (s *Server) listen(ctx context.Context) (err error) {
	s.waitGroup.Add(1)
	defer s.waitGroup.Done()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fmt.Println("listen loop")
		}
	}
}

func (s *Server) accept(ctx context.Context) (err error) {
	s.waitGroup.Add(1)
	defer s.waitGroup.Done()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fmt.Println("accept loop")
		}
	}
}

func (s *Server) close() (err error) {
	return
}
