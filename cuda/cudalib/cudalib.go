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
	sizef4  = C.sizeof_f4()
	sizeull = C.sizeof_ull()
)

// Float4
type Float4 struct {
	X, Y, Z, W float32
}

// Min implements f1 = f2 < f1 ? f2 : f1 componentwise
func (f1 *Float4) Min(f2 Float4) {
	if f2.X < f1.X {
		f1.X = f2.X
	}
	if f2.Y < f1.Y {
		f1.Y = f2.Y
	}
	if f2.Z < f1.Z {
		f1.Z = f2.Z
	}
	if f2.W < f1.W {
		f1.W = f2.W
	}
}

// Max implements f1 = f2 > f1 ? f2 : f1 componentwise
func (f1 *Float4) Max(f2 Float4) {
	if f2.X > f1.X {
		f1.X = f2.X
	}
	if f2.Y > f1.Y {
		f1.Y = f2.Y
	}
	if f2.Z > f1.Z {
		f1.Z = f2.Z
	}
	if f2.W > f1.W {
		f1.W = f2.W
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
	npart int32
}

// Allocate particle data
func NewDevParticle(npart int32) *DevParticle {
	if sizef4 != unsafe.Sizeof(Float4{}) {
		panic("Incorrect sized float4")
	}
	if unsafe.Sizeof(
	pp := new(DevParticle)
	pp.npart = npart
	pp.ptr = C.allocDev(C.long(npart * sizef4))
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
	C.copyToDevice(pp.ptr, unsafe.Pointer(&ff[0]), pp.npart*sizef4)
	CheckError("Error moving particle data to GPU")
}

// Histogram structure
type DevHist struct {
	ptr   unsafe.Pointer
	nbins int32
}

// Allocate histogram
func NewDevHist(nbins int32) *DevHist {
	if sizeull != unsafe.Sizeof(uint64{}) {
		panic("Incorrect sized unsigned long long")
	}
	h := new(DevHist)
	h.nbins = nbins
	h.ptr = C.allocDev(C.long(nbins * sizeull))
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
	C.copyFromDevice(unsafe.Pointer(&h1[0]),h.ptr,h.nbins*sizeull)
	CheckError("Error moving histogram data off device")
}
