/* Package particle defines a Particle struct and routines to manipulate these
 */
package particle

import (
	"fmt"
	"github.com/npadmana/npgo/textio"
	"math"
)

type Pos [3]float64

func Min(p1, p2 Pos) (m Pos) {
	for i := range p1 {
		m[i] = math.Min(p1[i], p2[i])
	}
	return
}

func Max(p1, p2 Pos) (m Pos) {
	for i := range p1 {
		m[i] = math.Max(p1[i], p2[i])
	}
	return
}

func Diff(p1, p2 Pos) (m Pos) {
	for i := range p1 {
		m[i] = p1[i] - p2[i]
	}
	return
}

type Particle struct {
	X Pos
	W float64
}

func (p Particle) String() string {
	return fmt.Sprintf("(%g,%g,%g),%g", p.X[0], p.X[1], p.X[2], p.W)
}

// ParticleArr is a storage container for Particles
type Particles struct {
	Data   []Particle
	BoxDim Pos
	Origin Pos
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
	minpos := Pos{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	maxpos := Pos{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}

	for _, p1 := range p.Data {
		minpos = Min(minpos, p1.X)
		maxpos = Max(maxpos, p1.X)
	}

	p.Origin = Diff(Pos{0, 0, 0}, minpos)
	p.BoxDim = Diff(maxpos, minpos)

	for i := range p.Data {
		p.Data[i].X = Diff(p.Data[i].X, minpos)
	}

}
