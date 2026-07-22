package popcount

import "testing"

func BenchmarkPopCountTalbe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountTable(0x1234567890ABCDEF)
	}
}

func BenchmarkPopcountShift(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountShift(0x1234567890ABCDEF)
	}
}

func BenchmarkPopCountClear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountClear(0x1234567890ABCDEF)
	}
}
