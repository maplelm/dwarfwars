package game

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/maplelm/dwarfwars/cmd/client/pkg/types"
	"github.com/maplelm/dwarfwars/pkg/cache"
	"github.com/maplelm/dwarfwars/pkg/command"
)

type Handler interface {
	Init(*Game) error
	UserInput(*Game) error
	Update(*Game, []*command.Command) error // Might need to add some arguements here so that networked data can be passed here.
	Draw() error

	IsInitialized() bool
}

type Game struct {
	Scenes      []Handler
	activeScene int

	ScreenSize rl.Vector2

	Opts *cache.Cache[types.Options]

	Paused bool

	networkWait sync.WaitGroup
	connected   bool
	ReadQueue   chan *command.Command
	WriteQueue  chan *command.Command
}

func New(screenx, screeny float32, title string, Scenes []Handler) *Game {
	rl.InitWindow(int32(screenx), int32(screeny), title)
	return &Game{
		Scenes:      Scenes,
		activeScene: 0,
		ScreenSize: rl.Vector2{
			X: screenx,
			Y: screeny,
		},
		Opts: cache.New(time.Duration(5)*time.Second, func(o *types.Options) error {
			if o == nil {
				return fmt.Errorf("Options pointer can not be nil")
			}
			exepath, err := os.Executable()
			if err != nil {
				return err
			}
			fullpath := filepath.Join(filepath.Dir(exepath), "config/General.toml")
			b, err := os.ReadFile(fullpath)
			if err != nil {
				return err
			}
			return toml.Unmarshal(b, o)
		}),
		ReadQueue:  make(chan *command.Command, 100),
		WriteQueue: make(chan *command.Command, 100),
	}
}

func (g *Game) IsConnected() bool {
	return g.connected
}

func (g *Game) SetScene(index int) error {
	if index < 0 || index >= len(g.Scenes) {
		return fmt.Errorf("specified index out of range")
	}
	g.activeScene = index
	if g.Scenes[g.activeScene].IsInitialized() {
		return nil
	}
	return g.Scenes[g.activeScene].Init(g)
}

func (g *Game) Run() {
	ctx, ctxclose := context.WithCancel(context.Background())

	defer rl.CloseWindow()
	defer g.networkWait.Wait()
	defer close(g.ReadQueue)
	defer close(g.WriteQueue)
	defer ctxclose()

	if !g.Scenes[g.activeScene].IsInitialized() {
		err := g.Scenes[g.activeScene].Init(g)
		if err != nil {
			fmt.Printf("Warning Error Initializing Scene (%d): %s\n", g.activeScene, err)
		}
	}
	go g.Network(ctx)
	for !rl.WindowShouldClose() {
		g.UserInput() // Pause State Agnostic
		g.Update()    // will not work while game is paused
		g.Draw()
	}
}
func (g *Game) UserInput() {

	if rl.IsKeyPressed(rl.KeyP) && !rl.IsKeyPressedRepeat(rl.KeyP) {
		g.Paused = !g.Paused
	}

	if !g.Paused {
		err := g.Scenes[g.activeScene].UserInput(g)
		if err != nil {
			fmt.Printf("Warning Game Scene (%d) User Input Error: %s\n", g.activeScene, err)
		}
	}
}

