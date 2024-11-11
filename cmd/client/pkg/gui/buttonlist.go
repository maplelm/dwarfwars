package gui

import (
	"math"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Button struct {
	Clicked bool
	Label   string
	Action  func()
}

func InitButton(l string, a func()) Button {
	return Button{
		Clicked: false,
		Label:   l,
		Action:  a,
	}
}

type ButtonList struct {
	Position   rl.Vector2
	Buttonsize rl.Vector2
	Width      int
	List       []Button
	Scale      *float32
}

func NewButtonList(dims rl.Rectangle, w int, scale *float32) *ButtonList {
	return &ButtonList{
		Buttonsize: rl.Vector2{X: dims.Width, Y: dims.Height},
		Position:   rl.Vector2{X: dims.X, Y: dims.Y},
		Width:      w,
		Scale:      scale,
		List:       make([]Button, 0),
	}
}

func (bl *ButtonList) Add(l string, a func()) int {
	bl.List = append(bl.List, Button{
		Clicked: false,
		Label:   l,
		Action:  a,
	})
	return len(bl.List) - 1
}

func (bl *ButtonList) AddMulti(bs []Button) int {
	bl.List = append(bl.List, bs...)
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

func (bl *ButtonList) Center() {
	if len(bl.List) < bl.Width {
		bl.Position.X = bl.Position.X - (((bl.Buttonsize.X * (*bl.Scale)) * float32(len(bl.List))) / 2)
	} else {
		bl.Position.X = bl.Position.X - (((bl.Buttonsize.X * (*bl.Scale)) * float32(bl.Width)) / 2)
	}
	bl.Position.Y = bl.Position.Y - ((bl.Buttonsize.Y*(*bl.Scale))*(float32(math.Floor(float64(len(bl.List))/float64(bl.Width)))))/2
}
