package ir

import (
	"strings"
)

type OptFunc func(instrs []Instruction, i int) (matched bool, consumed int, replacement []Instruction)

func OptimizeFunctionBody(f Function) Function {
	optimized := make([]Instruction, 0)

	opts := []OptFunc{
		optSkipMacroSet,
		optCollapseCopyChain,
		optSetCopyToSet,
		optDedupRet,
	}

	i := 0
	for i < len(f.Instructions) {
		applied := false

		for _, opt := range opts {
			if matched, consumed, repl := opt(f.Instructions, i); matched {
				optimized = append(optimized, repl...)
				i += consumed
				applied = true
				break
			}
		}

		if !applied {
			optimized = append(optimized, f.Instructions[i])
			i++
		}
	}

	f.Instructions = optimized
	return f
}

func isMacroPath(path string) bool {
	return strings.Contains(path, "$(")
}

func isMacroSetPattern(setInst Instruction) bool {
	if setInst.Type != Set || len(setInst.Args) != 3 {
		return false
	}
	path := setInst.Args[1]
	val := setInst.Args[2]

	if strings.HasPrefix(path, "vars.") && strings.HasPrefix(val, "$(") && strings.HasSuffix(val, ")") {
		varName := strings.TrimPrefix(path, "vars.")
		macroName := strings.TrimSuffix(strings.TrimPrefix(val, "$("), ")")
		return varName == macroName
	}
	return false
}

func sameLocation(storageA, pathA, storageB, pathB string) bool {
	return storageA == storageB && pathA == pathB
}

func optSkipMacroSet(instrs []Instruction, i int) (bool, int, []Instruction) {
	curr := instrs[i]
	if curr.Type == Set && isMacroSetPattern(curr) {
		return true, 1, []Instruction{curr}
	}
	return false, 0, nil
}

func optCollapseCopyChain(instrs []Instruction, i int) (bool, int, []Instruction) {
	if instrs[i].Type != Copy || len(instrs[i].Args) != 4 {
		return false, 0, nil
	}

	srcStorage, srcPath := instrs[i].Args[0], instrs[i].Args[1]
	curDstStorage, curDstPath := instrs[i].Args[2], instrs[i].Args[3]

	chainEnd := i + 1
	for chainEnd < len(instrs) {
		next := instrs[chainEnd]
		if next.Type != Copy || len(next.Args) != 4 {
			break
		}
		if !sameLocation(curDstStorage, curDstPath, next.Args[0], next.Args[1]) {
			break
		}
		curDstStorage, curDstPath = next.Args[2], next.Args[3]
		chainEnd++
	}

	if chainEnd > i+1 {
		return true, chainEnd - i, []Instruction{{
			Type: Copy,
			Args: []string{srcStorage, srcPath, curDstStorage, curDstPath},
		}}
	}
	return false, 0, nil
}

func optSetCopyToSet(instrs []Instruction, i int) (bool, int, []Instruction) {
	if i+1 >= len(instrs) {
		return false, 0, nil
	}

	curr, next := instrs[i], instrs[i+1]
	if curr.Type == Set && len(curr.Args) == 3 &&
		next.Type == Copy && len(next.Args) == 4 &&
		sameLocation(curr.Args[0], curr.Args[1], next.Args[0], next.Args[1]) {

		return true, 2, []Instruction{{
			Type: Set,
			Args: []string{next.Args[2], next.Args[3], curr.Args[2]},
		}}
	}
	return false, 0, nil
}

func optDedupRet(instrs []Instruction, i int) (bool, int, []Instruction) {
	if instrs[i].Type == Ret && i > 0 && instrs[i-1].Type == Ret {
		return true, 1, nil // skip duplicate Ret
	}
	return false, 0, nil
}
