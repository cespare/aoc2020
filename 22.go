package main

import (
	"fmt"
	"strings"
)

func init() {
	addSolutions(22, problem22)
}

func problem22(ctx *problemContext) {
	var cards [][]uint8
	scanner := ctx.scanner()
	scanner.s.Split(splitGroups)
	for scanner.scan() {
		header := fmt.Sprintf("Player %d:", len(cards)+1)
		s, ok := trimPrefix(scanner.text(), header)
		if !ok {
			panic("bad input")
		}
		var deck []uint8
		for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
			deck = append(deck, uint8(parseUint(line, 10, 8)))
		}
		cards = append(cards, deck)
	}
	if len(cards) != 2 {
		panic("bad input")
	}
	ctx.reportLoad()

	g := &cardGame1{copyCards(cards[0]), copyCards(cards[1])}
	for !g.step1() {
	}
	ctx.reportPart1(scoreGame(g.player1) + scoreGame(g.player2))

	g2 := newCardGame2(copyCards(cards[0]), copyCards(cards[1]))
	if g2.run() == player1Wins {
		ctx.reportPart2(scoreGame(g2.player1))
	} else {
		ctx.reportPart2(scoreGame(g2.player2))
	}
}

type cardGame1 struct {
	player1 []uint8
	player2 []uint8
}

func (g *cardGame1) step1() (done bool) {
	c1, c2 := g.player1[0], g.player2[0]
	g.player1 = g.player1[1:]
	g.player2 = g.player2[1:]
	if c1 > c2 {
		g.player1 = append(g.player1, c1, c2)
	} else {
		g.player2 = append(g.player2, c2, c1)
	}
	return len(g.player1) == 0 || len(g.player2) == 0
}

func scoreGame(cards []uint8) int {
	var n int
	for i := range cards {
		n += int(cards[len(cards)-i-1]) * (i + 1)
	}
	return n
}

func copyCards(cards []uint8) []uint8 {
	return append([]uint8(nil), cards...)
}

type cardGame2 struct {
	seen    map[string]struct{}
	player1 []uint8
	player2 []uint8
}

func newCardGame2(player1, player2 []uint8) *cardGame2 {
	return &cardGame2{
		seen:    make(map[string]struct{}),
		player1: player1,
		player2: player2,
	}
}

func (g *cardGame2) state() string {
	var b strings.Builder
	for _, card := range g.player1 {
		b.WriteByte(byte(card))
		b.WriteByte(0xfe)
	}
	b.WriteByte(0xff)
	for _, card := range g.player1 {
		b.WriteByte(byte(card))
		b.WriteByte(0xfe)
	}
	return b.String()
}

func (g *cardGame2) run() game2Outcome {
	for {
		if outcome := g.step(); outcome != noWinner {
			return outcome
		}
	}
}

func (g *cardGame2) step() game2Outcome {
	s := g.state()
	if _, ok := g.seen[s]; ok {
		return player1Wins
	}
	g.seen[s] = struct{}{}
	c1, c2 := g.player1[0], g.player2[0]
	g.player1 = g.player1[1:]
	g.player2 = g.player2[1:]
	if len(g.player1) < int(c1) || len(g.player2) < int(c2) {
		if c1 > c2 {
			g.player1 = append(g.player1, c1, c2)
		} else {
			g.player2 = append(g.player2, c2, c1)
		}
	} else {
		sub := newCardGame2(copyCards(g.player1[:c1]), copyCards(g.player2[:c2]))
		if sub.run() == player1Wins {
			g.player1 = append(g.player1, c1, c2)
		} else {
			g.player2 = append(g.player2, c2, c1)
		}
	}
	switch {
	case len(g.player1) == 0:
		return player2Wins
	case len(g.player2) == 0:
		return player1Wins
	default:
		return noWinner
	}
}

type game2Outcome int

const (
	noWinner game2Outcome = iota
	player1Wins
	player2Wins
)
