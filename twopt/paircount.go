package twopt

import (
	"github.com/npadmana/go-xi/mesh"
	"github.com/npadmana/go-xi/utils"
	"math"
)

// PairCounter implements the inner pair counting loop.
func PairCounter(h utils.Histogrammer, p1, p2 []mesh.Particle, scale float64) {
	var s2, l2, sl, s1, l1, mu, w1 float64
	//var x1 utils.Vector3D

	n1 := len(p1)
	n2 := len(p2)
	var ip1, ip2, i int

	for ip1=0; ip1 < n1; ip1++ {
		//x1 = p1[ip1].X
		w1 = p1[ip1].W * scale
		for ip2 = 0; ip2 < n2; ip2++ {
			s2, l2, sl = 0, 0, 0
			for i = 0; i < 3; i++ {
				s1 = p1[ip1].X[i] - p2[ip2].X[i]
				l1 = 0.5 * (p1[ip1].X[i] + p2[ip2].X[i])
				s2 += s1 * s1
				l2 += l1 * l1
				sl += s1 * l1
			}
			s1 = math.Sqrt(s2)
			l1 = math.Sqrt(l2)
			mu = math.Abs(sl) / (s1*l1 + 1.e-15)
			h.Add(w1*p2[ip2].W, s1, mu)
		}
	}

}
