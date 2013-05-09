package mesh

import (
	"bufio"
	"fmt"
	"github.com/npadmana/go-xi/cuda/cudalib"
	"io"
	"math/rand"
	"os"
)

// This is an almost direct copy of the original mesh package, with minor tweaks for
// the cuda code. We'll keep these separate for now, to reduce dependencies.

type ParticleArr []cudalib.Float4

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

func ReadParticles(fn string, subsample float32) (ParticleArr, error) {
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
		// Read in the first four columns
		_, err = fmt.Fscan(fbuf, &parr[iline][0], &parr[iline][1], &parr[iline][2], &parr[iline][3])
		if err != nil {
			return nil, err
		}
		if rand.Float32() < subsample {
			iline++
		}
		// Now fast-forward to the end of the line, discarding the output.
		_, err = fbuf.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
	}
	fmt.Printf("%d particles read in from %s\n", iline, fn)

	return parr[0:iline], nil
}

func (p ParticleArr) MinMax() (boxmin, boxmax cudalib.Float4) {
	boxmin = p[0]
	boxmax = p[0]
	for _, p1 := range p {
		(&boxmin).Min(p1)
		(&boxmax).Max(p1)
	}
	return
}
