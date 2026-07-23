package intset

import (
	"math/rand"
	"testing"
)

type MapSet struct {
	m map[int]bool
}

func NewMapSet() *MapSet {
	return &MapSet{m: make(map[int]bool)}
}

func (s *MapSet) Add(x int) {
	s.m[x] = true
}

func (s *MapSet) Has(x int) bool {
	return s.m[x]
}

func (s *MapSet) UnionWith(t *MapSet) {
	for x := range t.m {
		s.m[x] = true
	}
}

func (s *MapSet) IntersectWith(t *MapSet) {
	for x := range s.m {
		if !t.m[x] {
			delete(s.m, x)
		}
	}
}

func (s *MapSet) Dirrerencewith(t *MapSet) {
	for x := range t.m {
		delete(s.m, x)
	}
}

func (s *MapSet) SymmetricDifferenceWith(t *MapSet) {
	for x := range t.m {
		if s.m[x] {
			delete(s.m, x)
		} else {
			s.m[x] = true
		}
	}
}

func (s *MapSet) Len() int {
	return len(s.m)
}

func (s *MapSet) Elems() []int {
	elems := make([]int, 0, len(s.m))
	for x := range s.m {
		elems = append(elems, x)
	}
	return elems
}

func randomInts(n int) []int {
	r := rand.New(rand.NewSource(1))
	ints := make([]int, n)
	for i := range ints {
		ints[i] = r.Intn(1 << 20) // 0..1,048,575
	}
	return ints
}

func benchmarkAdd(b *testing.B, n int) {
	ints := randomInts(n)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := &IntSet{}
		for _, x := range ints {
			s.Add(x)
		}
	}
}

func BenchmarkAdd100(b *testing.B)   { benchmarkAdd(b, 100) }
func BenchmarkAdd1000(b *testing.B)  { benchmarkAdd(b, 1000) }
func BenchmarkAdd10000(b *testing.B) { benchmarkAdd(b, 10000) }

func benchmarkUnionWith(b *testing.B, n int) {
	ints1 := randomInts(n)
	ints2 := randomInts(n)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s1 := &IntSet{}
		s2 := &IntSet{}
		for _, x := range ints1 {
			s1.Add(x)
		}
		for _, x := range ints2 {
			s2.Add(x)
		}
		s1.UnionWith(s2)
	}
}

func BenchmarkUnionWith100(b *testing.B)   { benchmarkUnionWith(b, 100) }
func BenchmarkUnionWIth1000(b *testing.B)  { benchmarkUnionWith(b, 1000) }
func BenchmarkUnionWith10000(b *testing.B) { benchmarkUnionWith(b, 10000) }

func benchmarkMapAdd(b *testing.B, n int) {
	ints := randomInts(n)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := NewMapSet()
		for _, x := range ints {
			s.Add(x)
		}
	}
}

func BenchmarkMapAdd100(b *testing.B)   { benchmarkMapAdd(b, 100) }
func BenchmarkMapAdd1000(b *testing.B)  { benchmarkMapAdd(b, 1000) }
func BenchmarkMapAdd10000(b *testing.B) { benchmarkMapAdd(b, 10000) }

func benchmarkMapUnionWith(b *testing.B, n int) {
	ints1 := randomInts(n)
	ints2 := randomInts(n)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s1 := NewMapSet()
		s2 := NewMapSet()
		for _, x := range ints1 {
			s1.Add(x)
		}
		for _, x := range ints2 {
			s2.Add(x)
		}
		s1.UnionWith(s2)
	}
}

func BenchmarkMapUnionWith100(b *testing.B)   { benchmarkMapUnionWith(b, 100) }
func BenchmarkMapUnionWith1000(b *testing.B)  { benchmarkMapUnionWith(b, 1000) }
func BenchmarkMapUnionWith10000(b *testing.B) { benchmarkMapUnionWith(b, 10000) }
