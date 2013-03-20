package main 

import (
	//"fmt"
	"log"
	"github.com/npadmana/go-xi/particle"
	"github.com/npadmana/npgo/math/histogram"
	"github.com/npadmana/go-xi/twopt"
)

func main() {
	parr, err := particle.NewFromXYZW("boss_S.shift.small")
	if err != nil {
		log.Fatal(err)
	}
	
	h := histogram.NewUniform([]int{101, 100}, []float64{0, 0}, []float64{202, 1})
	h.Reset()
	twopt.PairCounter(h, parr.Data, parr.Data, twopt.SMu)

}
