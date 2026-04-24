package main

import (
	"bytes"
	"fmt"
	"math/bits"
)

type IntSet struct {
	words []uint64
}

func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *IntSet) Len() int {
	result := 0
	for _, word := range s.words {
		result += bits.OnesCount64(word)
	}
	return result
}

func (s *IntSet) Remove(x int) {
	word, bit := x/64, uint(x%64)
	if word >= len(s.words) {
		return
	}
	s.words[word] &^= 1 << bit
}

func (s *IntSet) Clear() {
	s.words = nil
}

func (s *IntSet) Copy() *IntSet {
	newSet := IntSet{}
	newSet.words = make([]uint64, len(s.words))
	copy(newSet.words, s.words)
	return &newSet
}

func (s *IntSet) AddAll(nums ...int) {
	for _, n := range nums {
		s.Add(n)
	}
}

func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) IntersectWith(t *IntSet) {
	for i, sword := range s.words {
		if i < len(t.words) {
			s.words[i] = sword & t.words[i]
		} else {
			s.words[i] = 0
		}
	}
	for i := len(s.words) - 1; i >= 0 && s.words[i] == 0; i-- {
		s.words = s.words[:i]
	}
}

func (s *IntSet) DifferenceWith(t *IntSet) {
	upper := len(s.words)
	if len(t.words) < upper {
		upper = len(t.words)
	}
	for i := 0; i < upper; i++ {
		s.words[i] &^= t.words[i]
	}
	for i := len(s.words) - 1; i >= 0 && s.words[i] == 0; i-- {
		s.words = s.words[:i]
	}
}

func (s *IntSet) SymmetricDifference(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] ^= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
	for i := len(s.words) - 1; i >= 0 && s.words[i] == 0; i-- {
		s.words = s.words[:i]
	}
}

func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func main() {
	s := &IntSet{}
	t := &IntSet{}
	s.AddAll(1, 2, 3, 4)
	t.AddAll(3, 4, 5, 6)
	s.IntersectWith(t)
	fmt.Println("IntersectWith:", s) // {3 4}

	s = &IntSet{}
	t = &IntSet{}
	s.AddAll(1, 2, 3, 4)
	t.AddAll(3, 4, 5, 6)
	fmt.Println("DifferenceWith:", s) // {1 2}

	s = &IntSet{}
	t = &IntSet{}
	s.AddAll(1, 2, 3)
	t.AddAll(3, 4, 5)
	s.SymmetricDifference(t)
	fmt.Println("SymmetricDiffence 1:", s) // {1 2 4 5}

	s = &IntSet{}
	t = &IntSet{}
	s.AddAll(1, 4, 9, 144)
	t.AddAll(9, 42, 144, 200)
	s.SymmetricDifference(t)
	fmt.Println("SymmetricDifference 2:", s) // {1 4 42 200}
}
