package main

import (
	"math/bits"
	"strings"
)

func init() {
	addSolutions(6, problem6)
}

func problem6(ctx *problemContext) {
	scanner := ctx.scanner()
	scanner.s.Split(splitGroups)
	var any, every int
	for scanner.scan() {
		var anySet, everySet uint32
		everySet = (1 << 27) - 1
		for _, line := range strings.Split(strings.TrimSpace(scanner.text()), "\n") {
			var set uint32
			for i := 0; i < len(line); i++ {
				set |= 1 << (line[i] - 'a')
			}
			anySet |= set
			everySet &= set
		}
		any += bits.OnesCount32(anySet)
		every += bits.OnesCount32(everySet)
	}
	ctx.reportLoad()

	ctx.reportPart1(any)
	ctx.reportPart2(every)
}
