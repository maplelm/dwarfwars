package server

/*
* WARNING: CURRENTLY NOT VALIDATING CONNECTIONS OR TRAFFIC.
* WARNING: No Authentication is required.
 *
 * FIX: need to acount for bytes not being sent all at once
 * FIX: Connections currently do not timeout leading to a runaway goruitines
*/

import (
	"context"
	"fmt"
	"github.com/maplelm/dwarfwars/pkg/tcp"
	"log"
	"net"
	"sync"
	"time"
)

const (
	PoolPublic = iota
)

type Server struct {
	Addr        string         // Address the server will listen on
	Port        string         // Listening Port
	Ln          net.Listener   // Used to listen for and accept connections
	Quitting    bool           // If true server will shutdown
	exit        chan struct{}  // used to close down the server
	waitGroup   sync.WaitGroup // tracking goruitines for connections so they can shutdown properly
	idleTimeout time.Duration
	connPool    map[int][]net.Conn
}

func New(addr, port string, timeout time.Duration) (s *Server, err error) {
	s = new(Server)
	s.Addr = addr
	s.Port = port
	s.Quitting = false
	s.exit = make(chan struct{})
	s.Ln, err = net.Listen("tcp", s.FullAddr())
	s.idleTimeout = timeout
	s.connPool = make(map[int][]net.Conn)
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
	log.Printf("Server shutting down, exiting the start command")
	return
}

func (s *Server) Stop() {
	log.Printf("Server.Stop: Called")
	s.exit <- struct{}{}
	s.waitGroup.Wait()
}

func (s *Server) listen(ctx context.Context) (err error) {
	s.waitGroup.Add(1)
	defer s.waitGroup.Done()
	go func(ctx context.Context) error {
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
			s.connPool[PoolPublic] = append(s.connPool[PoolPublic], conn)
			go s.accept(ctx, &s.connPool[PoolPublic][len(s.connPool[PoolPublic])-1])
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
			s.Ln.Close()
			return ctx.Err()
		}
	}
}

func (s *Server) accept(ctx context.Context, conn *net.Conn) (err error) {
	s.waitGroup.Add(1)
	defer s.waitGroup.Done()

	defer (*conn).Close()

	log.Printf("Server Accepted Connection from %s", (*conn).RemoteAddr())

	for {
		select {

		case <-ctx.Done():
			return ctx.Err()

		default:
			// Making the server echo for now
			var cmd tcp.Command
			var data []byte
			rn, err := (*conn).Read(data)
			if err != nil {
				log.Printf("Conn (%s) Error: %s", (*conn).RemoteAddr(), err)
			}
			err = cmd.UnmarshalBinary(data)
			if err != nil {
				log.Printf("Failed to Unmarshal command from client: %s", err)
				continue
			}

			response, err := cmd.MarshalBinary()
			if err != nil {
				log.Printf("Failed to Marshal Binary command for client: %s", err)
				continue
			}

			wn, err := (*conn).Write(response)
			if err != nil {
				log.Printf("Conn (%s) Error: %s", (*conn).RemoteAddr(), err)
			}
			if wn != rn {
				log.Printf("Conn (%s) Warning: read and write data lengths do not match R: %d, W: %d", (*conn).RemoteAddr(), rn, wn)
			}

		}
	}
}
