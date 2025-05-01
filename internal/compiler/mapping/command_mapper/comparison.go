package command_mapper

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

func (c *CommandMapper) NegateNumber(varName string) string {
	cmd := ""
	cmd += c.MoveScore(varName, varName)
	cmd += c.ExecCond(
		fmt.Sprintf("score %s %s matches 0", varName, c.Namespace),
		true,
		c.MakeConst(nbt.NewInt(1), varName),
	)
	cmd += c.ExecCond(
		fmt.Sprintf("score %s %s matches 1..", varName, c.Namespace),
		true,
		c.MakeConst(nbt.NewInt(0), varName),
	)
	return cmd
}

func (c *CommandMapper) CompNumbers(cond string, ifcond bool, ra string, rb string, rx string) string {
	cmd := ""
	// If the numbers match the condition, the result is 1, otherwise 0
	cmd += c.MakeConst(nbt.NewInt(0), rx)
	cmd += c.MoveScore(ra, ra)
	cmd += c.MoveScore(rb, rb)
	cmd += c.ExecCond(fmt.Sprintf("score %s %s %s %s %s", ra, c.Namespace, cond, rb, c.Namespace), ifcond, c.MakeConst(nbt.NewInt(1), rx))
	return cmd
}

func (c *CommandMapper) EqNumbers(ra string, rb string, rx string) string {
	return c.CompNumbers("=", true, ra, rb, rx)
}

func (c *CommandMapper) NeqNumbers(ra string, rb string, rx string) string {
	return c.CompNumbers("=", false, ra, rb, rx)
}

func (c *CommandMapper) GtNumbers(ra string, rb string, rx string) string {
	cmd := ""
	cmd += c.CompNumbers(">", true, ra, rb, rx)
	return cmd
}
func (c *CommandMapper) GteNumbers(ra string, rb string, rx string) string {
	cmd := ""
	cmd += c.CompNumbers(">=", true, ra, rb, rx)
	return cmd
}
func (c *CommandMapper) LtNumbers(ra string, rb string, rx string) string {
	cmd := ""
	cmd += c.CompNumbers("<", true, ra, rb, rx)
	return cmd
}
func (c *CommandMapper) LteNumbers(ra string, rb string, rx string) string {
	cmd := ""
	cmd += c.CompNumbers("<=", true, ra, rb, rx)
	return cmd
}

func (c *CommandMapper) EqStrings(ra string, rb string, rx string) string {
	cmd := ""
	// Two strings are equal if they are the same length and each character is the same.
	// First we check if the strings are the same length
	ralen := ra + "len"
	rblen := rb + "len"
	cmd += c.MakeConst(nbt.NewInt(1), rx)
	cmd += c.Size(ra, ralen)
	cmd += c.Size(rb, rblen)
	cmd += c.MoveScore(ralen, ralen)
	cmd += c.MoveScore(rblen, rblen)
	cmd += c.ExecCond(
		fmt.Sprintf("score %s %s = %s %s", ralen, c.Namespace, rblen, c.Namespace),
		false,
		c.MakeConst(nbt.NewInt(0), rx),
	)
	// If the strings are the same length, we compare the string values
	// The easiest way to do this with MC commands is to check the success of a data command trying to overwrite the value of a string
	// If the data command fails, the strings are equal

	condReg := rx + "__cond__"
	cmd += c.MoveScore(rx, condReg)

	stringComparison := ""
	tmpReg := "__tmp__"
	stringComparison += c.Move(ra, tmpReg)
	stringComparison += fmt.Sprintf(
		"execute store success storage %s:%s %s int 1 run data modify storage %s:%s %s set from storage %s:%s %s\n",
		c.Namespace, c.VarPath, rx,
		c.Namespace, c.VarPath, tmpReg,
		c.Namespace, c.VarPath, rb)
	cmd += c.ExecCond(
		fmt.Sprintf("score %s %s matches 1", condReg, c.Namespace),
		true,
		stringComparison,
	)
	cmd += c.NegateNumber(rx)
	return cmd
}

func (c *CommandMapper) NeqStrings(ra string, rb string, rx string) string {
	cmd := ""
	cmd += c.EqStrings(ra, rb, rx)
	cmd += c.NegateNumber(rx)
	return cmd
}
