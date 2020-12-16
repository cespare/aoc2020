package main

// solveSAT solves a boolean satisfiability formula presented in CNF using a
// simple DPLL approach. Atoms are arbitrary positive ints. Negated atoms are
// negative.
func solveSAT(cnf [][]int) ([]int, bool) {
	var atoms []int
	atomIdx := make(map[int]int)
	for _, clause := range cnf {
		for _, atom := range clause {
			if atom == 0 {
				panic("zero atom")
			}
			if atom < 0 {
				atom = -atom
			}
			if _, ok := atomIdx[atom]; !ok {
				atomIdx[atom] = len(atoms)
				atoms = append(atoms, atom)
			}
		}
	}
	state := &solverState{
		assn:         make([]solverAssn, len(atoms)),
		clauses:      make([][]uint32, len(cnf)),
		countScratch: make([]int, len(atoms)),
	}
	for i, disj := range cnf {
		clause := make([]uint32, len(disj))
		for j, atom := range disj {
			if atom > 0 {
				clause[j] = uint32(atomIdx[atom]) * 2
			} else {
				clause[j] = uint32(atomIdx[-atom])*2 + 1
			}
		}
		state.clauses[i] = clause
	}
	soln := state.solve()
	if soln == nil {
		return nil, false // unsatisfiable
	}
	var result []int
	for i, a := range soln {
		if a == solverTrue {
			result = append(result, atoms[i])
		}
	}
	return result, true
}

type solverAssn uint8

const (
	solverUnassigned solverAssn = iota
	solverFalse
	solverTrue
)

type solverState struct {
	assn         []solverAssn
	clauses      [][]uint32
	countScratch []int
}

func (s *solverState) solve() []solverAssn {
	// pretty.Println(s)
	// Check if we're done.
	// fmt.Println("Checking done")
	for _, clause := range s.clauses {
		if len(clause) == 0 {
			return nil
		}
	}
	done := true
	for _, assn := range s.assn {
		if assn == solverUnassigned {
			done = false
			break
		}
	}
	if done {
		return s.assn
	}

	// Unit propagation: iterate to fixpoint.
unitPropLoop:
	// fmt.Println("unit prop")
	for _, clause := range s.clauses {
		if len(clause) == 1 {
			if s.propagateUnit(clause[0]) {
				goto unitPropLoop
			}
		}
	}

	// Pure literal elimination: iterate to fixpoint again.
	// fmt.Println("pure lit elim")
pureLitElimLoop:
	for a, assn := range s.assn {
		if assn != solverUnassigned {
			continue
		}
		candidate := uint32(a) << 1
		found := false
		for _, clause := range s.clauses {
			for _, atom := range clause {
				if found {
					if atom^1 == candidate {
						// Not pure.
						continue pureLitElimLoop
					}
				} else {
					if atom == candidate {
						found = true
					} else if atom^1 == candidate {
						candidate ^= 1
						found = true
					}
				}
			}
		}
		if found {
			s.propagateUnit(candidate)
			if candidate&1 == 0 {
				s.assn[candidate>>1] = solverTrue
			} else {
				s.assn[candidate>>1] = solverFalse
			}
			goto pureLitElimLoop
		}
	}

	atom, ok := s.chooseNextAtom()
	if !ok {
		return s.assn // done
	}

	oldClauses := make([][]uint32, len(s.clauses))
	for i, clause := range s.clauses {
		oldClauses[i] = append([]uint32(nil), clause...)
	}

	// fmt.Println("left branch")
	s.assn[atom>>1] = solverTrue
	s.propagateUnit(atom)
	if soln := s.solve(); soln != nil {
		return soln
	}

	// fmt.Println("right branch")
	s.clauses = oldClauses
	s.assn[atom>>1] = solverFalse
	s.propagateUnit(atom ^ 1)
	if soln := s.solve(); soln != nil {
		return soln
	}
	s.assn[atom>>1] = solverUnassigned
	return nil
}

func (s *solverState) propagateUnit(unit uint32) bool {
	changed := false
	assigned := s.assn[unit>>1] != solverUnassigned
	ci := 0
conjLoop:
	for _, clause := range s.clauses {
		di := 0
	disjLoop:
		for _, atom := range clause {
			if atom == unit && (assigned || len(clause) > 1) {
				changed = true
				continue conjLoop // delete this clause
			}
			if atom^1 == unit {
				changed = true
				continue disjLoop // delete this atom
			}
			clause[di] = atom
			di++
		}
		s.clauses[ci] = clause[:di]
		ci++
	}
	s.clauses = s.clauses[:ci]
	return changed
}

func (s *solverState) eliminatePureLiteral(pure uint32) {
	ci := 0
conjLoop:
	for _, clause := range s.clauses {
		for _, atom := range clause {
			if atom == pure {
				continue conjLoop // delete this clause
			}
		}
		s.clauses[ci] = clause
		ci++
	}
	s.clauses = s.clauses[:ci]
}

func (s *solverState) chooseNextAtom() (uint32, bool) {
	for i := range s.countScratch {
		s.countScratch[i] = 0
	}
	for _, clause := range s.clauses {
		for _, atom := range clause {
			s.countScratch[atom>>1]++
		}
	}
	var atom int
	var max int
	for a, count := range s.countScratch {
		if count > max {
			atom = a
			max = count
		}
	}
	if max == 0 {
		return 0, false
	}
	return uint32(atom) << 1, true
}
