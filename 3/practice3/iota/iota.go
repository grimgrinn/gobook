package main

import "fmt"

const (
	KiB = 1 << (10 * (iota + 1))
	MiB
	GiB
	TiB
	PiB
	EiB
	ZiB
	YiB
)

func main() {
	fmt.Println(
		KiB,
		MiB,
		GiB,
		TiB,
		PiB,
		EiB,
		//ZiB,    //overflow
		//YiB,
	)
}
