package server

import (
	"context"
	"fmt"
	"net"
	"sync"
)

type server struct {
	Addr      string        // Address the server will listen on
	Port      string        // Listening Port
	Ln        net.Listener  // Used to listen for and accept connections
	Quitting  bool          // If true server will shutdown
	exit      chan struct{} // used to close down the server
	waitGroup sync.WaitGroup
}

func New(addr, port string) (s *server, err error) {
	s = new(server)
	s.Addr = addr
	s.Port = port
	s.Quitting = false
	s.exit = make(chan struct{})
	s.Ln, err = net.Listen("tcp", s.FullAddr())
	return
}

func (s *server) FullAddr() string {
	return fmt.Sprintf("%s:%s", s.Addr, s.Port)
}

func (s *server) Start() (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	go s.listen(ctx)
	<-s.exit
	cancel()
	return
}

func (s *server) Stop() {
	s.exit <- struct{}{}
	s.waitGroup.Wait()
}

func (s *server) listen(ctx context.Context) (err error) {
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

func (s *server) accept(ctx context.Context) (err error) {
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

func (s *server) close() (err error) {
	return
}
