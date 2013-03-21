package main

import (
	"fmt"
	//"github.com/npadmana/go-xi/particle"
	"github.com/npadmana/go-xi/mesh"
	"github.com/npadmana/go-xi/twopt"
	"github.com/npadmana/go-xi/utils"
	"log"
)

func main() {
	p, err := mesh.ReadParticles("test_N.dat")
	fmt.Println("Read in Particles")
	if err != nil {
		log.Fatal(err)
	}
	m := mesh.New(p, 50.0)
	fmt.Println("Mesh created")

	h := utils.NewUniform([]int{5, 5}, []float64{0, 0}, []float64{200, 1})
	h.Reset()
	for _, g1 := range m.Grid {
		if g1 != nil {
			for _, g2 := range m.Grid {
				if g2 != nil {
					twopt.PairCounter(h, g1.P, g2.P, twopt.SMu)
				}
			}
		}
	}

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			fmt.Print(h.Get(i, j), " ")
		}
		fmt.Println()
	}

}
