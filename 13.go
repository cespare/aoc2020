package main

import (
	"strings"
)

func init() {
	addSolutions(13, problem13)
}

func problem13(ctx *problemContext) {
	var start int64
	var busSchedule []int64
	scanner := ctx.scanner()
	var i int
	for scanner.scan() {
		if i == 0 {
			start = scanner.int64()
		} else if i == 1 {
			for _, s := range strings.Split(scanner.text(), ",") {
				if s == "x" {
					busSchedule = append(busSchedule, -1)
				} else {
					busSchedule = append(busSchedule, parseInt(s, 10, 64))
				}
			}
		}
		i++
	}
	ctx.reportLoad()

	min := int64(-1)
	var minID int64
	for _, n := range busSchedule {
		if n == -1 {
			continue
		}
		d := ((start+n-1)/n)*n - start
		if min < 0 || d < min {
			min = d
			minID = n
		}
	}
	ctx.reportPart1(min * minID)

	t := busSchedule[0]
	m := t
	for i, d := range busSchedule {
		if i == 0 || d == -1 {
			continue
		}
		t = busCombine(t, m, d, int64(i)%d)
		m *= d
	}
	ctx.reportPart2(t)
}

func busCombine(t, m, d, r int64) int64 {
	for c := t; ; c += m {
		if c < d {
			continue
		}
		if (c+r)%d == 0 {
			return c
		}
	}
}
