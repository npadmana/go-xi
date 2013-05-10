package cudalib

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L../build -lgpucount

#include <stdlib.h>
#include "gpucount.h"

*/
import "C"

import (
	"unsafe"
)

var (
	sizef4  = int(C.sizef4())
	sizeull = int(C.sizeull())
)

type Float4 [4]float32

// Min implements f1 = f2 < f1 ? f2 : f1 componentwise
func (f1 *Float4) Min(f2 Float4) {
	for ii := 0; ii < 4; ii++ {
		if f2[ii] < f1[ii] {
			f1[ii] = f2[ii]
		}
	}
}

// Max implements f1 = f2 > f1 ? f2 : f1 componentwise
func (f1 *Float4) Max(f2 Float4) {
	for ii := 0; ii < 4; ii++ {
		if f2[ii] > f1[ii] {
			f1[ii] = f2[ii]
		}
	}
}

// Check error checks the last CUDA error; if an error occurs, it aborts.
func CheckError(msg string) {
	cstr := C.CString(msg)
	C.checkCUDAError(cstr)
	C.free(unsafe.Pointer(cstr)) // This line will almost never be called
}

// Device particle structure
type DevParticle struct {
	ptr   unsafe.Pointer
	npart int
}

// Allocate particle data
func NewDevParticle(npart int) *DevParticle {
	if sizef4 != int(unsafe.Sizeof(Float4{})) {
		panic("Incorrect sized float4")
	}
	pp := new(DevParticle)
	pp.npart = npart
	C.allocDev(C.long(npart*sizef4), &pp.ptr)
	CheckError("Allocating particle data")
	return pp
}

// Free Particle data
func (pp *DevParticle) Free() {
	pp.npart = 0
	C.freeDev(pp.ptr)
	pp.ptr = nil
}

// CopyToDevice moves particle data onto the device
func (pp *DevParticle) CopyToDevice(ff []Float4) {
	C.copyToDevice(pp.ptr, unsafe.Pointer(&ff[0]), C.long(pp.npart*sizef4))
	CheckError("Error moving particle data to GPU")
}

// Histogram structure
type DevHist struct {
	ptr   unsafe.Pointer
	nbins int
}

// Allocate histogram
func NewDevHist(nbins int) *DevHist {
	if sizeull != int(unsafe.Sizeof(uint64(0))) {
		panic("Incorrect sized unsigned long long")
	}
	h := new(DevHist)
	h.nbins = nbins
	C.allocDev(C.long(nbins*sizeull), &h.ptr)
	CheckError("Allocating histogram data")
	return h
}

// Free histogram data
func (h *DevHist) Free() {
	h.nbins = 0
	C.freeDev(h.ptr)
	h.ptr = nil
}

// Copy histogram back from device; the slice you send is what will 
// be filled.
func (h *DevHist) CopyFromDevice(h1 []uint64) {
	C.copyFromDevice(unsafe.Pointer(&h1[0]), h.ptr, C.long(h.nbins*sizeull))
	CheckError("Error moving histogram data off device")
}

// Copy histogram to device; the slice you send is copied over
func (h *DevHist) CopyToDevice(h1 []uint64) {
	C.copyToDevice(h.ptr, unsafe.Pointer(&h1[0]), C.long(h.nbins*sizeull))
	CheckError("Error moving histogram data to device")
}
