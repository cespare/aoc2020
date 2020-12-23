package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func init() {
	addSolutions(23, problem23)
}

func problem23(ctx *problemContext) {
	b, err := ioutil.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	ctx.reportLoad()

	ns := make([]int, len(b)-1)
	for i, c := range bytes.TrimSpace(b) {
		ns[i] = int(c) - '0'
	}
	cups := newCups(append([]int(nil), ns...))
	for i := 0; i < 100; i++ {
		cups.step()
	}
	ctx.reportPart1(cups.order())

	millionCups := make([]int, 1e6)
	copy(millionCups, ns)
	var n int
	for i := range millionCups {
		if i < len(ns) {
			if ns[i] > n {
				n = ns[i]
			}
		} else {
			n++
			millionCups[i] = n
		}
	}
	cups = newCups(millionCups)
	for i := 0; i < 10e6; i++ {
		cups.step()
	}
	x, y := cups.nextAfter1()
	ctx.reportPart2(x * y)
}

type cups struct {
	next    []int
	cur     int
	scratch [3]int
}

func newCups(ns []int) *cups {
	c := &cups{
		next: make([]int, len(ns)+1),
		cur:  ns[0],
	}
	c.next[0] = -1
	for i, v := range ns {
		c.next[v] = ns[(i+1)%len(ns)]
	}
	return c
}

func (c *cups) step() {
	x := c.next[c.cur]
	for i := 0; i < 3; i++ {
		c.scratch[i] = x
		x = c.next[x]
	}
	c.next[c.cur] = x
	dest := c.cur
destLoop:
	for {
		dest--
		if dest == 0 {
			dest = len(c.next) - 1
		}
		for _, v := range c.scratch {
			if dest == v {
				continue destLoop
			}
		}
		break
	}
	after := c.next[dest]
	c.next[dest] = c.scratch[0]
	c.next[c.scratch[2]] = after
	c.cur = c.next[c.cur]
}

func (c *cups) order() string {
	var b strings.Builder
	for x := c.next[1]; x != 1; x = c.next[x] {
		fmt.Fprint(&b, x)
	}
	return b.String()
}

func (c *cups) nextAfter1() (int, int) {
	x := c.next[1]
	y := c.next[x]
	return x, y
}
