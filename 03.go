package main

import "log"

func init() {
	addSolutions(3, problem3)
}

func problem3(ctx *problemContext) {
	var m treeMap
	scanner := ctx.scanner()
	for scanner.scan() {
		m.addRow(scanner.text())
	}
	ctx.reportLoad()

	ctx.reportPart1(m.countForSlope(3, 1))

	p := 1
	p *= m.countForSlope(1, 1)
	p *= m.countForSlope(3, 1)
	p *= m.countForSlope(5, 1)
	p *= m.countForSlope(7, 1)
	p *= m.countForSlope(1, 2)
	ctx.reportPart2(p)
}

type treeMap struct {
	rows [][]bool
}

func (m *treeMap) addRow(s string) {
	if len(m.rows) > 0 && len(s) != len(m.rows[0]) {
		log.Fatal("mismatched row sizes")
	}
	row := make([]bool, len(s))
	for i := range row {
		row[i] = s[i] != '.'
	}
	m.rows = append(m.rows, row)
}

func (m *treeMap) countForSlope(dx, dy int) int {
	var n int
	var x, y int
	for y < len(m.rows) {
		x = x % len(m.rows[0])
		if m.rows[y][x] {
			n++
		}
		x += dx
		y += dy
	}
	return n
}