func (g *Game) Network(ctx context.Context) {
	g.networkWait.Add(1)
	defer g.networkWait.Done()

	var (
		conn net.Conn
		err  error
		wg   sync.WaitGroup
	)

	opts, err := g.Opts.Get()
	if err != nil {
		conn, err = net.Dial("tcp", "127.0.0.1:3000")
		if err != nil {
			panic(fmt.Errorf("failed to connect to server at %s:%d | %s", opts.Network.Addr, opts.Network.Port, err))
		}
	}
	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", opts.Network.Addr, opts.Network.Port))
	if err != nil {
		panic(fmt.Errorf("failed to connect to server at %s:%d | %s", opts.Network.Addr, opts.Network.Port, err))
	}

	g.connected = true

	// Reading to network connection
	go func(c *net.Conn, w *sync.WaitGroup, ctx context.Context) {
		timeoutCount := 0
		header := make([]byte, command.HeaderSize)
		buffer := make([]byte, opts.Network.BufferSize)
		msg := make([]byte, opts.Network.BufferSize*3)
		for !rl.WindowShouldClose() {
			select {
			case <-ctx.Done():
				return
			default:
				_, err := (*c).Read(header)
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					}
					var opErr net.Error
					if errors.As(err, &opErr) && opErr.Timeout() {
						timeoutCount++
						continue
					}
					if errors.Is(err, syscall.ECONNRESET) {
						return
					}
				}
				l, _, err := command.ValidateHeader(header)
				if err != nil {
					fmt.Printf("failed to validate header: %s\n", err)
					continue
				}

				buffer = make([]byte, l)
				n, err := (*c).Read(buffer)
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					}
					var opErr net.Error
					if errors.As(err, &opErr) && opErr.Timeout() {
						timeoutCount++
						continue
					}
					if errors.Is(err, syscall.ECONNRESET) {
						return
					}
				}
				if n != int(l) {
					fmt.Printf("Warning, did not get expected command length from server: %d, %d\n", n, l)
				}
				msg = make([]byte, int(l)+(int(command.HeaderSize)/8))
				for i, v := range header {
					msg[i] = v
				}
				for i, v := range buffer {
					msg[i+int(int(command.HeaderSize)/8)] = v
				}

				cmd, err := command.Unmarshal(msg)
				if err != nil {
					fmt.Printf("error Unmarshaling command: %s\n", err)
					continue
				}
				g.ReadQueue <- cmd

			}
		}
	}(&conn, &wg, ctx)

	// Writing to network connection
	go func(c *net.Conn, w *sync.WaitGroup, wc <-chan *command.Command, ctx context.Context) {
		for !rl.WindowShouldClose() {
			select {
			case <-ctx.Done():
				return
			case cmd := <-g.WriteQueue:
				_, err := (*c).Write(cmd.Marshal())
				if err != nil {
					fmt.Printf("%s\n", err)
				}
			}
		}
	}(&conn, &wg, g.WriteQueue, ctx)

	select {
	case <-ctx.Done():
		conn.Close()
		g.connected = false
		return
	}

}

func (g *Game) Update() {

	// allows the game engine itself to react to commands instead of strictly having the scnenes themselves handle network interactions. this also gerantees that we read incoming messages from the server regardless of what scene it is as the other way would have been accedent prone
	var inboundcommands []*command.Command
	if len(g.ReadQueue) > 0 {
		inboundcommands = make([]*command.Command, len(g.ReadQueue))
		count := 0
		for len(g.ReadQueue) > 0 {
			cmd := <-g.ReadQueue
			inboundcommands[count] = cmd
			count++
		}
	}

	for _, v := range inboundcommands {
		fmt.Printf("Incoming Command: %s\n", string(v.Marshal()))
	}

	if !g.Paused {
		err := g.Scenes[g.activeScene].Update(g, inboundcommands) // will need to repace [][]byte{} with network traffic
		if err != nil {
			fmt.Printf("Warning Game Scene (%d) Draw Error: %s\n", g.activeScene, err)
		}
	}
}

func (g *Game) Draw() {
	rl.ClearBackground(rl.RayWhite)
	rl.BeginDrawing()
	err := g.Scenes[g.activeScene].Draw()
	if err != nil {
		fmt.Printf("Warning Game Scene (%d) Draw Error: %s\n", g.activeScene, err)
	}

	if g.Paused {
		rl.DrawText("Paused", 0, 0, 12, rl.Black)
	} else {
		rl.DrawText("Unpaused", 0, 0, 12, rl.Black)
	}

	rl.EndDrawing()
}
