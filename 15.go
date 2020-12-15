package main

import (
	"io/ioutil"
	"log"
	"strings"
)

func init() {
	addSolutions(15, problem15)
}

func problem15(ctx *problemContext) {
	b, err := ioutil.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	var ns []uint32
	for _, s := range strings.Split(strings.TrimSpace(string(b)), ",") {
		ns = append(ns, uint32(parseUint(s, 10, 32)))
	}
	ctx.reportLoad()

	var m memGame
	var n uint32
	for i := uint32(0); i < 30e6; i++ {
		prev := n
		if i < uint32(len(ns)) {
			n = ns[i]
		} else {
			if i0, ok := m.get(n); ok {
				n = i - 1 - i0
			} else {
				n = 0
			}
		}
		if i > 0 {
			m.set(prev, i-1)
		}
		if i == 2019 {
			ctx.reportPart1(n)
		} else if i == 30e6-1 {
			ctx.reportPart2(n)
		}
	}
}

type memGame []uint32

func (m *memGame) get(n uint32) (uint32, bool) {
	if n >= uint32(len(*m)) {
		return 0, false
	}
	prev := (*m)[n]
	if prev == 0 {
		return 0, false
	}
	return prev - 1, true
}

func (m *memGame) set(n, step uint32) {
	if n >= uint32(len(*m)) {
		m1 := make([]uint32, (n+1)*2)
		copy(m1, *m)
		*m = m1
	}
	(*m)[n] = step + 1
}
