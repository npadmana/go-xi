package twopt

import (
	"github.com/npadmana/go-xi/mesh"
	"github.com/npadmana/go-xi/utils"
)

// PairCounter implements the inner pair counting loop.
func PairCounter(h utils.Histogrammer, p1, p2 []mesh.Particle, f DistFunc) {

	for _, ip1 := range p1 {
		x1 := ip1.X
		w1 := ip1.W
		for _, ip2 := range p2 {
			h.Add(w1*ip2.W, f(x1, ip2.X)...)
		}
	}

}
