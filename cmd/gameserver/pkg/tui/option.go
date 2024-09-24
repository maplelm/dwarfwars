package tui

type Option struct {
	Label, Desc string
	Action      func() error
}

func (o Option) Title() string {
	return o.Label
}

func (o Option) Description() string {
	return o.Desc
}

func (o Option) FilterValue() string {
	return o.Label
}
