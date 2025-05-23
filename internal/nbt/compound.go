package nbt

type Compound struct {
	values map[string]Value
}

func NewCompound() Compound {
	return Compound{
		values: make(map[string]Value),
	}
}

func (c Compound) ToString() string {
	var str string
	for k, v := range c.values {
		str += k + ": " + v.ToString() + ", "
	}
	if len(str) > 0 {
		str = str[:len(str)-2]
	}
	return "{" + str + "}"
}

func (c Compound) Get(key string) (Value, bool) {
	val, ok := c.values[key]
	return val, ok
}

func (c Compound) Set(key string, value Value) Compound {
	if c.values == nil {
		c.values = make(map[string]Value)
	}
	c.values[key] = value
	return c
}

func (c Compound) Size() int {
	return len(c.values)
}
