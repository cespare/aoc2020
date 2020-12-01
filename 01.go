package main

func init() {
	addSolutions(1, problem1)
}

func problem1(ctx *problemContext) {
	entries := make(map[int64]struct{})
	scanner := ctx.scanner()
	for scanner.scan() {
		entries[scanner.int64()] = struct{}{}
	}
	ctx.reportLoad()

	for n := range entries {
		m := 2020 - n
		if _, ok := entries[m]; ok {
			ctx.reportPart1(n * m)
			break
		}
	}

part2:
	for x := range entries {
		for y := range entries {
			z := 2020 - x - y
			if _, ok := entries[z]; ok {
				ctx.reportPart2(x * y * z)
				break part2
			}
		}
	}
}
