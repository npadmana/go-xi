package twopt

import (
	"github.com/npadmana/go-xi/particle"
	"github.com/npadmana/npgo/math/histogram"
)

// PairCounter implements the inner pair counting loop.
func PairCounter(h histogram.Histogrammer, p1, p2 []particle.Particle, f DistFunc) {

	for _, ip1 := range p1 {
		for _, ip2 := range p2 {
			h.Add(ip1.W*ip2.W, f(ip1.X, ip2.X)...)
		}
	}

}
