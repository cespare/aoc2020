package main

import (
	"bytes"
)

func init() {
	addSolutions(11, problem11)
}

func problem11(ctx *problemContext) {
	var s0 seats
	scanner := ctx.scanner()
	for scanner.scan() {
		s0 = append(s0, []byte(scanner.text()))
	}
	ctx.reportLoad()

	s1 := s0
	for {
		s2 := s1.iterate1()
		if s2.same(s1) {
			ctx.reportPart1(s2.occupied())
			break
		}
		s1 = s2
	}

	s1 = s0
	for {
		s2 := s1.iterate2()
		if s2.same(s1) {
			ctx.reportPart2(s2.occupied())
			break
		}
		s1 = s2
	}
}

type seats [][]byte

func (s seats) same(s1 seats) bool {
	for y, row := range s {
		if !bytes.Equal(row, s1[y]) {
			return false
		}
	}
	return true
}

func (s seats) cols() int { return len(s[0]) }

func (s seats) rows() int { return len(s) }

func (s seats) iterate1() seats {
	adj := make([][]int, s.rows())
	for y := range s {
		adj[y] = make([]int, s.cols())
	}
	for y, row := range s {
		for x := range row {
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					x0 := x + dx
					y0 := y + dy
					switch {
					case dx == 0 && dy == 0:
					case x0 < 0 || x0 >= s.cols():
					case y0 < 0 || y0 >= s.rows():
					case s[y0][x0] == '#':
						adj[y][x]++
					}
				}
			}
		}
	}
	s1 := make(seats, s.rows())
	for y, row := range s {
		s1[y] = make([]byte, s.cols())
		for x, cell := range row {
			switch {
			case cell == 'L' && adj[y][x] == 0:
				cell = '#'
			case cell == '#' && adj[y][x] >= 4:
				cell = 'L'
			}
			s1[y][x] = cell
		}
	}
	return s1
}

func (s seats) iterate2() seats {
	vis := make([][]int, s.rows())
	for y := range s {
		vis[y] = make([]int, s.cols())
	}
	for y, row := range s {
		for x := range row {
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					if dx == 0 && dy == 0 {
						continue
					}
					x0 := x
					y0 := y
				loop:
					x0 += dx
					y0 += dy
					switch {
					case x0 < 0 || x0 >= s.cols():
					case y0 < 0 || y0 >= s.rows():
					case s[y0][x0] == '#':
						vis[y][x]++
					case s[y0][x0] == 'L':
					default:
						goto loop
					}
				}
			}
		}
	}
	s1 := make(seats, s.rows())
	for y, row := range s {
		s1[y] = make([]byte, s.cols())
		for x, cell := range row {
			switch {
			case cell == 'L' && vis[y][x] == 0:
				cell = '#'
			case cell == '#' && vis[y][x] >= 5:
				cell = 'L'
			}
			s1[y][x] = cell
		}
	}
	return s1
}

func (s seats) occupied() int {
	var occ int
	for _, row := range s {
		occ += bytes.Count(row, []byte("#"))
	}
	return occ
}
