package gui

type Button struct {
	Label   string
	Clicked bool
	Action  func()
}

// func InitButton(l string, a func()) Button {
// 	return Button {
// 		Label: l,
// 		Clicked: false,
// 		Action: a,
// 	}
// }
