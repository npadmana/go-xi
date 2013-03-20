package main

import (
	"fmt"
	"github.com/npadmana/go-xi/particle"
	"github.com/npadmana/go-xi/twopt"
	"github.com/npadmana/npgo/math/histogram"
	"log"
)

func main() {
	parr, err := particle.NewFromXYZW("test_N.dat")
	if err != nil {
		log.Fatal(err)
	}

	h := histogram.NewUniform([]int{5, 5}, []float64{0, 0}, []float64{200, 1})
	h.Reset()
	twopt.PairCounter(h, parr.Data, parr.Data, twopt.SMu)

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			fmt.Print(h.Get(i, j)," ")
		}
		fmt.Println()
	}

}
