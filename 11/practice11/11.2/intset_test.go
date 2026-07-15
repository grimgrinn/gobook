package main

import (
	"math/rand"
	"testing"
)

type mapIntSet struct {
	m map[int]bool
}

func newMapIntSet() *mapIntSet {
	return &mapIntSet{m: make(map[int]bool)}
}

func (s *mapIntSet) Add(x int) {
	s.m[x] = true
}

func (s *mapIntSet) Has(x int) bool {
	return s.m[x]
}

func (s *mapIntSet) Remove(x int) {
	delete(s.m, x)
}

func (s *mapIntSet) Clear() {
	s.m = make(map[int]bool)
}

func (s *mapIntSet) Len() int {
	return len(s.m)
}

func (s *mapIntSet) Elems() []int {
	elems := make([]int, 0, len(s.m))
	for x := range s.m {
		elems = append(elems, x)
	}
	return elems
}

func (s *mapIntSet) UnionWith(t *mapIntSet) {
	for x := range t.m {
		s.m[x] = true
	}
}

func (s *mapIntSet) IntersectWith(t *mapIntSet) {
	for x := range s.m {
		if !t.m[x] {
			delete(s.m, x)
		}
	}
}

func (s *mapIntSet) DifferenceWith(t *mapIntSet) {
	for x := range t.m {
		delete(s.m, x)
	}
}

func (s *mapIntSet) SymmetricDifference(t *mapIntSet) {
	for x := range t.m {
		if s.m[x] {
			delete(s.m, x)
		} else {
			s.m[x] = true
		}
	}
}

func (s *mapIntSet) Copy() *mapIntSet {
	ns := newMapIntSet()
	for x := range s.m {
		ns.m[x] = true
	}
	return ns
}

func equal(a, b *IntSet) bool {
	elemsA := a.Elems()
	elemsB := b.Elems()

	if len(elemsA) != len(elemsB) {
		return false
	}

	seen := make(map[int]bool)
	for _, x := range elemsA {
		seen[x] = true
	}
	for _, x := range elemsB {
		if !seen[x] {
			return false
		}
	}
	return true
}

func equalIntSetMap(s *IntSet, ms *mapIntSet) bool {
	elemsA := s.Elems()
	elemsB := ms.Elems()

	if len(elemsA) != len(elemsB) {
		return false
	}

	seen := make(map[int]bool)
	for _, x := range elemsA {
		seen[x] = true
	}
	for _, x := range elemsB {
		if !seen[x] {
			return false
		}
	}
	return true
}

func TestIntSet(t *testing.T) {
	for i := 0; i < 1000; i++ {
		s := &IntSet{}
		ms := newMapIntSet()

		for j := 0; j < 10; j++ {
			x := rand.Intn(100)
			s.Add(x)
			ms.Add(x)
			if !equalIntSetMap(s, ms) {
				t.Errorf("после Add(%d): s=%v. ms=%v", x, s.Elems(), ms.Elems())
			}
		}

		for j := 0; j < 10; j++ {
			x := rand.Intn(100)
			if s.Has(x) != ms.Has(x) {
				t.Errorf("Has(%d): s=%v, ms=%v", x, s.Has(x), ms.Has(x))
			}
		}

		if s.Len() != ms.Len() {
			t.Errorf("Len: s=%d, ms=%d", s.Len(), ms.Len())
		}

		for j := 0; j < 5; j++ {
			x := rand.Intn(100)
			s.Remove(x)
			ms.Remove(x)
			if !equalIntSetMap(s, ms) {
				t.Errorf("после Remove(%d): s=%v, ms=%v", x, s.Elems(), ms.Elems())
			}
		}

		s2 := &IntSet{}
		ms2 := newMapIntSet()
		for j := 0; j < 5; j++ {
			x := rand.Intn(100)
			s2.Add(x)
			ms2.Add(x)
		}
		s.UnionWith(s2)
		ms.UnionWith(ms2)
		if !equalIntSetMap(s, ms) {
			t.Errorf("после UnionWIth: s=%v, ms=%v", s.Elems(), ms.Elems())
		}

		s3 := &IntSet{}
		ms3 := newMapIntSet()
		for j := 0; j < 5; j++ {
			x := rand.Intn(100)
			s3.Add(x)
			ms3.Add(x)
		}
		s.IntersectWith(s3)
		ms.IntersectWith(ms3)
		if !equalIntSetMap(s, ms) {
			t.Errorf("после InterSectWith: s=%v, ms=%v", s.Elems(), ms.Elems())
		}

		s4 := &IntSet{}
		ms4 := newMapIntSet()
		for j := 0; j < 5; j++ {
			x := rand.Intn(100)
			s4.Add(x)
			ms4.Add(x)
		}
		s.DifferenceWith(s4)
		ms.DifferenceWith(ms4)
		if !equalIntSetMap(s, ms) {
			t.Errorf("после DifferenceWith: s=%v, ms=%v", s.Elems(), ms.Elems())
		}

		s5 := &IntSet{}
		ms5 := newMapIntSet()
		for j := 0; j < 5; j++ {
			x := rand.Intn(100)
			s5.Add(x)
			ms5.Add(x)
		}
		s.SymmetricDifference(s5)
		ms.SymmetricDifference(ms5)
		if !equalIntSetMap(s, ms) {
			t.Errorf("после SymmmetricDifference: s=%v, ms=%v", s.Elems(), ms.Elems())
		}

		s.Clear()
		ms.Clear()
		if !equalIntSetMap(s, ms) {
			t.Errorf("после Clear: s=%v, ms=%v", s.Elems(), ms.Elems())
		}

		s.AddAll(1, 2, 3, 4, 5)
		ms.Add(1)
		ms.Add(2)
		ms.Add(3)
		ms.Add(4)
		ms.Add(5)
		copyS := s.Copy()
		copyMs := ms.Copy()
		if !equalIntSetMap(copyS, copyMs) {
			t.Errorf("после Copy: s=%v, ms=%v", copyS.Elems(), copyMs.Elems())
		}
	}
}
