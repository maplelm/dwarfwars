package textbox

import rl "github.com/gen2brain/raylib-go/raylib"

type Textbox struct {
	Position  rl.Vector2
	MaxLength int

	value    string
	selected bool
}

func New() *Textbox {
	return &Textbox{}
}

func (t *Textbox) Update() {
}

func (t *Textbox) Draw() {
}
