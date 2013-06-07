package twopt

/*
#cgo CFLAGS: -I . 

#include "counters.h"

*/
import "C"

import (
	"github.com/npadmana/go-xi/mesh"
	"unsafe"
)

func smucount(p1, p2 mesh.ParticleArr, s *SMuPairCounter, scale float64) {
	n1 := len(p1)
	n2 := len(p2)
	maxs2 := s.Maxs * s.Maxs
	invdmu := 1 / s.Dmu
	invds := 1 / s.Ds

	C.smu(unsafe.Pointer(&p1[0]), unsafe.Pointer(&p2[0]),
		C.int(n1), C.int(n2),
		*C.double(&s.Data[0]), C.int(s.Nmu),
		C.double(maxs2), C.double(invdmu), C.double(invds),
		C.double(scale))

}
