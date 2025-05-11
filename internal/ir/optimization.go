package ir

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"strings"
)

type OptFunc func(instrs []interfaces.Instruction, i int) (matched bool, consumed int, replacement []interfaces.Instruction)

func OptimizeFunctionBody(f interfaces.Function) interfaces.Function {
	optimized := make([]interfaces.Instruction, 0)

	opts := []OptFunc{
		optOmitAfterRet,
		optSkipMacroSet,
		optCollapseCopyChain,
		optSetCopyToSet,
		optDedupRet,
	}

	i := 0
	for i < f.GetCode().Len() {
		applied := false

		for _, opt := range opts {
			instructions := f.GetCode().GetInstructions()
			if matched, consumed, repl := opt(instructions, i); matched {
				optimized = append(optimized, repl...)
				i += consumed
				applied = true
				break
			}
		}

		if !applied {
			optimized = append(optimized, f.GetCode().GetInstructions()[i])
			i++
		}
	}

	f.GetCode().SetInstructions(optimized)
	return f
}

func isMacroSetPattern(setInst interfaces.Instruction) bool {
	if setInst.GetType() != Set || len(setInst.GetArgs()) != 3 {
		return false
	}
	path := setInst.GetArgs()[1]
	val := setInst.GetArgs()[2]

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

func optOmitAfterRet(instrs []interfaces.Instruction, i int) (bool, int, []interfaces.Instruction) {
	if i == 0 || instrs[i].GetType() != Ret {
		return false, 0, nil
	}
	if instrs[i-1].GetType() == Ret && i < len(instrs) {
		return true, len(instrs) - i, nil // skip all instructions after Ret
	}
	return false, 0, nil
}

func optSkipMacroSet(instrs []interfaces.Instruction, i int) (bool, int, []interfaces.Instruction) {
	curr := instrs[i]
	if curr.GetType() == Set && isMacroSetPattern(curr) {
		return true, 1, []interfaces.Instruction{curr}
	}
	return false, 0, nil
}

func optCollapseCopyChain(instrs []interfaces.Instruction, i int) (bool, int, []interfaces.Instruction) {
	if instrs[i].GetType() != Copy || len(instrs[i].GetArgs()) != 4 {
		return false, 0, nil
	}

	srcStorage, srcPath := instrs[i].GetArgs()[0], instrs[i].GetArgs()[1]
	curDstStorage, curDstPath := instrs[i].GetArgs()[2], instrs[i].GetArgs()[3]

	intermediate := [][2]string{{curDstStorage, curDstPath}}

	chainEnd := i + 1
	for chainEnd < len(instrs) {
		next := instrs[chainEnd]
		if next.GetType() != Copy || len(next.GetArgs()) != 4 {
			break
		}
		if !sameLocation(curDstStorage, curDstPath, next.GetArgs()[0], next.GetArgs()[1]) {
			break
		}
		curDstStorage, curDstPath = next.GetArgs()[2], next.GetArgs()[3]
		intermediate = append(intermediate, [2]string{curDstStorage, curDstPath})
		chainEnd++
	}

	if chainEnd > i+1 {
		// Ensure none of the intermediate destinations are used later
		for _, loc := range intermediate[:len(intermediate)-1] {
			if isValueUsedLater(instrs, chainEnd, loc[0], loc[1]) {
				return false, 0, nil
			}
		}
		return true, chainEnd - i, []interfaces.Instruction{Instruction{
			Type: Copy,
			Args: []string{srcStorage, srcPath, curDstStorage, curDstPath},
		}}
	}
	return false, 0, nil
}

func optSetCopyToSet(instrs []interfaces.Instruction, i int) (bool, int, []interfaces.Instruction) {
	if i+1 >= len(instrs) {
		return false, 0, nil
	}

	curr, next := instrs[i], instrs[i+1]
	if curr.GetType() == Set && len(curr.GetArgs()) == 3 &&
		next.GetType() == Copy && len(next.GetArgs()) == 4 &&
		sameLocation(curr.GetArgs()[0], curr.GetArgs()[1], next.GetArgs()[0], next.GetArgs()[1]) {

		srcStorage, srcPath := curr.GetArgs()[0], curr.GetArgs()[1]

		// Check if value is used later
		if isValueUsedLater(instrs, i+2, srcStorage, srcPath) {
			return false, 0, nil
		}

		return true, 2, []interfaces.Instruction{Instruction{
			Type: Set,
			Args: []string{next.GetArgs()[2], next.GetArgs()[3], curr.GetArgs()[2]},
		}}
	}
	return false, 0, nil
}

func optDedupRet(instrs []interfaces.Instruction, i int) (bool, int, []interfaces.Instruction) {
	if instrs[i].GetType() == Ret && i > 0 && instrs[i-1].GetType() == Ret {
		return true, 1, nil // skip duplicate Ret
	}
	return false, 0, nil
}

func isValueUsedLater(instrs []interfaces.Instruction, start int, storage, path string) bool {
	for _, instr := range instrs[start:] {
		args := instr.GetArgs()
		switch instr.GetType() {
		case Set:
			// If this Set writes to the target, it overwrites it, so it's safe.
			if len(args) >= 2 && sameLocation(args[0], args[1], storage, path) {
				return false
			}
		case Copy:
			if len(args) == 4 {
				// If this Copy reads from the target, it's a use.
				if sameLocation(args[0], args[1], storage, path) {
					return true
				}
				// If it overwrites the target, it stops our concern.
				if sameLocation(args[2], args[3], storage, path) {
					return false
				}
			}
		case Ret, Call: // Anything that might read a value
			for _, arg := range args {
				if arg == path {
					return true
				}
			}
		default:
			for _, arg := range args {
				if arg == path {
					return true
				}
			}
		}
	}
	return false
}
