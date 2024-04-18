package ops

import (
	"fmt"
)

func (o *Op) NegateNumber(varName string) string {
	cmd := ""
	cmd += o.MoveScore(varName, varName)
	cmd += o.ExecCond(
		fmt.Sprintf("score %s %s matches 0", varName, o.Namespace),
		true,
		o.MoveConst("1", varName),
	)
	cmd += o.ExecCond(
		fmt.Sprintf("score %s %s matches 1..", varName, o.Namespace),
		true,
		o.MoveConst("0", varName),
	)
	return cmd
}

func (o *Op) CompNumbers(cond string, ifcond bool, ra string, rb string, rx string) string {
	cmd := ""
	// If the numbers are equal, the result is 1, otherwise 0
	cmd += o.MoveConst("0", rx)
	cmd += o.MoveScore(ra, ra)
	cmd += o.MoveScore(rb, rb)
	cmd += o.ExecCond(fmt.Sprintf("score %s %s %s %s %s", ra, o.Namespace, cond, rb, o.Namespace), ifcond, o.MoveConst("1", rx))
	return cmd
}

func (o *Op) EqNumbers(ra string, rb string, rx string) string {
	return o.CompNumbers("=", true, ra, rb, rx)
}

func (o *Op) NeqNumbers(ra string, rb string, rx string) string {
	return o.CompNumbers("=", false, ra, rb, rx)
}

func (o *Op) GtNumbers(ra string, rb string, rx string) string {
	cmd := ""
	cmd += o.CompNumbers(">", true, ra, rb, rx)
	return cmd
}
func (o *Op) GteNumbers(ra string, rb string, rx string) string {
	cmd := ""
	cmd += o.CompNumbers(">=", true, ra, rb, rx)
	return cmd
}
func (o *Op) LtNumbers(ra string, rb string, rx string) string {
	cmd := ""
	cmd += o.CompNumbers("<", true, ra, rb, rx)
	return cmd
}
func (o *Op) LteNumbers(ra string, rb string, rx string) string {
	cmd := ""
	cmd += o.CompNumbers("<=", true, ra, rb, rx)
	return cmd
}

func (o *Op) EqStrings(ra string, rb string, rx string) string {
	cmd := ""
	// Two strings are equal if they are the same length and each character is the same.
	// First we check if the strings are the same length
	ralen := ra + "len"
	rblen := rb + "len"
	cmd += o.MoveConst("1", rx)
	cmd += o.LenString(ra, ralen)
	cmd += o.LenString(rb, rblen)
	cmd += o.MoveScore(ralen, ralen)
	cmd += o.MoveScore(rblen, rblen)
	cmd += o.ExecCond(
		fmt.Sprintf("score %s %s = %s %s", ralen, o.Namespace, rblen, o.Namespace),
		false,
		o.MoveConst("0", rx),
	)
	// If the strings are the same length, we compare the string values
	// The easiest way to do this with MC commands is to check the success of a data command trying to overwrite the value of a string
	// If the data command fails, the strings are equal

	condReg := rx + "__cond__"
	cmd += o.MoveScore(rx, condReg)

	stringComparison := ""
	tmpReg := "__tmp__"
	stringComparison += o.Move(ra, tmpReg)
	stringComparison += fmt.Sprintf(
		"execute store success storage %s:%s %s int 1 run data modify storage %s:%s %s set from storage %s:%s %s\n",
		o.Namespace, VarPath, rx,
		o.Namespace, VarPath, tmpReg,
		o.Namespace, VarPath, rb)
	cmd += o.ExecCond(
		fmt.Sprintf("score %s %s matches 1", condReg, o.Namespace),
		true,
		stringComparison,
	)
	cmd += o.NegateNumber(rx)
	return cmd
}

func (o *Op) NeqStrings(ra string, rb string, rx string) string {
	cmd := ""
	cmd += o.EqStrings(ra, rb, rx)
	cmd += o.NegateNumber(rx)
	return cmd
}