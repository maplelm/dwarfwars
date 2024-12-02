package command

type DataWrapper[T any] struct {
	Secret []byte
	Data   T
}

func NewDataWrapper[T any](s []byte, data T) *DataWrapper[T] {
	return &DataWrapper[T]{
		Secret: s,
		Data:   data,
	}
}
