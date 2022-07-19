package data

type Val[T any] interface {
	InDB() bool
	Val() T
}

// Value 用于替代nil来确认各种数据库中是否存在某个值
type Value[T any] struct {
	inDB  bool
	value T
}

func NewValue[T any](inDB bool, value T) Val[T] {
	return Value[T]{
		inDB:  inDB,
		value: value,
	}
}

// InDB 是否存在与数据库中
func (v Value[T]) InDB() bool {
	return v.inDB
}

// Val 获取数据的值
func (v Value[T]) Val() T {
	return v.value
}
