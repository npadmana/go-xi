/* Package particle defines a Particle struct and routines to manipulate these
 */
package particle

import (
	"fmt"
	"github.com/npadmana/npgo/textio"
	"github.com/npadmana/npgo/math/vector"
)

type Particle struct {
	X vector.Vector3D
	W float64
}

func (p Particle) String() string {
	return fmt.Sprintf("(%g,%g,%g),%g", p.X[0], p.X[1], p.X[2], p.W)
}

// ParticleArr is a storage container for Particles
type Particles []Particle

func NewFromXYZW(fn string) (Particles, error) {
	var parr Particles
	var err error
	var p1 Particle

	out := make(chan textio.Line, 100)
	go textio.FileLineReader(fn, out)
	for l1 := range out {
		if l1.Err != nil {
			return nil, l1.Err
		}
		_, err = fmt.Sscan(l1.Str, &p1.X[0], &p1.X[1], &p1.X[2], &p1.W)
		parr = append(parr, p1)
		if err != nil {
			return nil, err
		}
	}

	return parr, nil
}

