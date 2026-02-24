package main

import (
	"gobook/2/popcount"
	"testing"
)

func BenchmarkPopcount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		popcount.PopCount(1234567890)
	}
}

func BenchmarkPopcountLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		popcount.PopCountLoop(1234567890)
	}
}

func BenchmarkPopcountBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		popcount.PopCountBytes(1234567890)
	}
}

func BenchmarkPopcountClear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		popcount.PopCountClear(1234567890)
	}
}
