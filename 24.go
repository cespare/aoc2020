package main

func init() {
	addSolutions(24, problem24)
}

func problem24(ctx *problemContext) {
	var coords [][]hexDir
	scanner := ctx.scanner()
	for scanner.scan() {
		line := scanner.text()
		var dirs []hexDir
		for line != "" {
			var dir hexDir
			dir, line = parseHexDir(line)
			dirs = append(dirs, dir)
		}
		coords = append(coords, dirs)
	}
	ctx.reportLoad()

	grid := make(hexGrid)
	for _, coord := range coords {
		var pos hexPos
		for _, dir := range coord {
			pos = pos.move(dir)
		}
		grid.flip(pos)
	}
	ctx.reportPart1(len(grid))

	for i := 0; i < 100; i++ {
		grid = grid.step()
	}
	ctx.reportPart2(len(grid))
}

type hexDir int

const (
	hexE hexDir = iota
	hexSE
	hexSW
	hexW
	hexNW
	hexNE
)

func parseHexDir(s string) (hexDir, string) {
	for _, d := range []struct {
		prefix string
		dir    hexDir
	}{
		{"e", hexE},
		{"se", hexSE},
		{"sw", hexSW},
		{"w", hexW},
		{"nw", hexNW},
		{"ne", hexNE},
	} {
		if s, ok := trimPrefix(s, d.prefix); ok {
			return d.dir, s
		}
	}
	panic("unreached")
}

type hexPos struct {
	x int
	y int
}

func (p hexPos) move(dir hexDir) hexPos {
	switch dir {
	case hexE:
		return hexPos{x: p.x + 1, y: p.y}
	case hexSE:
		return hexPos{x: p.x, y: p.y + 1}
	case hexSW:
		return hexPos{x: p.x - 1, y: p.y + 1}
	case hexW:
		return hexPos{x: p.x - 1, y: p.y}
	case hexNW:
		return hexPos{x: p.x, y: p.y - 1}
	case hexNE:
		return hexPos{x: p.x + 1, y: p.y - 1}
	default:
		panic("unreached")
	}
}

func (p hexPos) neighbors() []hexPos {
	return []hexPos{
		{x: p.x + 1, y: p.y},
		{x: p.x, y: p.y + 1},
		{x: p.x - 1, y: p.y + 1},
		{x: p.x - 1, y: p.y},
		{x: p.x, y: p.y - 1},
		{x: p.x + 1, y: p.y - 1},
	}
}

type hexGrid map[hexPos]struct{} // only store black

func (g hexGrid) flip(p hexPos) {
	if _, ok := g[p]; ok {
		delete(g, p)
	} else {
		g[p] = struct{}{}
	}
}

func (g hexGrid) step() hexGrid {
	next := make(hexGrid)
	// Rough cut, then refine.
	for p := range g {
		next[p] = struct{}{}
		for _, n := range p.neighbors() {
			next[n] = struct{}{}
		}
	}
	for p := range next {
		var bn int
		for _, n := range p.neighbors() {
			if _, ok := g[n]; ok {
				bn++
			}
		}
		if _, ok := g[p]; ok {
			if bn > 0 && bn <= 2 {
				next[p] = struct{}{}
			} else {
				delete(next, p)
			}
		} else {
			if bn == 2 {
				next[p] = struct{}{}
			} else {
				delete(next, p)
			}
		}
	}
	return next
}
