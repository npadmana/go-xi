package mesh

import (
	"sort"
)

func (m *Mesh) Len() int {
	return len(m.Particles)
}

func (m *Mesh) Less(i, j int) bool {
	return m.Ndx[i] < m.Ndx[j]
}

func (m *Mesh) flat2three(ndx int) (out Index3D) {
	for i, s := range m.Stride {
		out[i] = ndx / s
		ndx = ndx % s
	}
	return
}

func (m *Mesh) Swap(i, j int) {
	m.Ndx[j], m.Ndx[i] = m.Ndx[i], m.Ndx[j]
	m.Particles[j], m.Particles[i] = m.Particles[i], m.Particles[j]
}

func (m *Mesh) DoSort() {
	// Sort the particle data
	sort.Sort(m)
	
	// Set up the Grid
	m.Grid = make([]*GridPoint, m.Ngrid)

	// Now set the indices
	ndxprev := m.Ndx[0]
	min := 0
	m.Grid[ndxprev] = &GridPoint{N: ndxprev, I : m.flat2three(0)}
	for ii, ndx1 := range m.Ndx {
		if ndx1 != ndxprev {
			m.Grid[ndxprev].P = m.Particles[min:ii]
			m.Grid[ndx1] = &GridPoint{N: ndx1, I: m.flat2three(ndx1)}
			min = ii
			ndxprev = ndx1

		}
	}
	m.Grid[ndxprev].P = m.Particles[min:m.Npart]
}
