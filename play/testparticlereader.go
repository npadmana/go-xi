package main

import (
	"fmt"
	"github.com/npadmana/go-xi/particle"
	"github.com/npadmana/npgo/math/histogram"
	"github.com/npadmana/go-xi/twopt"
	"log"
)

func main() {
	parr, err := particle.NewFromXYZW("testparticlereader.dat")
	if err != nil {
		log.Fatal(err)
	}
	
	parr.SetBox()
	
	for _,p1 := range parr.Data {
		fmt.Println(p1)
	}
	
	fmt.Println("BoxDim=", parr.BoxDim)
	fmt.Println("BoxMin=", parr.BoxMin)

	h := histogram.NewUniform([]int{10, 2}, []float64{0, 0}, []float64{100, 1})
	h.Reset()
	twopt.PairCounter(h, parr.Data, parr.Data, twopt.SMu)

}
