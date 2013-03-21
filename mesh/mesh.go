package mesh

import (
	"fmt"
	"github.com/npadmana/npgo/math/vector"
	"io"
	"math"
	"os"
	"sort"
	//"github.com/npadmana/go-xi/twopt"
	//"github.com/npadmana/npgo/math/histogram"
	//"log"
)

// This code is NOT for anything but 3D
const (
	NDIM = 3
)

type Particle struct {
	X vector.Vector3D
	W float64
}

func (p Particle) String() string {
	return fmt.Sprintf("(%g,%g,%g),%g", p.X[0], p.X[1], p.X[2], p.W)
}

func (p *Particle) Fscan(ff io.Reader) (int, error) {
	return fmt.Fscan(ff, &p.X[0], &p.X[1], &p.X[2], &p.W)
}

type Mesh struct {
	Particles          []Particle
	Ndx                []int
	Dx                 float64
	Dim, Stride        [NDIM]int
	MinIndex, MaxIndex []int
	BoxDim             vector.Vector3D
	BoxMin             vector.Vector3D
}

func (m *Mesh) Len() int {
	return len(m.Particles)
}

func (m *Mesh) Less(i, j int) bool {
	return m.Ndx[i] < m.Ndx[j]
}

func (m *Mesh) Swap(i, j int) {
	m.Ndx[j], m.Ndx[i] = m.Ndx[i], m.Ndx[j]
	m.Particles[j], m.Particles[i] = m.Particles[i], m.Particles[j]
}

func (m *Mesh) Read(fn string) error {
	// Open file
	ff, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer ff.Close()

	// Read loop
	for err == nil {
		p1 := new(Particle)
		if _, err = p1.Fscan(ff); err == nil {
			m.Particles = append(m.Particles, *p1)
		}
	}
	if err != io.EOF {
		return err
	}
	return nil
}

func New(fn string, dx float64) (*Mesh, error) {

	// Setup
	m := new(Mesh)
	m.Dx = dx

	// Read file
	err := m.Read(fn)
	if err != nil {
		return nil, err
	}

	// Set box dimensions
	m.BoxMin = m.Particles[0].X
	maxpos := m.Particles[0].X
	for _, p1 := range m.Particles {
		m.BoxMin = m.BoxMin.Min(p1.X)
		maxpos = maxpos.Max(p1.X)
	}
	m.BoxDim = maxpos.Sub(m.BoxMin)

	// Set dimensions and strides
	for i, ll := range m.BoxDim {
		m.Dim[i] = int(math.Ceil(ll/m.Dx)) + 1
	}
	// Since we're only doing this in 3 dimensions.
	m.Stride[2] = 1
	m.Stride[1] = m.Dim[2]
	m.Stride[0] = m.Dim[1] * m.Stride[1]
	size := m.Stride[0] * m.Dim[0]

	// Allocate storage for indices
	m.MinIndex = make([]int, size)
	m.MaxIndex = make([]int, size)
	m.Ndx = make([]int, len(m.Particles))

	// Compute indices
	for i, pp := range m.Particles {
		for j, xx := range pp.X {
			m.Ndx[i] += int(math.Floor((xx-m.BoxMin[j])/m.Dx)) * m.Stride[j]
		}
	}

	// Sort the 
	sort.Sort(m)

	return m, nil
}
