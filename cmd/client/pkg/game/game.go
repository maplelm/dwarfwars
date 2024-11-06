package game

import (
	"context"
	"encoding/json"
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

	Ctx      context.Context
	ctxClose func()
}

func New(screenx, screeny float32, title string, opts *cache.Cache[types.Options], Scenes []Handler) *Game {
	rl.InitWindow(int32(screenx), int32(screeny), title)
	return &Game{
		Scenes:      Scenes,
		activeScene: 0,
		ScreenSize: rl.Vector2{
			X: screenx,
			Y: screeny,
		},
		Opts: func() *cache.Cache[types.Options] {

			if opts != nil {
				return opts
			}

			return cache.New(time.Duration(5)*time.Second, func(o *types.Options) error {
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
			})
		}(),
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
	g.Ctx, g.ctxClose = context.WithCancel(context.Background())

	defer g.networkWait.Wait()
	defer rl.CloseWindow()
	defer g.ctxClose()

	if !g.Scenes[g.activeScene].IsInitialized() {
		err := g.Scenes[g.activeScene].Init(g)
		if err != nil {
			fmt.Printf("Warning Error Initializing Scene (%d): %s\n", g.activeScene, err)
		}
	}
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

func (g *Game) Network(ctx context.Context, addr string, port int) error {
	g.networkWait.Add(1)
	defer g.networkWait.Done()

	rerrchan := make(chan error)
	werrchan := make(chan error)

	var (
		conn net.Conn
		err  error
		wg   sync.WaitGroup
	)

	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return err
	}
	defer conn.Close()

	cmd, _ := command.New(0, command.FormatText, command.TypeWelcome, []byte{})
	_, err = cmd.Send(conn)
	if err != nil {
		fmt.Printf("Failed to connect to Server %s:%d, %s\n", addr, port, err)
		return err
	}

	// Getting ID from server
	cmd, err = command.Recieve(conn)
	if err != nil {
		fmt.Printf("Failed to recieve welcome from Server %s:%d, %s\n", addr, port, err)
		return err
	}
	var id uint32
	if cmd.Type == command.TypeWelcome {
		id = cmd.ClientID
		g.connected = true
		defer func() { g.connected = false }()
	} else {
		fmt.Printf("Invalid Welcome Response from Server %s:%d!", addr, port)
		return fmt.Errorf("invalid welcome response")
	}

	// Reading to network connection
	go func(c *net.Conn, w *sync.WaitGroup, ctx context.Context, errchan chan<- error) {
		defer close(errchan)
		var timeoutCount int = 0
		for !rl.WindowShouldClose() {
			select {
			case <-ctx.Done():
				return
			default:
				cmd, err := command.Recieve(*c)
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					} else if opErr, ok := err.(net.Error); ok && opErr.Timeout() {
						timeoutCount++
						continue
					} else if errors.Is(err, syscall.ECONNRESET) {
						return
					} else {
						fmt.Printf("Network Read Error: %s", err)
					}
				}
				g.ReadQueue <- cmd
			}
		}
	}(&conn, &wg, ctx, rerrchan)

	// Writing to network connection
	go func(c *net.Conn, w *sync.WaitGroup, wc <-chan *command.Command, ctx context.Context, errchan chan<- error) {
		defer close(errchan)
		for !rl.WindowShouldClose() {
			select {
			case <-ctx.Done():
				return
			case cmd := <-g.WriteQueue:
				cmd.ClientID = id
				fmt.Printf("Sending data to server (%d), %s", cmd.ClientID, string(cmd.Data))
				n, err := cmd.Send(*c)
				if err != nil {
					fmt.Printf("Network Write Error: %s\n", err)
				} else {
					fmt.Printf("Network Write: %d bytes\n", n)
				}
			}
		}
	}(&conn, &wg, g.WriteQueue, ctx, werrchan)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-rerrchan:
		return err
	case err := <-werrchan:
		return err
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
		b, _ := json.Marshal(v)
		fmt.Printf("Incoming Command: %s\n", b)
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
