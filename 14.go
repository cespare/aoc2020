package main

import (
	"math/bits"
	"strings"
)

func init() {
	addSolutions(14, problem14)
}

func problem14(ctx *problemContext) {
	var code []interface{}
	scanner := ctx.scanner()
	for scanner.scan() {
		line := scanner.text()
		if strings.HasPrefix(line, "mask") {
			code = append(code, parseMask36(line))
		} else {
			code = append(code, parseMem36Assign(line))
		}
	}
	ctx.reportLoad()

	prog := newMem36Program(code)
	prog.run1()
	ctx.reportPart1(prog.sum())

	prog = newMem36Program(code)
	prog.run2()
	ctx.reportPart2(prog.sum())
}

type mem36 uint64

type mask36 struct {
	// As a "part 2" mask, x is unused.
	x    mem36
	zero mem36
	one  mem36
}

func parseMask36(s string) mask36 {
	s = strings.TrimPrefix(s, "mask = ")
	var mask mask36
	for i := 0; i < len(s); i++ {
		switch s[len(s)-i-1] {
		case 'X':
			mask.x |= 1 << i
		case '0':
			mask.zero |= 1 << i
		case '1':
			mask.one |= 1 << i
		}
	}
	return mask
}

func (m mask36) expand2() []mask36 {
	var masks []mask36
	var emit func(mask36, mem36)
	emit = func(m mask36, x mem36) {
		if x == 0 {
			masks = append(masks, m)
			return
		}
		i := bits.TrailingZeros64(uint64(x))
		x &^= (1 << i)
		m0 := m
		m0.zero |= (1 << i)
		emit(m0, x)
		m1 := m
		m1.one |= (1 << i)
		emit(m1, x)
	}
	emit(mask36{}, m.x)
	return masks
}

func (m mem36) applyMask1(mask mask36) mem36 {
	m &^= mask.zero
	m |= mask.one
	return m
}

func (m mem36) applyMask2(mask1, mask2 mask36) mem36 {
	m |= mask1.one
	return m.applyMask1(mask2)
}

type mem36Assign struct {
	addr mem36
	val  mem36
}

func parseMem36Assign(s string) mem36Assign {
	parts := strings.SplitN(s, " = ", 2)
	if len(parts) != 2 {
		panic(s)
	}
	addrStr := strings.TrimSuffix(strings.TrimPrefix(parts[0], "mem["), "]")
	return mem36Assign{
		addr: mem36(parseUint(addrStr, 10, 36)),
		val:  mem36(parseUint(parts[1], 10, 36)),
	}

}

type mem36Program struct {
	mask1  mask36
	masks2 []mask36
	mem    map[mem36]mem36
	code   []interface{}
}

func newMem36Program(code []interface{}) *mem36Program {
	return &mem36Program{
		mem:  make(map[mem36]mem36),
		code: code,
	}
}

func (p *mem36Program) run1() {
	for _, op := range p.code {
		if mask, ok := op.(mask36); ok {
			p.mask1 = mask
			continue
		}
		assn := op.(mem36Assign)
		p.mem[assn.addr] = assn.val.applyMask1(p.mask1)
	}
}

func (p *mem36Program) run2() {
	for _, op := range p.code {
		if mask, ok := op.(mask36); ok {
			p.mask1 = mask
			p.masks2 = mask.expand2()
			continue
		}
		assn := op.(mem36Assign)
		for _, mask2 := range p.masks2 {
			addr := assn.addr.applyMask2(p.mask1, mask2)
			p.mem[addr] = assn.val
		}
	}
}

func (p *mem36Program) sum() int64 {
	var sum int64
	for _, v := range p.mem {
		sum += int64(v)
	}
	return sum
}
