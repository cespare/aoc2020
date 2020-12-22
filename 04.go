package main

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	addSolutions(4, problem4)
}

func problem4(ctx *problemContext) {
	var passports []passport
	scanner := ctx.scanner()
	scanner.s.Split(splitGroups)
	for scanner.scan() {
		passports = append(passports, parsePassport(scanner.text()))
	}
	ctx.reportLoad()

	var valid1, valid2 int
	for _, pp := range passports {
		if len(pp) != 7 {
			continue
		}
		valid1++
		if pp.valid() {
			valid2++
		}
	}
	ctx.reportPart1(valid1)
	ctx.reportPart2(valid2)
}

func splitGroups(data []byte, atEOF bool) (advance int, token []byte, err error) {
	i := bytes.Index(data, []byte("\n\n"))
	if i < 0 {
		if atEOF && len(data) > 0 {
			if data[len(data)-1] == '\n' {
				return len(data), data[:len(data)-1], nil
			}
			return len(data), data, nil
		}
		return 0, nil, nil
	}
	return i + 2, data[:i], nil
}

type passport map[string]string

func parsePassport(s string) map[string]string {
	m := make(map[string]string)
	for _, field := range strings.Fields(s) {
		parts := strings.SplitN(field, ":", 2)
		if len(parts) != 2 {
			continue
		}
		switch parts[0] {
		case "byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid":
			m[parts[0]] = parts[1]
		}
	}
	return m
}

var (
	hclRegexp = regexp.MustCompile(`^#[0-9a-f]{6}$`)
	pidRegexp = regexp.MustCompile(`^[0-9]{9}$`)
)

func (pp passport) valid() bool {
	inRange := func(s string, min, max int) bool {
		n, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		return n >= min && n <= max
	}
	for _, check := range []struct {
		field string
		f     func(s string) bool
	}{
		{"byr", func(s string) bool { return inRange(s, 1920, 2002) }},
		{"iyr", func(s string) bool { return inRange(s, 2010, 2020) }},
		{"eyr", func(s string) bool { return inRange(s, 2020, 2030) }},
		{
			"hgt",
			func(s string) bool {
				if h, ok := trimSuffix(pp["hgt"], "cm"); ok {
					return inRange(h, 150, 193)
				}
				if h, ok := trimSuffix(pp["hgt"], "in"); ok {
					return inRange(h, 59, 76)
				}
				return false
			},
		},
		{"hcl", hclRegexp.MatchString},
		{
			"ecl",
			func(s string) bool {
				switch s {
				case "amb", "blu", "brn", "gry", "grn", "hzl", "oth":
					return true
				default:
					return false
				}
			},
		},
		{"pid", pidRegexp.MatchString},
	} {
		if !check.f(pp[check.field]) {
			return false
		}
	}
	return true
}
