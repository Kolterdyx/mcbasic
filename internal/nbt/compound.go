package nbt

type Compound struct {
	Values map[string]Value
}

func NewCompound() *Compound {
	return &Compound{
		Values: make(map[string]Value),
	}
}

func (c *Compound) ToString() string {
	var str string
	for k, v := range c.Values {
		str += k + ": " + v.ToString() + ", "
	}
	if len(str) > 0 {
		str = str[:len(str)-2]
	}
	return "{" + str + "}"
}

func (c *Compound) Get(key string) Value {
	return c.Values[key]
}

func (c *Compound) Set(key string, value Value) {
	if c.Values == nil {
		c.Values = make(map[string]Value)
	}
	c.Values[key] = value
}
