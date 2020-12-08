package main

import (
	"strconv"
	"strings"
)

func init() {
	addSolutions(8, problem8)
}

func problem8(ctx *problemContext) {
	var insts []bootInst
	scanner := ctx.scanner()
	for scanner.scan() {
		insts = append(insts, parseBootInst(scanner.text()))
	}
	ctx.reportLoad()

	m := &bootMachine{code: insts}
	m.run()
	ctx.reportPart1(m.acc)

	for i, inst := range insts {
		switch inst.op {
		case nop:
			insts[i].op = jmp
			if m.run() {
				ctx.reportPart2(m.acc)
				return
			}
		case jmp:
			insts[i].op = nop
			if m.run() {
				ctx.reportPart2(m.acc)
				return
			}
		}
		insts[i] = inst
		m.reset()
	}
}

type bootOp int

const (
	nop bootOp = iota
	acc
	jmp
)

type bootInst struct {
	op bootOp
	v  int
}

func parseBootInst(s string) bootInst {
	var inst bootInst
	parts := strings.SplitN(s, " ", 2)
	switch parts[0] {
	case "nop":
		inst.op = nop
	case "acc":
		inst.op = acc
	case "jmp":
		inst.op = jmp
	default:
		panic("bad opcode")
	}
	n, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	inst.v = n
	return inst
}

type bootMachine struct {
	code []bootInst
	pc   int
	acc  int64
}

func (m *bootMachine) run() (terminated bool) {
	seen := make([]bool, len(m.code))
	for {
		if m.pc == len(m.code) {
			return true
		}
		if seen[m.pc] {
			return false
		}
		seen[m.pc] = true
		inst := m.code[m.pc]
		switch inst.op {
		case nop:
		case acc:
			m.acc += int64(inst.v)
		case jmp:
			m.pc += inst.v - 1
		}
		m.pc++
	}
}

func (m *bootMachine) reset() {
	m.pc = 0
	m.acc = 0
}
