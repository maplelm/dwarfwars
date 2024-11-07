package gui

import (
	"math"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type button struct {
	Clicked bool
	Label   string
	Action  func()
}

type ButtonList struct {
	Position   rl.Vector2
	Buttonsize rl.Vector2
	Width      int
	List       []button
	Scale      *float32
}

func NewButtonList(s rl.Vector2, p rl.Vector2, w int, sc *float32) *ButtonList {
	return &ButtonList{
		Buttonsize: s,
		Position:   p,
		Width:      w,
		Scale:      sc,
		List:       make([]button, 0),
	}
}

func (bl *ButtonList) Add(l string, a func()) int {
	bl.List = append(bl.List, button{
		Clicked: false,
		Label:   l,
		Action:  a,
	})
	return len(bl.List) - 1
}

func (bl *ButtonList) Draw() {
	for i, v := range bl.List {
		if bl.Scale != nil {
			v.Clicked = raygui.Button(rl.NewRectangle(bl.Position.X+bl.Buttonsize.X*float32(i%bl.Width), bl.Position.Y+((bl.Buttonsize.Y**bl.Scale)*float32(math.Floor(float64(i)/float64(bl.Width)))), bl.Buttonsize.X, bl.Buttonsize.Y), v.Label)
		} else {
			v.Clicked = raygui.Button(rl.NewRectangle(bl.Position.X+bl.Buttonsize.X*float32(i%bl.Width), bl.Position.Y+bl.Buttonsize.Y*float32(math.Floor(float64(i)/float64(bl.Width))), bl.Buttonsize.X, bl.Buttonsize.Y), v.Label)
		}

		/*
			if bl.Scale != nil {
				v.Clicked = raygui.Button(rl.NewRectangle(bl.Position.X, bl.Position.Y+((bl.Buttonsize.Y**bl.Scale)*float32(i)), bl.Buttonsize.X, bl.Buttonsize.Y), v.Label)
			} else {
				v.Clicked = raygui.Button(rl.NewRectangle(bl.Position.X, bl.Position.Y+bl.Buttonsize.Y*float32(i), bl.Buttonsize.X, bl.Buttonsize.Y), v.Label)
			}
		*/
		bl.List[i] = v
	}
}

func (bl *ButtonList) Execute() {
	for _, v := range bl.List {
		if v.Clicked {
			v.Action()
		}
	}
}

func (bl *ButtonList) Centered() rl.Vector2 {
	var (
		x float32
		y float32
	)
	if len(bl.List) < bl.Width {
		x = bl.Position.X - (((bl.Buttonsize.X * (*bl.Scale)) * float32(len(bl.List))) / 2)
	} else {
		x = bl.Position.X - (((bl.Buttonsize.X * (*bl.Scale)) * float32(bl.Width)) / 2)
	}
	y = bl.Position.Y - ((bl.Buttonsize.Y * (*bl.Scale) * float32(len(bl.List))) / 2)
	return rl.Vector2{
		X: x,
		Y: y,
	}
}
