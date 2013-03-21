package main

import (
	"fmt"
	//"github.com/npadmana/go-xi/particle"
	"github.com/npadmana/go-xi/twopt"
	"github.com/npadmana/npgo/math/histogram"
	"github.com/npadmana/go-xi/mesh"
	"log"
)

func main() {
	m, err := mesh.New("test_N.dat", 50.0)
	if err != nil {
		log.Fatal(err)
	}

	h := histogram.NewUniform([]int{5, 5}, []float64{0, 0}, []float64{200, 1})
	h.Reset()
	twopt.PairCounter(h, m.Particles, m.Particles, twopt.SMu)

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			fmt.Print(h.Get(i, j)," ")
		}
		fmt.Println()
	}

}
