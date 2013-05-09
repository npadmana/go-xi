package main

import (
	"fmt"
	"github.com/npadmana/go-xi/cuda/cudalib"
)

func main() {
	pp := make([]cudalib.Float4, 100)
	// Allocate particles on GPU
	ppd := cudalib.NewDevParticle(100)
	fmt.Printf("%v \n", ppd)
	ppd.CopyToDevice(pp)

	// Do the same with the histogram
	hh := make([]uint64, 70)
	hhd := cudalib.NewDevHist(70)
	fmt.Printf("%v\n", hhd)
	hhd.CopyFromDevice(hh)

	// Free stuff
	hhd.Free()
	ppd.Free()

}
