/* Package particle defines a Particle struct and routines to manipulate these
 */
package particle

import (
	"fmt"
	"github.com/npadmana/npgo/textio"
	//"strings"
)

type Particle struct {
	Pos    [3]float64
	Weight float64
}

func (p Particle) String() string {
	return fmt.Sprintf("(%g,%g,%g),%g", p.Pos[0], p.Pos[1], p.Pos[2], p.Weight)
}

// ParticleArr is a storage container for Particles
type ParticleArr []Particle

func NewFromXYZW(fn string) (ParticleArr, error) {
	var parr ParticleArr
	var err error
	var p1 Particle

	out := make(chan textio.Line, 100)
	go textio.FileLineReader(fn, out)
	for l1 := range out {
		if l1.Err != nil {
			return nil, l1.Err
		}
		_, err = fmt.Sscan(l1.Str, &p1.Pos[0], &p1.Pos[1], &p1.Pos[2], &p1.Weight)
		parr = append(parr, p1)
		if err != nil {
			return nil, err
		}
	}

	return parr, nil
}
