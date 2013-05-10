package cudalib

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L../build -lgpucount

#include <stdlib.h>
#include "gpucount.h"

*/
import "C"

import (
	"errors"
	"fmt"
	"io"
)

// PairCounter is a basic interface for the paircounting codes
type CudaPairCounter interface {
	Count(*DevParticle, *DevParticle, int, int, int, int, float32)
	Get(is, imu int) float64
	PPrint(io.Writer)
	PullFromDevice()
}

// SMuPairCounter 
type SMuCudaPairCounter struct {
	Data          []uint64
	Nmu, Ns       int
	Dmu, Ds, Maxs float32
	IntScale      float32
	devh          *DevHist
}

// NewSMuPairCounter returns an SMu paircounter
func NewSMuCudaPairCounter(Ns, Nmu int, Maxs float32, IntScale float32) (smu *SMuCudaPairCounter) {
	smu = new(SMuCudaPairCounter)
	smu.Ns = Ns
	smu.Nmu = Nmu
	smu.Maxs = Maxs
	smu.IntScale = IntScale
	smu.Dmu = 1 / float32(Nmu)
	smu.Ds = Maxs / float32(Ns)
	smu.Data = make([]uint64, Ns*Nmu)
	for ii := 0; ii < Ns*Nmu; ii++ {
		smu.Data[ii] = 0
	}

	smu.devh = NewDevHist(Ns * Nmu)
	smu.devh.CopyToDevice(smu.Data)

	return
}

func (smu *SMuCudaPairCounter) Free() {
	smu.devh.Free()
}

func (smu *SMuCudaPairCounter) Get(is, imu int) float64 {
	return float64(smu.Data[is*smu.Nmu+imu]) / float64(smu.IntScale)
}

func (smu *SMuCudaPairCounter) PPrint(ff io.Writer) {
	for i := 0; i <= smu.Ns; i++ {
		fmt.Fprintf(ff, "%.3f ", float32(i)*smu.Ds)
	}
	fmt.Fprintln(ff)
	for i := 0; i <= smu.Nmu; i++ {
		fmt.Fprintf(ff, "%.3f ", float32(i)*smu.Dmu)
	}
	fmt.Fprintln(ff)

	for i := 0; i < smu.Ns; i++ {
		for j := 0; j < smu.Nmu; j++ {
			fmt.Fprintf(ff, "%25.15e ", smu.Get(i, j))
		}
		fmt.Fprintln(ff)
	}

}

func (smu *SMuCudaPairCounter) PullFromDevice() {
	smu.devh.CopyFromDevice(smu.Data)
}

func (smu *SMuCudaPairCounter) Count(p1, p2 *DevParticle, lo1, hi1, lo2, hi2 int, fac float32) error {
	ii := C.smu(p1.ptr, C.int(lo1), C.int(hi1), p2.ptr, C.int(lo2), C.int(hi2), C.float(fac*smu.IntScale),
		C.int(smu.Ns), C.int(smu.Nmu), C.float(1/smu.Ds), smu.devh.ptr, 7, 4, 16, 32)
	if ii != 0 {
		return errors.New("An error occurred!")
	}
	return nil
}
