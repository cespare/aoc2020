package main

import (
	"regexp"
	"strings"
)

func init() {
	addSolutions(19, problem19)
}

func problem19(ctx *problemContext) {
	var rules []*messageRule
	scanner := ctx.scanner()
	for scanner.scan() {
		if scanner.text() == "" {
			break
		}
		rules = append(rules, parseMessageRule(scanner.text()))
	}
	var inputs []string
	for scanner.scan() {
		inputs = append(inputs, scanner.text())
	}
	ctx.reportLoad()

	m := make(map[int]*messageRule)
	for _, rule := range rules {
		if _, ok := m[rule.id]; ok {
			panic("duplicate rule IDs")
		}
		m[rule.id] = rule
	}
	for _, rule := range rules {
		if rule.lit != "" {
			continue
		}
		rule.opts = make([][]*messageRule, len(rule.optsIdx))
		for i, idxs := range rule.optsIdx {
			rule.opts[i] = make([]*messageRule, len(idxs))
			for j, idx := range idxs {
				r, ok := m[idx]
				if !ok {
					panic("nonexistent index")
				}
				rule.opts[i][j] = r
			}
		}
	}

	re := regexp.MustCompile("^" + m[0].regexpText() + "$")
	var n int64
	for _, in := range inputs {
		if re.MatchString(in) {
			n++
		}
	}
	ctx.reportPart1(n)

	rule42 := regexp.MustCompile("^" + m[42].regexpText())
	rule31 := regexp.MustCompile(m[31].regexpText() + "$")
	n = 0
	for _, in := range inputs {
		if messageMatchPart2(rule42, rule31, in) {
			n++
		}
	}
	ctx.reportPart2(n)
}

func messageMatchPart2(rule42, rule31 *regexp.Regexp, s string) bool {
	m := 0
	for ; ; m++ {
		suf := rule31.FindString(s)
		if suf == "" {
			break
		}
		s = strings.TrimSuffix(s, suf)
	}
	n := 0
	for ; ; n++ {
		pre := rule42.FindString(s)
		if pre == "" {
			break
		}
		s = strings.TrimPrefix(s, pre)
	}
	return s == "" && m > 0 && n > m
}

type messageRule struct {
	id      int
	optsIdx [][]int
	opts    [][]*messageRule
	lit     string // if opts is nil
}

func (r *messageRule) regexpText() string {
	if r.opts == nil {
		return r.lit
	}
	var b strings.Builder
	if len(r.opts) > 1 {
		b.WriteByte('(')
	}
	for i, opt := range r.opts {
		if i > 0 {
			b.WriteByte('|')
		}
		for _, r := range opt {
			b.WriteString(r.regexpText())
		}
	}
	if len(r.opts) > 1 {
		b.WriteByte(')')
	}
	return b.String()
}

func parseMessageRule(s string) *messageRule {
	r := new(messageRule)
	parts := strings.SplitN(s, ": ", 2)
	r.id = int(parseInt(parts[0], 10, 64))
	if strings.Contains(parts[1], `"`) {
		r.lit = strings.TrimFunc(parts[1], func(c rune) bool { return c == '"' })
		return r
	}
	for _, opt := range strings.Split(parts[1], "|") {
		var subs []int
		for _, f := range strings.Fields(opt) {
			subs = append(subs, int(parseInt(f, 10, 64)))
		}
		r.optsIdx = append(r.optsIdx, subs)
	}
	return r
}
