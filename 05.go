package main

import (
	"log"
	"strconv"
	"strings"
)

func init() {
	addSolutions(5, problem5)
}

func problem5(ctx *problemContext) {
	var maxID int
	ids := make(map[int]struct{})
	replacer := strings.NewReplacer("F", "0", "B", "1", "L", "0", "R", "1")
	scanner := ctx.scanner()
	for scanner.scan() {
		s := replacer.Replace(scanner.text())
		n, err := strconv.ParseInt(s, 2, 32)
		if err != nil {
			log.Fatal(err)
		}
		id := int(n)
		ids[id] = struct{}{}
		if id > maxID {
			maxID = id
		}
	}
	ctx.reportLoad()
	ctx.reportPart1(maxID)

	for n := maxID - 1; ; n-- {
		if _, ok := ids[n]; ok {
			continue
		}
		ctx.reportPart2(n)
		break
	}

}
