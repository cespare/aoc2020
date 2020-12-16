package main

import (
	"regexp"
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
	table := make([][]bool, k) // k x k truth table describing the constraints
	for i := range table {
		row := make([]bool, k)
		for j := range row {
			row[j] = true
		}
		table[i] = row
	}
	for _, tick := range tickets {
		for i, n := range tick {
			for j, rule := range rules {
				table[i][j] = table[i][j] && rule.matches(n)
			}
		}
	}

	// atom # (1-indexed) is fieldIdx * k + ruleIdx + 1

	var cnfClauses [][]int
	// For each row, use a clause to indicate that one of the available
	// rules must match and multiple negative atoms to mark the positions
	// that are disallowed.
	for i, row := range table {
		var allowed []int
		for j, ok := range row {
			atom := i*k + j + 1
			if ok {
				allowed = append(allowed, atom)
			} else {
				cnfClauses = append(cnfClauses, []int{-atom})
			}
		}
		cnfClauses = append(cnfClauses, allowed)
		// Now add the constraint that each field index must correspond
		// to only one rule. For each pair of valid fields A, B:
		// ¬(A ∧ B) -> ¬A ∨ ¬B.
		for d, a0 := range allowed {
			for _, a1 := range allowed[d+1:] {
				pair := []int{-a0, -a1}
				cnfClauses = append(cnfClauses, pair)
			}
		}
	}
	// For each column, indicate that each rule must correspond to only one
	// field by using negative pairs the same way.
	for j := 0; j < k; j++ {
		var allowed []int
		for i := 0; i < k; i++ {
			if table[i][j] {
				atom := i*k + j + 1
				allowed = append(allowed, atom)
			}
		}
		for d, a0 := range allowed {
			for _, a1 := range allowed[d+1:] {
				pair := []int{-a0, -a1}
				cnfClauses = append(cnfClauses, pair)
			}
		}
	}
	var t int
	for _, clause := range cnfClauses {
		t += len(clause)
	}

	soln, ok := solveSAT(cnfClauses)
	if !ok {
		panic("no solution")
	}

	m := int64(1)
	for _, atom := range soln {
		atom--
		i, j := atom/k, atom%k
		rule := rules[j]
		if !strings.HasPrefix(rule.name, "departure ") {
			continue
		}
		m *= int64(yourTicket[i])
	}
	ctx.reportPart2(m)
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
