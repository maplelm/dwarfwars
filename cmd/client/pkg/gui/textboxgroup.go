package gui

import (
	"fmt"

	rlgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Textbox struct {
	Label        string
	Password     bool
	Email        bool
	Data         string
	Active       bool
	LineSize     rl.Vector2
	LabelScaling float32
	Gap          float32
}

func InitTextbox(linesize rl.Vector2, label string, ispassword, isemail, startactive bool, labelscale, margin float32) Textbox {
	return Textbox{
		Label:        label,
		Password:     ispassword,
		Email:        isemail,
		Data:         "",
		Active:       startactive,
		LineSize:     linesize,
		LabelScaling: labelscale,
		Gap:          margin,
	}
}

func (tb *Textbox) Draw(Position rl.Vector2, charlim int, fc rl.Color) {
	rl.DrawText(tb.Label, int32(Position.X), int32(Position.Y), int32((tb.LineSize.Y * tb.LabelScaling)), fc)
	if rlgui.TextBox(rl.NewRectangle(Position.X, Position.Y+tb.LineSize.Y+tb.Gap, tb.LineSize.X, tb.LineSize.Y), &tb.Data, charlim, tb.Active) {
		tb.Active = !tb.Active
	}
}

func (tb *Textbox) Bounds(Position rl.Vector2) rl.Rectangle {
	return rl.NewRectangle(Position.X, Position.Y, tb.LineSize.X, tb.LineSize.Y+(tb.LineSize.Y*tb.LabelScaling)+tb.Gap)
}

type TextBoxGroup struct {
	Position       rl.Vector2
	Size           rl.Vector2
	CharacterLimit int
	LabelColor     rl.Color
	List           []Textbox
	Gap            float32
}

func NewTextBoxGroup(rect rl.Rectangle, charlimit int, g float32, lc rl.Color) *TextBoxGroup {
	return &TextBoxGroup{
		Position:       rl.Vector2{X: rect.X, Y: rect.Y},
		Size:           rl.Vector2{X: rect.Width, Y: rect.Height},
		CharacterLimit: charlimit,
		LabelColor:     lc,
		Gap:            g,
		List:           make([]Textbox, 0),
	}
}

func (tbg *TextBoxGroup) Add(t Textbox) {
	tbg.List = append(tbg.List, t)
}

func (tbg *TextBoxGroup) AddMulti(t []Textbox) {
	tbg.List = append(tbg.List, t...)
}

func (tbg *TextBoxGroup) ValueByLabel(label string) (string, error) {
	for _, v := range tbg.List {
		if label == v.Label {
			return v.Data, nil
		}
	}
	return "", fmt.Errorf("label not found: %s", label)
}

func (tbg *TextBoxGroup) Draw() {
	cursor := tbg.Position

	for i := range tbg.List {
		tbg.List[i].Draw(cursor, tbg.CharacterLimit, tbg.LabelColor)
		b := tbg.List[i].Bounds(cursor)
		cursor.Y += b.Height + tbg.Gap
	}
}

func (tbg *TextBoxGroup) Center() {
}
