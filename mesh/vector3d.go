package mesh

import (
	"math"
)

/* Vector 3D is a 3D vector. 
 */
type Vector3D [3]float64

func (v1 Vector3D) Add(v2 Vector3D) (v3 Vector3D) {
	v3[0] = v1[0] + v2[0]
	v3[1] = v1[1] + v2[1]
	v3[2] = v1[2] + v2[2]
	return
}

func (v1 Vector3D) Sub(v2 Vector3D) (v3 Vector3D) {
	v3[0] = v1[0] - v2[0]
	v3[1] = v1[1] - v2[1]
	v3[2] = v1[2] - v2[2]
	return
}

func (v1 Vector3D) Mul(v2 Vector3D) (v3 Vector3D) {
	v3[0] = v1[0] * v2[0]
	v3[1] = v1[1] * v2[1]
	v3[2] = v1[2] * v2[2]
	return
}

func (v1 Vector3D) Div(v2 Vector3D) (v3 Vector3D) {
	v3[0] = v1[0] / v2[0]
	v3[1] = v1[1] / v2[1]
	v3[2] = v1[2] / v2[2]
	return
}

func (v1 Vector3D) Min(v2 Vector3D) (v3 Vector3D) {
	v3[0] = math.Min(v1[0], v2[0])
	v3[1] = math.Min(v1[1], v2[1])
	v3[2] = math.Min(v1[2], v2[2])
	return
}

func (v1 Vector3D) Max(v2 Vector3D) (v3 Vector3D) {
	v3[0] = math.Max(v1[0], v2[0])
	v3[1] = math.Max(v1[1], v2[1])
	v3[2] = math.Max(v1[2], v2[2])
	return
}

func (v1 Vector3D) Dot(v2 Vector3D) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
}

func (v1 Vector3D) Norm() float64 {
	return math.Sqrt(v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2])
}

func (v1 Vector3D) Scale(a float64) (v2 Vector3D) {
	v2[0] = a * v1[0]
	v2[1] = a * v1[0]
	v2[2] = a * v1[0]
	return
}
