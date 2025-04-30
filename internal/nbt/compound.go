package nbt

import (
	"github.com/elliotchance/orderedmap/v3"
)

type Compound struct {
	Value
	values *orderedmap.OrderedMap[string, Value]
}

func NewCompound() *Compound {
	return &Compound{
		values: orderedmap.NewOrderedMap[string, Value](),
	}
}

func (c *Compound) ToString() string {
	var str string
	for k, v := range c.values.AllFromFront() {
		str += k + ": " + v.ToString() + ", "
	}
	if len(str) > 0 {
		str = str[:len(str)-2]
	}
	return "{" + str + "}"
}

func (c *Compound) Get(key string) (Value, bool) {
	return c.values.Get(key)
}

func (c *Compound) Set(key string, value Value) {
	if c.values == nil {
		c.values = orderedmap.NewOrderedMap[string, Value]()
	}
	c.values.Set(key, value)
}

func (c *Compound) Size() int {
	return c.values.Len()
}
