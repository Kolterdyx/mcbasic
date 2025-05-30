package ir

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type Code struct {
	interfaces.IRCode

	Instructions []interfaces.Instruction
	Namespace    string
	storage      string
}

func NewCode(namespace, storage string) interfaces.IRCode {
	return &Code{
		Instructions: make([]interfaces.Instruction, 0),
		Namespace:    namespace,
		storage:      storage,
	}
}

func (c *Code) AddInstruction(inst interfaces.Instruction) interfaces.IRCode {
	c.Instructions = append(c.Instructions, inst)
	return c
}

func (c *Code) addInst(instType interfaces.InstructionType, args ...string) interfaces.IRCode {
	return c.AddInstruction(c.makeInst(instType, args...))
}

func (c *Code) makeInst(instType interfaces.InstructionType, args ...string) interfaces.Instruction {
	return Instruction{
		Type:        instType,
		DPNamespace: c.Namespace,
		Storage:     c.storage,
		Args:        args,
	}
}

func (c *Code) ToString() string {
	body := ""
	for _, instr := range c.Instructions {
		body += instr.ToString() + "\n"
	}
	return body
}

func (c *Code) ToMCCode() string {
	body := ""
	for _, instr := range c.Instructions {
		body += instr.ToMCCommand()
	}
	return body
}

func (c *Code) Extend(code interfaces.IRCode) interfaces.IRCode {
	c.Instructions = append(c.Instructions, code.GetInstructions()...)
	return c
}

func (c *Code) GetNamespace() string {
	return c.Namespace
}

func (c *Code) GetStorage() string {
	return c.storage
}

func (c *Code) GetInstructions() []interfaces.Instruction {
	return c.Instructions
}

func (c *Code) SetInstructions(instructions []interfaces.Instruction) {
	c.Instructions = instructions
}

func (c *Code) Len() int {
	return len(c.Instructions)
}
