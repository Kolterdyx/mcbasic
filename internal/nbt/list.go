package nbt

type List struct {
	Value

	values []Value
}

func NewList(values ...Value) *List {
	return &List{
		values: values,
	}
}

func (l *List) ToString() string {
	var str string
	for _, v := range l.values {
		str += v.ToString() + ", "
	}
	if len(str) > 0 {
		str = str[:len(str)-2]
	}
	return "[" + str + "]"
}
