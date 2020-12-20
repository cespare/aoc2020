package main

import (
	"bytes"
	"strings"
)

func init() {
	addSolutions(20, problem20)
}

func problem20(ctx *problemContext) {
	var tiles []*imageTile
	scanner := ctx.scanner()
	scanner.s.Split(splitGroups)
	for scanner.scan() {
		group := scanner.text()
		tiles = append(tiles, parseImageTile(group))
	}
	ctx.reportLoad()

	edges := make(map[string][]tileAndOrientation)
	for _, tile := range tiles {
		for _, to := range allTileOrientations {
			edge := tile.edge(to)
			tao := append(edges[edge], tileAndOrientation{tile, to})
			edges[edge] = tao
		}
	}

	var topLeft tileAndOrientation
	corners := make(map[*imageTile]struct{})
	for _, tile := range tiles {
		var singles, flippedSingles int
		for _, to := range allTileOrientations {
			if len(edges[tile.edge(to)]) == 1 {
				if to.flipped {
					flippedSingles++
				} else {
					singles++
				}
			}
		}
		if singles == 2 || flippedSingles == 2 {
			if singles < 2 || flippedSingles < 2 {
				panic("unexpected")
			}
			corners[tile] = struct{}{}
			if topLeft.tile == nil {
				topLeft.tile = tile
				flipped := true // choose to match example
				to0 := tileOrientation{0, flipped}
				to1 := tileOrientation{1, flipped}
				to2 := tileOrientation{2, flipped}
				to3 := tileOrientation{3, flipped}
				single0 := len(edges[tile.edge(to0)]) == 1
				single1 := len(edges[tile.edge(to1)]) == 1
				single2 := len(edges[tile.edge(to2)]) == 1
				single3 := len(edges[tile.edge(to3)]) == 1
				switch {
				case single0 && single1:
					topLeft.to = to1
				case single1 && single2:
					topLeft.to = to2
				case single2 && single3:
					topLeft.to = to3
				case single3 && single0:
					topLeft.to = to0
				default:
					panic("unexpected")
				}
			}
		}
		if singles > 2 || flippedSingles > 2 {
			panic(tile.id)
		}
	}

	if len(corners) != 4 {
		panic("bad")
	}

	part1 := int64(1)
	for tile := range corners {
		part1 *= tile.id
	}
	ctx.reportPart1(part1)

	state := &arrangementState{
		edges: edges,
		placed: map[vec2]tileAndOrientation{
			{0, 0}: topLeft,
		},
		corners: corners,
		used:    map[*imageTile]struct{}{topLeft.tile: {}},
	}
	if !state.arrange(vec2{1, 0}) {
		panic("no solution")
	}

	image := make([][]byte, state.height*8)
	for y := range image {
		row := make([]byte, state.width*8)
		image[y] = row
	}
	for v, tao := range state.placed {
		tileImg := make([][]byte, 8)
		for y := range tileImg {
			row := make([]byte, 8)
			tileImg[y] = row
			for x := range row {
				row[x] = tao.tile.array[y+1][x+1]
			}
		}
		tileImg = reorientImage(tileImg, tao.to)
		for y := range tileImg {
			row := tileImg[y]
			for x := range row {
				v1 := v.mul(8).add(vec2{int64(x), int64(y)})
				image[v1.y][v1.x] = row[x]
			}
		}
	}

	for _, to := range allTileOrientations {
		image1 := reorientImage(image, to)
		n := countSeaMonsters(image1)
		if n == 0 {
			continue
		}
		var water int
		for _, row := range image1 {
			water += bytes.Count(row, []byte("#"))
		}
		ctx.reportPart2(water - n*len(seaMonster))
		return
	}
}

type arrangementState struct {
	edges map[string][]tileAndOrientation

	placed  map[vec2]tileAndOrientation
	width   int64
	height  int64
	corners map[*imageTile]struct{}
	used    map[*imageTile]struct{}
}

func (s *arrangementState) arrange(p vec2) bool {
	var horizCandidates []tileAndOrientation
	if p.x > 0 {
		leftNeighbor := s.placed[vec2{p.x - 1, p.y}]
		leftNeighborRightSide := tileOrientation{
			side:    (leftNeighbor.to.side + 1) % 4,
			flipped: leftNeighbor.to.flipped,
		}
		edge := reverse(leftNeighbor.tile.edge(leftNeighborRightSide))
		for _, cand := range s.candidates(edge) {
			cand.to.side = (cand.to.side + 1) % 4
			horizCandidates = append(horizCandidates, cand)
		}
		if len(horizCandidates) == 0 {
			return false
		}
	}
	var vertCandidates []tileAndOrientation
	if p.y > 0 {
		topNeighbor := s.placed[vec2{p.x, p.y - 1}]
		topNeighborBottomSide := tileOrientation{
			side:    (topNeighbor.to.side + 2) % 4,
			flipped: topNeighbor.to.flipped,
		}
		edge := reverse(topNeighbor.tile.edge(topNeighborBottomSide))
		for _, cand := range s.candidates(edge) {
			vertCandidates = append(vertCandidates, cand)
		}
		if len(vertCandidates) == 0 {
			return false
		}
	}
	var candidates []tileAndOrientation
	switch {
	case len(horizCandidates) == 0:
		candidates = vertCandidates
	case len(vertCandidates) == 0:
		candidates = horizCandidates
	default:
		for _, hc := range horizCandidates {
			for _, vc := range vertCandidates {
				if hc == vc {
					candidates = append(candidates, hc)
					break
				}
			}
		}
		if len(candidates) == 0 {
			return false
		}
	}

	for _, cand := range candidates {
		// Apply our candidate.
		s.placed[p] = cand
		s.used[cand.tile] = struct{}{}
		atRight := p.x == s.width-1
		if _, ok := s.corners[cand.tile]; ok {
			if p.y == 0 {
				s.width = p.x + 1
				atRight = true
			}
			if p.x == 0 {
				s.height = p.y + 1
			}
		}
		if p.x == s.width-1 && p.y == s.height-1 {
			return true // done
		}
		next := vec2{p.x + 1, p.y}
		if atRight {
			next = vec2{0, p.y + 1}
		}
		if s.arrange(next) {
			return true
		}

		// Undo.
		delete(s.placed, p)
		delete(s.used, cand.tile)
	}
	return false
}

