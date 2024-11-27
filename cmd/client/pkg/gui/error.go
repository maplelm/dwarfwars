package gui

type ErrNoGraphics string

func (eng ErrNoGraphics) Error() string {
	return string(eng)
}
