package twopt

import (
	"github.com/npadmana/go-xi/utils"
	"math"
)

type DistFunc func(p1, p2 *utils.Vector3D) []float64

func SMu(p1, p2 *utils.Vector3D) []float64 {
	var s2, l2, sl, s1, l1, mu float64
	for i := 0; i < 3; i++ {
		s1 = p1[i] - p2[i]
		l1 = 0.5*(p1[i]+p2[i])
		s2 += s1 * s1
		l2 += l1 * l1
		sl += s1 * l1
	}
	s1 = math.Sqrt(s2)
	l1 = math.Sqrt(l2)
	if (s1 > 0) && (l1 > 0) {
		mu = math.Abs(sl) / (s1 * l1)
	}

	return []float64{s1, mu}
}
