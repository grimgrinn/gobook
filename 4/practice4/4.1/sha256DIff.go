package main

import (
	"crypto/sha256"
	"fmt"
	"gobook/2/popcount"
)

func Sha256Diff(a, b [32]byte) int {
	diff := 0
	for i := 0; i < 32; i++ {
		diff += popcount.PopCount(uint64(a[i] ^ b[i]))
	}
	return diff
}

func main() {
	h1 := sha256.Sum256([]byte("hello"))
	h2 := sha256.Sum256([]byte("world"))
	h3 := sha256.Sum256([]byte("hello"))
	a := [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	b := [32]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	fmt.Printf("Различающихся битов: %d\n", Sha256Diff(h1, h2))
	fmt.Printf("Различающихся битов: %d\n", Sha256Diff(h1, h3))
	fmt.Printf("Различающихся битов: %d\n", Sha256Diff(a, b))
}
