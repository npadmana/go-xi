package mesh

import (
	"math"
	//"github.com/npadmana/go-xi/twopt"
)

const (
	GeneratorBuffer = 100 // The buffer size for our generators
)

type Index3D [3]int

type Mesh struct {
	Particles   ParticleArr  // Storage for particle data -- size Nparticles
	Npart       int          // Number of particles
	Ndx         []int        // Storage for particle grid index -- size Nparticles
	Dx          float64      // Grid spacing
	Dim, Stride Index3D      // Grid dimensions and strides
	Ngrid       int          // Number of grid points
	Grid        []*GridPoint // Slice containing pointers to gridpoints with data
}

type GridPoint struct {
	N int
	I Index3D
	P ParticleArr
}

func New(p ParticleArr, dx float64, boxmin, boxmax Vector3D) (m *Mesh) {
	var x Vector3D

	// Setup
	m = new(Mesh)
	m.Dx = dx
	m.Particles = p
	m.Npart = len(p)

	// Set dimensions and strides
	m.Ngrid = 1
	boxdim := boxmax.Sub(boxmin)
	for i, ll := range boxdim {
		m.Dim[i] = int(math.Ceil(ll/m.Dx)) + 1
		m.Ngrid *= m.Dim[i]
	}
	tmp := m.Ngrid
	for i, ll := range m.Dim {
		m.Stride[i] = tmp / ll
		tmp /= ll
	}

	m.Ndx = make([]int, m.Npart)

	// Compute indices
	for i, pp := range m.Particles {
		x = pp.X
		for j, x1 := range x {
			m.Ndx[i] += int(math.Floor((x1-boxmin[j])/m.Dx)) * m.Stride[j]
		}
	}

	m.DoSort()

	return
}

func (m *Mesh) LoopAll() chan *GridPoint {
	c := make(chan *GridPoint, GeneratorBuffer)
	go func() {
		for _, g := range m.Grid {
			if g != nil {
				c <- g
			}
		}
		close(c)
	}()
	return c
}

func (m *Mesh) LoopNear(p Index3D, dist float64) chan *GridPoint {
	c := make(chan *GridPoint, GeneratorBuffer)
	// Compute the index to search over
	ndx := int(math.Ceil(dist / m.Dx))

	// Sanity checks and deal with degenerate case 
	diff := float64(ndx)*m.Dx - dist
	if diff < -1.e-2 {
		panic("Failure in search radius")
	}
	if math.Abs(diff) < 1.e-2 {
		ndx++
	}

	var imin, imax Index3D

	for i, p1 := range p {
		imin[i] = p1 - ndx
		if imin[i] < 0 {
			imin[i] = 0
		}
		imax[i] = p1 + ndx + 1 // Exclusive bound
		if imax[i] > m.Dim[i] {
			imax[i] = m.Dim[i]
		}
	}

	go func() {
		pos := 0
		for i1 := imin[0]; i1 < imax[0]; i1++ {
			for j1 := imin[1]; j1 < imax[1]; j1++ {
				for k1 := imin[2]; k1 < imax[2]; k1++ {
					pos = i1*m.Stride[0] + j1*m.Stride[1] + k1*m.Stride[2]
					if m.Grid[pos] != nil {
						c <- m.Grid[pos]
					}
				}
			}
		}
		close(c)
	}()
	return c
}
