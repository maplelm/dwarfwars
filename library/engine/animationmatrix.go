package engine

import (
	"fmt"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type AnimationMatrix struct {
	// Grid Setup
	width   int        // Grid Width
	height  int        // Grid Height
	frames  int        // How many idexes are used on the grid
	boxsize rl.Vector2 // size of each grid box

	// Animation State
	fps               int          // how often frames change
	current           int          // current frame index
	lastframetick     time.Time    // last time the frame was changed
	spritesheet       rl.Texture2D // spritesheet to use grid and index with
	rotationanimation func(int, float32) float32
	rotation          float32

	// Positions
	Bounds rl.Rectangle

	// Texture Moifications
	tint rl.Color
}

func NewAnimationMatrix(w, h, frames, fps int, ss rl.Texture2D, boxsize rl.Vector2, tint rl.Color, rf func(int, float32) float32) (*AnimationMatrix, error) {
	if frames > w*h {
		return nil, fmt.Errorf("matrix size does not accomidate frame count")
	}
	if w <= 0 || h <= 0 || frames <= 0 || fps < 0 {
		return nil, fmt.Errorf("negative values are not allowed in an animation matrix")
	}

	if rf == nil {
		rf = func(c int, r float32) float32 { return 0 }
	}

	return &AnimationMatrix{
		width:             w,
		height:            h,
		frames:            frames,
		fps:               fps,
		current:           0,
		rotationanimation: rf,
		spritesheet:       ss,
	}, nil
}

func (am *AnimationMatrix) Width() int {
	return am.width
}

func (am *AnimationMatrix) Height() int {
	return am.height
}

func (am *AnimationMatrix) Frames() int {
	return am.frames
}

func (am *AnimationMatrix) Fps() int {
	return am.fps
}

func (am *AnimationMatrix) CurrentFrame() int {
	return am.current
}

func (am *AnimationMatrix) SetWidth(v int) error {
	if v > 0 {
		am.width = v
		return nil
	}
	return fmt.Errorf("width can't be set to a negative number")
}

func (am *AnimationMatrix) SetHeight(v int) error {
	if v > 0 {
		am.height = v
		return nil
	}
	return fmt.Errorf("height can't be set to a negative number")
}

func (am *AnimationMatrix) SetFrames(v int) error {
	if v >= 0 && v > am.width*am.height {
		am.frames = v
		return nil
	}
	if v <= 0 {
		return fmt.Errorf("frames can't be negative or 0")
	}
	return fmt.Errorf("frames has to be less then width * height (%d), frames value given = %d", am.width*am.height, v)
}

func (am *AnimationMatrix) SetFps(v int) error {
	if v >= 0 {
		am.fps = v
		return nil
	}
	return fmt.Errorf("fps can't be negative")
}

func (am *AnimationMatrix) SetCurrentFrame(v int) error {
	if v >= am.width*am.height || v < 0 {
		return fmt.Errorf("frame %d does not exist", v)
	}
	am.current = v
	return nil
}

func (am *AnimationMatrix) NextFrame() {
	if am.current == am.frames-1 {
		am.current = 0
	} else {
		am.current++
	}
}

func (am *AnimationMatrix) DrawAnimationFrame() {
	x := am.current % am.width
	y := int(math.Floor(float64(am.current) / float64(am.width)))
	if am.fps > 0 && time.Since(am.lastframetick) >= time.Second/time.Duration(am.fps) {
		am.NextFrame()
		am.lastframetick = time.Now()
	}

	rl.DrawTexturePro(am.spritesheet, rl.NewRectangle(float32(x)*am.boxsize.X, float32(y)*am.boxsize.Y, am.boxsize.X, am.boxsize.Y), am.Bounds, rl.Vector2{X: am.Bounds.X + (am.Bounds.Width / 2), Y: am.Bounds.Y + (am.Bounds.Height / 2)}, am.rotationanimation(am.current, am.rotation), am.tint)
}
