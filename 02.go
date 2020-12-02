package main

import (
	"regexp"
	"strings"

	"github.com/cespare/hasty"
)

func init() {
	addSolutions(2, problem2)
}

func problem2(ctx *problemContext) {
	var policies []pwPolicy
	scanner := ctx.scanner()
	for scanner.scan() {
		var pol pwPolicy
		hasty.MustParse(scanner.s.Bytes(), &pol, pwPolicyRegexp)
		policies = append(policies, pol)
	}
	ctx.reportLoad()

	var numValid int
	for _, pol := range policies {
		if pol.valid1() {
			numValid++
		}
	}
	ctx.reportPart1(numValid)

	numValid = 0
	for _, pol := range policies {
		if pol.valid2() {
			numValid++
		}
	}
	ctx.reportPart2(numValid)
}

var pwPolicyRegexp = regexp.MustCompile(`^(?P<Min>\d+)-(?P<Max>\d+) (?P<Letter>.): (?P<Password>.*)$`)

type pwPolicy struct {
	Min      int
	Max      int
	Letter   string
	Password string
}

func (p pwPolicy) valid1() bool {
	n := strings.Count(p.Password, p.Letter)
	return n >= p.Min && n <= p.Max
}

func (p pwPolicy) valid2() bool {
	c := p.Letter[0]
	return (p.Password[p.Min-1] == c) != (p.Password[p.Max-1] == c)
}
