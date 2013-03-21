package mesh

import (
	"bufio"
	"fmt"
	"github.com/npadmana/go-xi/utils"
	"io"
	"os"
)

type Particle struct {
	X utils.Vector3D
	W float64
}

func (p *Particle) String() string {
	return fmt.Sprintf("(%g,%g,%g),%g", p.X[0], p.X[1], p.X[2], p.W)
}

func countlines(fn string) (n int, err error) {
	n = 0
	ff, err := os.Open(fn)
	if err != nil {
		return 0, err
	}
	defer ff.Close()
	fbuf := bufio.NewReader(ff)

	err = nil
	for err == nil {
		n += 1
		_, err = fbuf.ReadString('\n')
	}
	if err != io.EOF {
		return 0, err
	}
	n -= 1
	return n, nil
}

func ReadParticles(fn string) ([]Particle, error) {
	// Get the number of lines and allocate
	nlines, err := countlines(fn)
	if err != nil {
		return nil, err
	}
	fmt.Println("Expect to get nlines =", nlines)

	// Allocate particle data
	parr := make([]Particle, nlines)

	// Open file
	ff, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer ff.Close()
	fbuf := bufio.NewReader(ff)

	// Read loop
	for ii := 0; ii < nlines; ii++ {
		_, err = fmt.Fscan(fbuf, &parr[ii].X[0], &parr[ii].X[1], &parr[ii].X[2], &parr[ii].W)
		if err != nil {
			return nil, err
		}
	}

	return parr, nil
}
