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

func (c *Code) addInstruction(inst interfaces.Instruction) {
	c.Instructions = append(c.Instructions, inst)
}

func (c *Code) addInst(instType interfaces.InstructionType, args ...string) interfaces.IRCode {
	c.addInstruction(c.makeInst(instType, args...))
	return c
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
		body += instr.ToString()
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
