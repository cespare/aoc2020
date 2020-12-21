package main

import (
	"sort"
	"strings"
)

func init() {
	addSolutions(21, problem21)
}

func problem21(ctx *problemContext) {
	var foods []*foodList
	scanner := ctx.scanner()
	for scanner.scan() {
		foods = append(foods, parseFoodList(scanner.text()))
	}
	ctx.reportLoad()

	allergenIngreds := make(map[string]stringSet) // allergen to set of ingredients
	for _, food := range foods {
		for allergen := range food.allergens {
			if s, ok := allergenIngreds[allergen]; ok {
				s.intersect(food.ingredients)
			} else {
				allergenIngreds[allergen] = food.ingredients.copy()
			}
		}
	}

	badIngreds := make(stringSet)
	for _, s := range allergenIngreds {
		badIngreds.union(s)
	}

	var part1 int
	for _, food := range foods {
		for ingred := range food.ingredients {
			if _, ok := badIngreds[ingred]; !ok {
				part1++
			}
		}
	}
	ctx.reportPart1(part1)

	ingredAllergens := make(map[string]string) // ingredient to known allergen
	for len(allergenIngreds) > 0 {
		found := make(stringSet)
		for allergen, ingreds := range allergenIngreds {
			if len(ingreds) == 1 {
				var ingred string
				for k := range ingreds {
					ingred = k
				}
				ingredAllergens[ingred] = allergen
				found[ingred] = struct{}{}
			}
		}
		for allergen, ingreds := range allergenIngreds {
			for ingred := range ingreds {
				if _, ok := found[ingred]; ok {
					delete(ingreds, ingred)
				}
				if len(ingreds) == 0 {
					delete(allergenIngreds, allergen)
				}
			}
		}
	}

	var ingredAllergenPairs [][2]string
	for ingred, allergen := range ingredAllergens {
		ingredAllergenPairs = append(ingredAllergenPairs, [2]string{ingred, allergen})
	}
	sort.Slice(ingredAllergenPairs, func(i, j int) bool {
		return ingredAllergenPairs[i][1] < ingredAllergenPairs[j][1]
	})
	var b strings.Builder
	for i, pair := range ingredAllergenPairs {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(pair[0])
	}
	ctx.reportPart2(b.String())
}

type foodList struct {
	ingredients stringSet
	allergens   stringSet
}

func parseFoodList(s string) *foodList {
	fl := &foodList{
		ingredients: make(stringSet),
		allergens:   make(stringSet),
	}
	parts := strings.SplitN(s, "(contains ", 2)
	for _, ingred := range strings.Fields(parts[0]) {
		fl.ingredients[ingred] = struct{}{}
	}
	for _, s1 := range strings.Split(strings.TrimSuffix(parts[1], ")"), ",") {
		fl.allergens[strings.TrimSpace(s1)] = struct{}{}
	}
	return fl
}

type stringSet map[string]struct{}

func (s stringSet) copy() stringSet {
	s1 := make(stringSet, len(s))
	for v := range s {
		s1[v] = struct{}{}
	}
	return s1
}

func (s stringSet) union(s1 stringSet) {
	for v := range s1 {
		s[v] = struct{}{}
	}
}

func (s stringSet) intersect(s1 stringSet) {
	for v := range s {
		if _, ok := s1[v]; !ok {
			delete(s, v)
		}
	}
}
