package main

import "math"

func init() {
	addSolutions(9, problem9)
}

func problem9(ctx *problemContext) {
	var ns []int64
	scanner := ctx.scanner()
	for scanner.scan() {
		ns = append(ns, scanner.int64())
	}
	ctx.reportLoad()

	fail := xmasFail(ns, 25)
	ctx.reportPart1(fail)

	ctx.reportPart2(xmasContig(ns, fail))
}

func xmasFail(ns []int64, p int) int64 {
	pairs := make(map[int64]int)
	for i := 0; i < p; i++ {
		for j := i + 1; j < p; j++ {
			n0, n1 := ns[i], ns[j]
			if n0 != n1 {
				pairs[n0+n1]++
			}
		}
	}
	for i := p; i < len(ns); i++ {
		n2 := ns[i]
		if pairs[n2] == 0 {
			return n2
		}
		n0 := ns[i-p]
		for j := i - p + 1; j < i; j++ {
			n1 := ns[j]
			pairs[n0+n1]--
			pairs[n1+n2]++
		}
	}
	panic("fail")
}

func xmasContig(ns []int64, targ int64) int64 {
	i := 0
	j := 1
	sum := ns[0]
	for {
		switch {
		case sum < targ:
			sum += ns[j]
			j++
		case sum > targ:
			sum -= ns[i]
			i++
		default:
			min := int64(math.MaxInt64)
			max := int64(0)
			for k := i; k < j; k++ {
				n := ns[k]
				if n < min {
					min = n
				}
				if n > max {
					max = n
				}
			}
			return min + max
		}
	}
}
