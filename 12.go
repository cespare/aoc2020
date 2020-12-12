package main

import (
	"log"
	"strconv"
)

func init() {
	addSolutions(12, problem12)
}

func problem12(ctx *problemContext) {
	var dirs []shipDir
	scanner := ctx.scanner()
	for scanner.scan() {
		dirs = append(dirs, parseShipDir(scanner.text()))
	}
	ctx.reportLoad()

	s := newShipState()
	for _, d := range dirs {
		s.do(d)
	}
	ctx.reportPart1(s.dist())

	w := newWaypointState()
	for _, d := range dirs {
		w.do(d)
	}
	ctx.reportPart2(w.dist())
}

type shipDir struct {
	c byte
	d int64
}

func parseShipDir(s string) shipDir {
	var sd shipDir
	sd.c = s[0]
	var err error
	sd.d, err = strconv.ParseInt(s[1:], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return sd
}

type shipState struct {
	p  vec2
	oi cardinal
}

func newShipState() *shipState {
	return &shipState{
		oi: east,
	}
}

func (s *shipState) dist() int64 {
	return s.p.mag()
}

func (s *shipState) do(d shipDir) {
	var v vec2
	switch d.c {
	case 'N':
		v = vnorth
	case 'S':
		v = vsouth
	case 'E':
		v = veast
	case 'W':
		v = vwest
	case 'L':
		if d.d%90 != 0 {
			panic(d.d)
		}
		s.oi = (s.oi + 4 - cardinal(d.d/90)) % 4
		return
	case 'R':
		if d.d%90 != 0 {
			panic(d.d)
		}
		s.oi = (s.oi + cardinal(d.d/90)) % 4
		return
	case 'F':
		v = cardinals[s.oi]
	default:
		panic(d.c)
	}
	s.p = s.p.add(v.mul(d.d))
}

type waypointState struct {
	ship vec2
	wp   vec2 // relative
}

func newWaypointState() *waypointState {
	return &waypointState{wp: vec2{10, 1}}
}

func (w *waypointState) dist() int64 {
	return w.ship.mag()
}

func (w *waypointState) do(d shipDir) {
	var v vec2
	switch d.c {
	case 'N':
		v = vnorth
	case 'S':
		v = vsouth
	case 'E':
		v = veast
	case 'W':
		v = vwest
	default:
		goto after
	}
	w.wp = w.wp.add(v.mul(d.d))
	return

after:
	var r int
	switch d.c {
	case 'L':
		r = 4 - (int(d.d / 90))
	case 'R':
		r = int(d.d / 90)
	default:
		goto forward
	}
	w.wp = w.wp.matMul(rotations[r])
	return

forward:
	if d.c != 'F' {
		panic(d.c)
	}
	w.ship = w.ship.add(w.wp.mul(d.d))
}
