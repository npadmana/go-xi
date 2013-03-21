package mesh

import (
	"fmt"
	"github.com/npadmana/npgo/math/vector"
	"io"
	"os"
)

type Particle struct {
	X vector.Vector3D
	W float64
}

func (p *Particle) String() string {
	return fmt.Sprintf("(%g,%g,%g),%g", p.X[0], p.X[1], p.X[2], p.W)
}

func ReadParticles(fn string) ([]Particle, error) {
	var parr []Particle
	var p1 *Particle

	// Open file
	ff, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer ff.Close()

	// Read loop
	for err == nil {
		p1 = new(Particle)
		_, err = fmt.Fscan(ff, &p1.X[0], &p1.X[1], &p1.X[2], &p1.W)
		if err == nil {
			parr = append(parr, *p1)
		}
	}
	if err != io.EOF {
		return nil,err
	}
	return parr, nil
}
