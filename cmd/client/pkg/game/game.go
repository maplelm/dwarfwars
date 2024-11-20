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

/*
Needed functions for a struct to be a scene
*/
type Scene interface {
	Init(*Game) error                             // Should be called lazy right before scene goes into use
	UserInput(*Game) error                        // Handles any scene specific input requirements
	Update(*Game, []*command.Command) error       // Might need to add some arguements here so that networked data can be passed here.
	PausedUpdate(*Game, []*command.Command) error // This function will be the alternate update function for what will happen to scene logic if the game is paused
	Draw() error                                  // Logic for how to draw everythin in the scene to the screen
	Deconstruct() error                           // Should be used to clean up the Scene before it gets deleted from the queue
	IsInitialized() bool                          // Check with a scene has been initialized before forcusing it.
	OnResize() error                              // Scene behavior when window is resized. do things need to be scaled? that sort of thig
}

type Game struct {
	Scenes      []Scene
	activeScene int

	ScreenSize rl.Vector2

	Opts *cache.Cache[types.Options]

	Paused bool

	networkWait sync.WaitGroup
	connected   bool
	connecting  bool
	ReadQueue   chan *command.Command
	WriteQueue  chan *command.Command

	Ctx      context.Context
	ctxClose func()

	NetworkCtx      context.Context
	NetworkCtxClose func()

	SWidth  int
	SHeight int

	ServerID uint32

	Scale float32
}

func New(screenx, screeny float32, title string, opts *cache.Cache[types.Options], scale float32, Scenes []Scene) *Game {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(int32(screenx), int32(screeny), title)
	ctx, cc := context.WithCancel(context.Background())
	nctx, ncc := context.WithCancel(ctx)
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

		ReadQueue:       make(chan *command.Command, 100),
		WriteQueue:      make(chan *command.Command, 100),
		Scale:           scale,
		Ctx:             ctx,
		ctxClose:        cc,
		NetworkCtx:      nctx,
		NetworkCtxClose: ncc,
	}
}

func (g *Game) IsConnected() bool {
	return g.connected
}

func (g *Game) IsConnecting() bool {
	return g.connecting
}

func (g *Game) SetScene(index int) error {
	if index < 0 || index >= len(g.Scenes) {
		return OOB{}
	}

	g.activeScene = index
	if g.Scenes[g.activeScene].IsInitialized() {
		return nil
	}
	return g.Scenes[g.activeScene].Init(g)
}

// remove current scene from ram and then move to the next one on the stack
func (g *Game) PopScene() error {
	if len(g.Scenes) <= 1 {
		return fmt.Errorf("can't pop last scene off stack")
	}
	if err := g.Scenes[g.activeScene].Deconstruct(); err != nil {
		fmt.Printf("Warning: Active Scene %d failed to Deconstruct!, %s\n", g.activeScene, err)
	}
	g.Scenes = append(g.Scenes[:g.activeScene], g.Scenes[g.activeScene+1:]...)
	g.SetScene(g.activeScene - 1)
	return nil
}

// Goes to a new scene while keeping the current in ram
func (g *Game) PushScene(s Scene) error {
	g.Scenes = append(g.Scenes[:g.activeScene+1], append([]Scene{s}, g.Scenes[g.activeScene+1:]...)...)
	err := g.SetScene(g.activeScene + 1)
	if err != nil {
		fmt.Printf("failing to set scene, %s\n", err)
	}
	return err
}

func (g *Game) NextScene(shift int) error {
	if shift+g.activeScene < 0 || shift+g.activeScene >= len(g.Scenes) {
		return fmt.Errorf("the spcified shift would escape the bounds of the slice")
	}
	return g.SetScene(g.activeScene + shift)
}

func (g *Game) ShiftRight(shift int) error {
	if shift+g.activeScene < 0 || shift+g.activeScene >= len(g.Scenes) {
		return OOB{Details: "shift width invalid"}
	}
	as := g.Scenes[g.activeScene]
	if err := g.PopScene(); err != nil {
		return err
	}
	if err := g.NextScene(shift); err != nil {
		return err
	}
	return g.PushScene(as)
}

func (g *Game) ShiftLeft(shift int) error {
	return g.ShiftRight(shift * -1)
}

// replaces the current seen with a new one rather then saving the current one in ram while going to a new one
func (g *Game) ReplaceScene(s Scene) error {
	g.Scenes[g.activeScene] = s
	return g.SetScene(g.activeScene)
}

func (g *Game) Run() {
	// Make sure any running network code has properly closed down before returning
	defer g.networkWait.Wait()
	// Make sure that the Raylib Window is closed
	defer rl.CloseWindow()
	// Cancel the Game Context so any go routines will start to shutdown
	defer g.ctxClose()
	// Cancel the Networking context so networking code will start to shutdown
	defer g.NetworkCtxClose()

	// Making sure the current scene has been initialized before running.
	if !g.Scenes[g.activeScene].IsInitialized() {
		if err := g.Scenes[g.activeScene].Init(g); err != nil {
			fmt.Printf("Warning Error Initializing Scene (%d): %s\n", g.activeScene, err)
		}
	}

	// Main Game Loop
	for !rl.WindowShouldClose() || (rl.WindowShouldClose() && rl.IsKeyPressed(rl.KeyEscape)) {
		g.UserInput() // will not work while game is paused
		g.Update()    // will not work while game is paused
		g.Draw()      // Draw the overlay after calling the scene's draw function
	}

}
func (g *Game) UserInput() {

	// Pause the game if the p key is pressed
	if rl.IsKeyPressed(rl.KeyP) && !rl.IsKeyPressedRepeat(rl.KeyP) {
		g.Paused = !g.Paused
	}

	// Run the Sence Update function if game is not paused
	if !g.Paused {
		err := g.Scenes[g.activeScene].UserInput(g)
		if err != nil {
			fmt.Printf("Warning Game Scene (%d) User Input Error: %s\n", g.activeScene, err)
		}
	}
}

func (g *Game) Network(ctx context.Context, addr string, port int) error {
	g.connecting = true
	defer func() { g.connecting = false }()
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
	if cmd.Type == command.TypeWelcome && cmd.Format == command.FormatText {
		g.ServerID = cmd.ClientID
		g.connected = true
		g.connecting = false
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
				cmd.ClientID = g.ServerID
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

	// Printing commands recieved from server before sending to scene
	for _, v := range inboundcommands {
		fmt.Printf("Command From Client: %d\n\tFormat: %d\n\type:%d\n\tData: %s\n\n", v.ClientID, v.Format, v.Type, string(v.Data))
	}

	if rl.IsWindowResized() {
		g.SWidth = rl.GetScreenWidth()
		g.SHeight = rl.GetScreenHeight()
		for i := range g.Scenes {
			if err := g.Scenes[i].OnResize(); err != nil {
				fmt.Printf("Warning: Error on OnResize %d, %s\n", g.activeScene, err)
			}
		}
	}

	// Do not update scene if paused
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

	// Drawing Scene
	err := g.Scenes[g.activeScene].Draw()
	if err != nil {
		fmt.Printf("Warning Game Scene (%d) Draw Error: %s\n", g.activeScene, err)
	}

	/////////////////////
	// Drawing Overlay //
	/////////////////////

	var str string
	var col rl.Color
	if g.Paused {
		str = "Paused"
		col = rl.Red
		rl.DrawText(str, 100, 100, 12, col)
	}

	if !g.connected {
		str = "Disconnected"
		col = rl.Red
		rl.DrawText(str, 100, 100, 12, col)
	}

	rl.EndDrawing()
}
