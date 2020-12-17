package main

func init() {
	addSolutions(17, problem17)
}

func problem17(ctx *problemContext) {
	p3 := newPocket3()
	p4 := newPocket4()
	scanner := ctx.scanner()
	y := int64(0)
	for scanner.scan() {
		for x, c := range scanner.text() {
			v3 := vec3{x: int64(x), y: y, z: 0}
			p3.set(v3, c == '#')
			v4 := vec4{x: int64(x), y: y, z: 0, w: 0}
			p4.set(v4, c == '#')
		}
		y++
	}
	ctx.reportLoad()

	for i := 0; i < 6; i++ {
		p3.step()
	}
	ctx.reportPart1(p3.countActive())

	for i := 0; i < 6; i++ {
		p4.step()
	}
	ctx.reportPart2(p4.countActive())
}

type pocket3 struct {
	m map[vec3]bool
}

func newPocket3() *pocket3 {
	return &pocket3{m: make(map[vec3]bool)}
}

func (p *pocket3) set(v vec3, b bool) {
	p.m[v] = b
}

func (p *pocket3) activeNeighbors(v vec3) int {
	var n int
	for _, neighbor := range v.neighbors() {
		if p.m[neighbor] {
			n++
		}
	}
	return n
}

func (p *pocket3) step() {
	m := make(map[vec3]bool, len(p.m)*2)
	for v, active := range p.m {
		m[v] = false
		if active {
			for _, neighbor := range v.neighbors() {
				m[neighbor] = false
			}
		}
	}
	for v := range m {
		active := false
		na := p.activeNeighbors(v)
		if p.m[v] {
			active = na == 2 || na == 3
		} else {
			active = na == 3
		}
		m[v] = active
	}
	p.m = m
}

func (p *pocket3) countActive() int {
	var n int
	for _, active := range p.m {
		if active {
			n++
		}
	}
	return n
}

type pocket4 struct {
	m map[vec4]bool
}

func newPocket4() *pocket4 {
	return &pocket4{m: make(map[vec4]bool)}
}

func (p *pocket4) set(v vec4, b bool) {
	p.m[v] = b
}

func (p *pocket4) activeNeighbors(v vec4) int {
	var n int
	for _, neighbor := range v.neighbors() {
		if p.m[neighbor] {
			n++
		}
	}
	return n
}

func (p *pocket4) step() {
	m := make(map[vec4]bool, len(p.m)*3)
	for v, active := range p.m {
		m[v] = false
		if active {
			for _, neighbor := range v.neighbors() {
				m[neighbor] = false
			}
		}
	}
	for v := range m {
		active := false
		na := p.activeNeighbors(v)
		if p.m[v] {
			active = na == 2 || na == 3
		} else {
			active = na == 3
		}
		m[v] = active
	}
	p.m = m
}

func (p *pocket4) countActive() int {
	var n int
	for _, active := range p.m {
		if active {
			n++
		}
	}
	return n
}
