/* Package particle defines a Particle struct and routines to manipulate these
 */
package particle

import (
	"fmt"
	"github.com/npadmana/npgo/textio"
	"github.com/npadmana/npgo/vector"
)

type Particle struct {
	X vector.Vector3D
	W float64
}

func (p Particle) String() string {
	return fmt.Sprintf("(%g,%g,%g),%g", p.X[0], p.X[1], p.X[2], p.W)
}

// ParticleArr is a storage container for Particles
type Particles struct {
	Data   []Particle
	BoxDim vector.Vector3D
	Origin vector.Vector3D
}

func NewFromXYZW(fn string) (*Particles, error) {
	var parr Particles
	var err error
	var p1 Particle

	out := make(chan textio.Line, 100)
	go textio.FileLineReader(fn, out)
	for l1 := range out {
		if l1.Err != nil {
			return &Particles{}, l1.Err
		}
		_, err = fmt.Sscan(l1.Str, &p1.X[0], &p1.X[1], &p1.X[2], &p1.W)
		parr.Data = append(parr.Data, p1)
		if err != nil {
			return &Particles{}, err
		}
	}

	return &parr, nil
}

func (p *Particles) Normalize() {
	if len(p.Data) == 0 {
		return
	}

	minpos := p.Data[0].X
	maxpos := p.Data[0].X

	for _, p1 := range p.Data {
		minpos = minpos.Min(p1.X)
		maxpos = maxpos.Max(p1.X)
	}

	p.Origin = vector.Vector3D{0,0,0}.Sub(minpos)
	p.BoxDim = maxpos.Sub(minpos)

	for i := range p.Data {
		p.Data[i].X = p.Data[i].X.Sub(minpos)
	}

}
