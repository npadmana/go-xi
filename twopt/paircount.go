package twopt

import (
	"github.com/npadmana/go-xi/mesh"
	"github.com/npadmana/go-xi/utils"
	"math"
)

// PairCounter is a basic interface for the paircounting codes
type PairCounter interface {
	Count([]mesh.Particle, []mesh.Particle, float64)
	Add(PairCounter)
	Get(is, imu int) float64
}

// SMuPairCounter 
type SMuPairCounter struct {
	Data          []float64
	Nmu, Ns       int
	Dmu, Ds, Maxs float64
}

// NewSMuPairCounter returns an SMu paircounter
func NewSMuPairCounter(Ns, Nmu int, Maxs float64) (smu *SMuPairCounter) {
	smu = new(SMuPairCounter)
	smu.Ns = Ns
	smu.Nmu = Nmu
	smu.Maxs = Maxs
	smu.Dmu = 1 / float64(Nmu)
	smu.Ds = Maxs / float64(Ns)
	smu.Data = make([]float64, Ns*Nmu)

	return
}

// PairCounter implements the inner pair counting loop.
func (smu *SMuPairCounter) Count(p1, p2 []mesh.Particle, scale float64) {
	var s2, l2, sl, s1, l1, mu, w1 float64
	var imu, is, i int
	var x1 utils.Vector3D

	n1 := len(p1)
	n2 := len(p2)
	var ip1, ip2 int
	
	maxs2 := smu.Maxs*smu.Maxs

	for ip1 = 0; ip1 < n1; ip1++ {
		x1 = p1[ip1].X
		w1 = p1[ip1].W * scale
		for ip2 = 0; ip2 < n2; ip2++ {
//			s2 = (x1[0]-p2[ip2].X[0])*(x1[0]-p2[ip2].X[0]) + 
//				 (x1[1]-p2[ip2].X[1])*(x1[1]-p2[ip2].X[1]) +
//				 (x1[2]-p2[ip2].X[2])*(x1[2]-p2[ip2].X[2])
//			l2 = (x1[0]+p2[ip2].X[0])*(x1[0]+p2[ip2].X[0]) + 
//				 (x1[1]+p2[ip2].X[1])*(x1[1]+p2[ip2].X[1]) +
//				 (x1[2]+p2[ip2].X[2])*(x1[2]+p2[ip2].X[2])
//			sl = (x1[0]-p2[ip2].X[0])*(x1[0]+p2[ip2].X[0]) + 
//				 (x1[1]-p2[ip2].X[1])*(x1[1]+p2[ip2].X[1]) +
//				 (x1[2]-p2[ip2].X[2])*(x1[2]+p2[ip2].X[2])
			s2, l2, sl = 0, 0, 0
			for i = 0; i < 3; i++ {
				s1 = x1[i] - p2[ip2].X[i]
				l1 = 0.5 * (x1[i] + p2[ip2].X[i])
				s2 += s1 * s1
				l2 += l1 * l1
				sl += s1 * l1
			}
			if s2 < maxs2 {
				s1 = math.Sqrt(s2)
				l1 = math.Sqrt(l2)
				mu = sl / (s1*l1 + 1.e-15)
				if mu < 0 {
					mu = -mu
				}
				imu = int(mu / smu.Dmu)
				is = int(s1 / smu.Ds)
				smu.Data[is*smu.Nmu+imu] += w1 * p2[ip2].W
			}
		}
	}

}

func (smu *SMuPairCounter) Get(is, imu int) float64 {
	return smu.Data[is*smu.Nmu+imu]
}

func (smu *SMuPairCounter) Add(src PairCounter) {
	switch src1 := src.(type) {
	case *SMuPairCounter:
		for i, val := range src1.Data {
			smu.Data[i] += val
		}
	}
}