func (s *arrangementState) candidates(edge string) []tileAndOrientation {
	var candidates []tileAndOrientation
	for _, tao := range s.edges[edge] {
		if _, ok := s.used[tao.tile]; ok {
			continue
		}
		candidates = append(candidates, tao)
	}
	return candidates
}

type imageTile struct {
	id    int64
	array []string
}

type tileAndOrientation struct {
	tile *imageTile
	to   tileOrientation
}

func parseImageTile(s string) *imageTile {
	parts := strings.SplitN(s, ":\n", 2)
	id := parseInt(strings.TrimPrefix(parts[0], "Tile "), 10, 64)
	array := strings.Split(parts[1], "\n")
	if len(array) != 10 {
		panic("malformed image")
	}
	return &imageTile{id, array}
}

func (t *imageTile) edge(to tileOrientation) string {
	var e string
	switch to.side {
	case 0:
		e = t.array[0]
		if to.flipped {
			return reverse(e)
		}
		return e
	case 1:
		var b strings.Builder
		if to.flipped {
			for _, row := range t.array {
				b.WriteByte(row[0])
			}
		} else {
			for _, row := range t.array {
				b.WriteByte(row[len(t.array[0])-1])
			}
		}
		return b.String()
	case 2:
		e = t.array[len(t.array)-1]
		if !to.flipped {
			return reverse(e)
		}
		return e
	case 3:
		var b strings.Builder
		if to.flipped {
			for y := len(t.array) - 1; y >= 0; y-- {
				b.WriteByte(t.array[y][len(t.array[0])-1])
			}
		} else {
			for y := len(t.array) - 1; y >= 0; y-- {
				b.WriteByte(t.array[y][0])
			}
		}
		return b.String()
	default:
		panic("unreached")
	}
}

func reverse(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := len(s) - 1; i >= 0; i-- {
		b.WriteByte(s[i])
	}
	return b.String()
}

type tileOrientation struct {
	side    int  // 0 - 3, starting at top, clockwise
	flipped bool // whether it's flipped left to right
}

// reorientImage flip/rotates the square image such the the edge indicated by to
// is at the top.
func reorientImage(img [][]byte, to tileOrientation) [][]byte {
	var start vec2
	var minor, major vec2
	switch to {
	case tileOrientation{0, false}:
		start = vec2{0, 0}
		minor = vec2{1, 0}
		major = vec2{0, 1}
	case tileOrientation{0, true}:
		start = vec2{int64(len(img)) - 1, 0}
		minor = vec2{-1, 0}
		major = vec2{0, 1}
	case tileOrientation{1, false}:
		start = vec2{int64(len(img)) - 1, 0}
		minor = vec2{0, 1}
		major = vec2{-1, 0}
	case tileOrientation{1, true}:
		start = vec2{0, 0}
		minor = vec2{0, 1}
		major = vec2{1, 0}
	case tileOrientation{2, false}:
		start = vec2{int64(len(img)) - 1, int64(len(img)) - 1}
		minor = vec2{-1, 0}
		major = vec2{0, -1}
	case tileOrientation{2, true}:
		start = vec2{0, int64(len(img)) - 1}
		minor = vec2{1, 0}
		major = vec2{0, -1}
	case tileOrientation{3, false}:
		start = vec2{0, int64(len(img)) - 1}
		minor = vec2{0, -1}
		major = vec2{1, 0}
	case tileOrientation{3, true}:
		start = vec2{int64(len(img)) - 1, int64(len(img)) - 1}
		minor = vec2{0, -1}
		major = vec2{-1, 0}
	}
	img1 := make([][]byte, len(img))
	for y := range img1 {
		img1[y] = make([]byte, len(img))
	}
	p := start
	for y := 0; y < len(img); y++ {
		p0 := p
		for x := 0; x < len(img); x++ {
			img1[y][x] = img[p.y][p.x]
			p = p.add(minor)
		}
		p = p0.add(major)
	}
	return img1
}

func countSeaMonsters(img [][]byte) int {
	var n int
	for y, row := range img {
		for x := range row {
			if seaMonsterAt(img, vec2{int64(x), int64(y)}) {
				n++
			}
		}
	}
	return n
}

func seaMonsterAt(img [][]byte, p vec2) bool {
	for _, v := range seaMonster {
		v1 := p.add(v)
		if v1.x >= int64(len(img[0])) || v1.y >= int64(len(img)) {
			return false
		}
		if img[v1.y][v1.x] != '#' {
			return false
		}
	}
	return true
}

var seaMonsterText = `
                  #
#    ##    ##    ###
 #  #  #  #  #  #
`

var seaMonster []vec2

func init() {
	s := strings.TrimLeft(seaMonsterText, "\n")
	s = strings.TrimRight(s, "\n")
	rows := strings.Split(s, "\n")
	for y, row := range rows {
		for x, c := range row {
			if c == '#' {
				seaMonster = append(seaMonster, vec2{int64(x), int64(y)})
			}
		}
	}
}

var allTileOrientations = []tileOrientation{
	{0, false},
	{0, true},
	{1, false},
	{1, true},
	{2, false},
	{2, true},
	{3, false},
	{3, true},
}
