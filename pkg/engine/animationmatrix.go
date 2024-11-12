package engine

import (
	"fmt"
	"math"
)

type AnimationMatrix struct {
	width   int
	height  int
	frames  int
	fps     int
	current int
}

func NewAnimationMatrix(w, h, frames, fps int) (*AnimationMatrix, error) {
	if frames > w*h {
		return nil, fmt.Errorf("matrix size does not accomidate frame count")
	}
	if w <= 0 || h <= 0 || frames <= 0 || fps < 0 {
		return nil, fmt.Errorf("negative values are not allowed in an animation matrix")
	}
	return &AnimationMatrix{
		width:   w,
		height:  h,
		frames:  frames,
		fps:     fps,
		current: 0,
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

func (am *AnimationMatrix) AnimationFrame() (x int, y int) {
	x = am.current % am.width
	y = int(math.Floor(float64(am.current) / float64(am.width)))

	if am.current >= am.frames-1 {
		am.current = 0
	} else {
		am.current++
	}

	return
}
