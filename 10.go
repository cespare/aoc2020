package main

import (
	"sort"
)

func init() {
	addSolutions(10, problem10)
}

func problem10(ctx *problemContext) {
	jolts := []int64{0}
	scanner := ctx.scanner()
	for scanner.scan() {
		jolts = append(jolts, scanner.int64())
	}
	ctx.reportLoad()

	sort.Slice(jolts, func(i, j int) bool { return jolts[i] < jolts[j] })
	jolts = append(jolts, jolts[len(jolts)-1]+3)

	ctx.reportPart1(joltDelta(jolts))
	ctx.reportPart2(joltCombos(jolts))
}

func joltDelta(jolts []int64) int64 {
	var d [4]int64
	for i, jolt := range jolts[:len(jolts)-1] {
		d[jolts[i+1]-jolt]++
	}
	return d[1] * d[3]
}

func joltCombos(jolts []int64) int64 {
	counts := make([]int64, len(jolts))
	counts[len(jolts)-1] = 1
	for i := len(jolts) - 2; i >= 0; i-- {
		jolt := jolts[i]
		for j := i + 1; j < len(jolts) && jolts[j]-jolt <= 3; j++ {
			counts[i] += counts[j]
		}
	}
	return counts[0]
}
