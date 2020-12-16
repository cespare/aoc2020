package main

import (
	"regexp"
	"sort"
	"strings"

	"github.com/cespare/hasty"
)

func init() {
	addSolutions(16, problem16)
}

func problem16(ctx *problemContext) {
	var rules []ticketRule
	var yourTicket []int64
	var nearbyTickets [][]int64
	scanner := ctx.scanner()
	var state int
	for scanner.scan() {
		switch state {
		case 0:
			if scanner.text() == "" {
				state++
				continue
			}
			rules = append(rules, parseTicketRule(scanner.text()))
		case 1:
			if scanner.text() != "your ticket:" {
				panic("malformed")
			}
			state++
		case 2:
			yourTicket = commaSeparatedInt64s(scanner.text())
			state++
		case 3:
			if scanner.text() != "" {
				panic("malformed")
			}
			state++
		case 4:
			if scanner.text() != "nearby tickets:" {
				panic("malformed")
			}
			state++
		case 5:
			nearbyTickets = append(nearbyTickets, commaSeparatedInt64s(scanner.text()))
		}
	}
	ctx.reportLoad()

	tickets := [][]int64{yourTicket}
	var errRate int64
	for _, tick := range nearbyTickets {
		valid := true
	ticketLoop:
		for _, n := range tick {
			for _, rule := range rules {
				if rule.matches(n) {
					continue ticketLoop
				}
			}
			valid = false
			errRate += n
		}
		if valid {
			tickets = append(tickets, tick)
		}
	}
	ctx.reportPart1(errRate)

	k := len(tickets[0])
	candidates := make([]map[int]struct{}, k)
	for i := range candidates {
		candidate := make(map[int]struct{})
		for j := range rules {
			candidate[j] = struct{}{}
		}
		candidates[i] = candidate
	}
	for _, tick := range tickets {
		for i, n := range tick {
			for j, rule := range rules {
				if !rule.matches(n) {
					delete(candidates[i], j)
				}
			}
		}
	}

	order := make([]int, k)
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		return len(candidates[i]) < len(candidates[j])
	})

	soln := make([]int, k)
	for i := range soln {
		soln[i] = -1
	}
	var solve func(int) bool
	solve = func(n int) bool {
		if n == k {
			return true
		}
		i := order[n]
	candidateLoop:
		for c := range candidates[i] {
			for _, c1 := range soln {
				if c == c1 { // already used
					continue candidateLoop
				}
			}
			soln[i] = c
			if solve(n + 1) {
				return true
			}
		}
		soln[i] = -1
		return false
	}
	solve(0)

	m := int64(1)
	for i, j := range soln {
		rule := rules[j]
		if !strings.HasPrefix(rule.name, "departure ") {
			continue
		}
		m *= int64(yourTicket[i])
	}
	ctx.reportPart2(m)
}

func indexIn(ns []int, n int) int {
	for i, n1 := range ns {
		if n1 == n {
			return i
		}
	}
	panic("not found")
}

type ticketRule struct {
	name   string
	range0 [2]int64
	range1 [2]int64
}

var ticketRuleRegexp = regexp.MustCompile(`^(?P<Name>.*): (?P<N0>\d+)-(?P<N1>\d+) or (?P<N2>\d+)-(?P<N3>\d+)$`)

func parseTicketRule(s string) ticketRule {
	var v struct {
		Name           string
		N0, N1, N2, N3 int64
	}
	hasty.MustParse([]byte(s), &v, ticketRuleRegexp)
	return ticketRule{
		name:   v.Name,
		range0: [2]int64{v.N0, v.N1},
		range1: [2]int64{v.N2, v.N3},
	}
}

func (r ticketRule) matches(n int64) bool {
	return n >= r.range0[0] && n <= r.range0[1] || n >= r.range1[0] && n <= r.range1[1]
}

func commaSeparatedInt64s(s string) []int64 {
	fields := strings.Split(s, ",")
	ns := make([]int64, len(fields))
	for i, f := range fields {
		ns[i] = parseInt(f, 10, 64)
	}
	return ns
}
