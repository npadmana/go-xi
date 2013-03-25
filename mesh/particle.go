package mesh

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
)

type Particle struct {
	X Vector3D
	W float64
}

type ParticleArr []Particle

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

func ReadParticles(fn string, subsample float64) (ParticleArr, error) {
	// Get the number of lines and allocate
	nlines, err := countlines(fn)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Expecting to get %d nlines from %s (if not subsampled)...\n", nlines, fn)

	// Allocate particle data
	parr := make(ParticleArr, nlines)

	// Open file
	ff, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer ff.Close()
	fbuf := bufio.NewReader(ff)

	if subsample < 1 {
		fmt.Println("Subsampling down by ", subsample)
	}
	// Read loop
	iline := 0
	for ii := 0; ii < nlines; ii++ {
		_, err = fmt.Fscan(fbuf, &parr[iline].X[0], &parr[iline].X[1], &parr[iline].X[2], &parr[iline].W)
		if err != nil {
			return nil, err
		}
		if rand.Float64() < subsample {
			iline++
		}
	}
	fmt.Printf("%d particles read in from %s\n", iline, fn)

	return parr[0:iline], nil
}

func (p ParticleArr) MinMax() (boxmin, boxmax Vector3D) {
	boxmin = p[0].X
	boxmax = p[0].X
	for _, p1 := range p {
		boxmin = boxmin.Min(p1.X)
		boxmax = boxmax.Max(p1.X)
	}
	return
}
