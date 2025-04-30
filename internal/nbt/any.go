package nbt

type Any struct {
	Value

	value string
}

func NewAny(value string) *Any {
	return &Any{
		value: value,
	}
}

func (a *Any) ToString() string {
	return a.value
}
