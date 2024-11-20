package button

import rl "github.com/gen2brain/raylib-go/raylib"

type List struct {
	Position rl.Vector2
	ButtonSize rl.Vector2
	BorderWidth int32
	BorderColor rl.Color
	FontSize int32
	FontColor rl.Color
	ForegroundColor rl.Color
	Buttons  []Button
	Columns  int32
	Scale    float32
	Font     rl.Font
}

func NewList(p rl.Vector2, f rl.Font, cols int32, s float32, borderwidth int32, bordercolor rl.Color, foreground rl.Color,  fontsize int32, fontcolor rl.Color buttonsize rl.Vector2) *List {
	return &List{
		Position: p,
		Font:     f,
		Columns:  cols,
		Scale:    s,
		Buttons:  make([]Button, 0),
		BorderWidth: borderwidth,
		BorderColor: bordercolor,
		ForegroundColor: foreground,
		FontSize: fontsize,
		FontColor: fontcolor,
		ButtonSize: buttonsize,
	}
}

func (l *List) Add(lab string, a func(), s rl.Vector2) int {
	
	// Calculate Posistion of button
	buttonCount := len(l.Buttons)
	x := buttonCount % l.Columns
	y := math.Floor(buttonCount/l.Columns)

	l.Buttons = append(l.Buttons, *New(
			lab,
			a,
			rl.NewRectangle(
				l.Position.X + (x * l.ButtonSize.X) + ( x * (l.BorderWidth*2)),
				l.Position.Y + (y * l.ButtonSize.Y) + (y * l.BorderWidth*2), 
				l.ButtonSize.X,
				l.ButtonSize.Y
			),
			l.BorderWidth,
			l.BorderColor,
			l.ForegroundColor,
			l.Font,
			l.FontSize,
			l.FontSize,
			l.FontColor,
			
		))
}

func (l *List) Move(v rl.Vector2) {
	l.Position = rl.Vector2Add(l.Position, v)
	for i, b := range l.Buttons {
		x := i % l.Columns
		y := math.Floor(i/l.Columns)
		l.Buttons[i].Bounds.X = l.Position.X + ( x * l.ButtonSize.X) + (x * (l.BorderWidth*2))
		l.Buttons[i].Bounds.Y = l.Position.Y + ( Y * l.ButtonSize.Y) + (y * (l.ButtonSize*2))
	}
}

func (l *List) Draw() {
}
