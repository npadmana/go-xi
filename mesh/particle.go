package mesh

import (
	"bufio"
	"fmt"
	"github.com/npadmana/go-xi/utils"
	"io"
	"math/rand"
	"os"
)

type Particle struct {
	X utils.Vector3D
	W float64
	R float64
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

func ReadParticles(fn string, subsample float64) ([]Particle, error) {
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

	fmt.Println("Subsampling down by ", subsample)
	// Read loop
	iline := 0
	for ii := 0; ii < nlines; ii++ {
		_, err = fmt.Fscan(fbuf, &parr[iline].X[0], &parr[iline].X[1], &parr[iline].X[2], &parr[iline].W)
		if err != nil {
			return nil, err
		}
		if rand.Float64() < subsample {
			parr[iline].R = parr[iline].X.Norm()
			iline++
		}
	}
	fmt.Println("Final size =", iline)

	return parr[0:iline], nil
}
