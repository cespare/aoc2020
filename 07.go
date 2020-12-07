package main

import (
	"regexp"
	"strings"

	"github.com/cespare/hasty"
)

func init() {
	addSolutions(7, problem7)
}

func problem7(ctx *problemContext) {
	rules := make(map[string]bagRule)
	scanner := ctx.scanner()
	for scanner.scan() {
		br := parseBagRule(scanner.text())
		rules[br.typ] = br
	}
	ctx.reportLoad()

	var gold int
	for _, rule := range rules {
		if hasShinyGold(rules, rule.typ, make(map[string]struct{})) {
			gold++
		}
	}
	ctx.reportPart1(gold)

	ctx.reportPart2(countContainedBags(rules, "shiny gold"))
}

type bagRule struct {
	typ      string
	contains []numBags
}

type numBags struct {
	N    int64
	Type string
}

var numBagsRegexp = regexp.MustCompile(`^(?P<N>\d+) (?P<Type>\w+ \w+) bags?$`)

func parseNumBags(s string) numBags {
	var nb numBags
	hasty.MustParse([]byte(s), &nb, numBagsRegexp)
	return nb
}

func parseBagRule(s string) bagRule {
	s = strings.TrimSuffix(s, ".")
	parts := strings.SplitN(s, " contain ", 2)
	br := bagRule{
		typ: strings.TrimSuffix(parts[0], " bags"),
	}
	if parts[1] == "no other bags" {
		return br
	}
	for _, c := range strings.Split(parts[1], ",") {
		nb := parseNumBags(strings.TrimSpace(c))
		br.contains = append(br.contains, nb)
	}
	return br
}

func hasShinyGold(rules map[string]bagRule, typ string, checked map[string]struct{}) bool {
	rule := rules[typ]
	for _, nb := range rule.contains {
		if nb.Type == "shiny gold" {
			return true
		}
		if _, ok := checked[nb.Type]; ok {
			continue
		}
		checked[nb.Type] = struct{}{}
		if hasShinyGold(rules, nb.Type, checked) {
			return true
		}
	}
	return false
}

func countContainedBags(rules map[string]bagRule, typ string) int64 {
	rule := rules[typ]
	var total int64
	for _, nb := range rule.contains {
		total += nb.N + nb.N*countContainedBags(rules, nb.Type)
	}
	return total
}
