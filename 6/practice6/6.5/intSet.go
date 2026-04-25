package main

import (
	"bytes"
	"fmt"
	"math/bits"
)

const wordSize = 32 << (^uint(0) >> 63)

type IntSet struct {
	words []uint
}

func (s *IntSet) Has(x int) bool {
	word, bit := x/wordSize, uint(x%wordSize)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

func (s *IntSet) Add(x int) {
	word, bit := x/wordSize, uint(x%wordSize)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *IntSet) Len() int {
	result := 0
	for _, word := range s.words {
		result += bits.OnesCount(uint(word))
	}
	return result
}

func (s *IntSet) Remove(x int) {
	word, bit := x/wordSize, uint(x%wordSize)
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
	newSet.words = make([]uint, len(s.words))
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

func (s *IntSet) Elems() []int {
	result := make([]int, 0, s.Len())

	for i, word := range s.words {
		for bit := 0; bit < wordSize; bit++ {
			if word&(1<<uint(bit)) != 0 {
				result = append(result, i*wordSize+bit)
			}
		}
	}
	return result
}

func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < wordSize; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", wordSize*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func main() {
	s := &IntSet{}
	s.Add(1)
	s.Add(243)
	s.Add(8)
	fmt.Println(s.String()) // "{1 8 243}"
	fmt.Println(s.Len())    // 3

}
