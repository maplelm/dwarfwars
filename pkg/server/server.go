package server

import (
	"context"
	"fmt"
	"log"
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
	go func(ctx context.Context) error {
		s.waitGroup.Add(1)
		defer s.waitGroup.Done()
		for {
			conn, err := s.Ln.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					// Server is shutting down. no error to report
					return ctx.Err()
				default:
					log.Printf("Server Listener Failed: %s", err)
				}
				continue
			}
			go s.accept(ctx, conn)
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				log.Printf("Server Acccepted connection: %s", conn.RemoteAddr())
			}
		}
	}(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			s.Ln.Accept()
		}
	}
}

func (s *Server) accept(ctx context.Context, conn net.Conn) (err error) {
	s.waitGroup.Add(1)
	defer s.waitGroup.Done()
	defer conn.Close()
	log.Printf("Server Accepted Connection from %s", conn.RemoteAddr())
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Making the server echo for now
			var req []byte = []byte{}
			rn, err := conn.Read(req)
			if err != nil {
				log.Printf("Conn (%s) Error: %s", conn.RemoteAddr(), err)
			}
			wn, err := conn.Write(req)
			if err != nil {
				log.Printf("Conn (%s_ Error: %s", conn.RemoteAddr(), err)
			}
			if wn != rn {
				log.Printf("Conn (%s) Warning: read and write data lengths do not match R: %d, W: %d", conn.RemoteAddr(), rn, wn)
			}

		}
	}
}
