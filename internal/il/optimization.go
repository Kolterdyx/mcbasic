package il

import (
	"strings"
)

func OptimizeFunctionBody(f Function) Function {
	return peepholeOptimize(f)
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

func peepholeOptimize(f Function) Function {
	optimized := make([]Instruction, 0)

	i := 0
	for i < len(f.Instructions) {
		curr := f.Instructions[i]

		if curr.Type == Set && isMacroSetPattern(curr) {
			// Skip macro var copy-in from $(...) to vars.X
			optimized = append(optimized, curr)
			i++
			continue
		}

		// Try to collapse copy chains
		if curr.Type == Copy && len(curr.Args) == 4 {
			srcStorage, srcPath := curr.Args[0], curr.Args[1]
			curDstStorage, curDstPath := curr.Args[2], curr.Args[3]

			// Begin chain
			chainEnd := i + 1

			for chainEnd < len(f.Instructions) {
				next := f.Instructions[chainEnd]
				if next.Type != Copy || len(next.Args) != 4 {
					break
				}

				// Check if next source matches current destination
				if !sameLocation(curDstStorage, curDstPath, next.Args[0], next.Args[1]) {
					break
				}

				// Check for macros
				if isMacroPath(next.Args[1]) || isMacroPath(next.Args[3]) {
					break
				}

				// Extend chain
				curDstStorage, curDstPath = next.Args[2], next.Args[3]
				chainEnd++
			}

			// Collapse to a single copy if chain length > 1
			if chainEnd > i+1 {
				optimized = append(optimized, Instruction{
					Type: Copy,
					Args: []string{srcStorage, srcPath, curDstStorage, curDstPath},
				})
				i = chainEnd
				continue
			}
		}

		// Pattern: set + copy → set
		if i+1 < len(f.Instructions) &&
			curr.Type == Set && len(curr.Args) == 3 &&
			f.Instructions[i+1].Type == Copy && len(f.Instructions[i+1].Args) == 4 {
			next := f.Instructions[i+1]

			if sameLocation(curr.Args[0], curr.Args[1], next.Args[0], next.Args[1]) &&
				!isMacroPath(curr.Args[1]) && !isMacroPath(next.Args[3]) {

				optimized = append(optimized, Instruction{
					Type: Set,
					Args: []string{next.Args[2], next.Args[3], curr.Args[2]},
				})
				i += 2
				continue
			}
		}

		if curr.Type == Ret && len(optimized) > 0 && optimized[len(optimized)-1].Type == Ret {
			i++
			continue
		}

		// Otherwise keep instruction as-is
		optimized = append(optimized, curr)
		i++
	}

	f.Instructions = optimized
	return f
}
