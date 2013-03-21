package utils

import (
	"math"
)

type Histogrammer interface {
	Reset()
	Add(float64, ...float64)
	Get(...int) float64
}

// UniformHistogram is a histogram type
type UniformHistogram struct {
	Ndim          int
	Nbins, Stride []int
	Min, Dx, Data []float64
}

// NewUniform creates a new UniformHistogram.
// The histogram has nbins bins running from min to max
func NewUniform(nbins []int, min, max []float64) (h *UniformHistogram) {
	// Create a new histogram
	h = new(UniformHistogram)

	// Copy parameters
	h.Ndim = len(nbins)
	h.Nbins = make([]int, h.Ndim)
	copy(h.Nbins, nbins)
	h.Min = make([]float64, h.Ndim)
	copy(h.Min, min)
	h.Dx = make([]float64, h.Ndim)
	for i := 0; i < h.Ndim; i++ {
		h.Dx[i] = (max[i] - min[i]) / float64(h.Nbins[i])
	}

	// Set data and strides
	h.Stride = make([]int, h.Ndim)
	h.Stride[h.Ndim-1] = 1
	for i := h.Ndim - 1; i > 0; i-- {
		h.Stride[i-1] = h.Stride[i] * h.Nbins[i]
	}
	h.Data = make([]float64, h.Stride[0]*h.Nbins[0])

	return
}

// Reset zeros the histogram
func (h *UniformHistogram) Reset() {
	for i := range h.Data {
		h.Data[i] = 0.0
	}
}

// Index returns the N-D and flattened index corresponding to f.
// If f does not lie within the histogram ndx is nil, and flat = -1
func (h *UniformHistogram) Index(f ...float64) (ndx []int, flat int) {
	if len(f) != h.Ndim {
		panic("Incorrect dimension of point")
	}

	ndx = make([]int, h.Ndim)
	flat = 0
	for i, dx := range h.Dx {
		ndx[i] = int(math.Floor((f[i] - h.Min[i]) / dx))
		if (ndx[i] < 0) || (ndx[i] >= h.Nbins[i]) {
			return nil, -1
		}
		flat = flat + ndx[i]*h.Stride[i]
	}
	return
}

// Add weight to the histogram at point f. 
// If f is out of bounds, ignore the value
func (h *UniformHistogram) Add(weight float64, f ...float64) {
	_, i := h.Index(f...)
	if i != -1 {
		h.Data[i] += weight
	}
}

// Get returns the current value of the bin at index ndx
func (h *UniformHistogram) Get(ndx ...int) float64 {
	if len(ndx) != h.Ndim {
		panic("Incorrect dimension of point")
	}

	flat := 0
	for i, s := range h.Stride {
		flat = flat + ndx[i]*s
	}
	return h.Data[flat]
}

// Bins returns the bin edges along dimension idim
func (h *UniformHistogram) Bins(idim int) (b []float64) {
	b = make([]float64, h.Nbins[idim]+1)
	for i := 0; i <= h.Nbins[idim]; i++ {
		b[i] = h.Min[idim] + h.Dx[idim]*float64(i)
	}
	return
}

// AddHist adds histogram h1 to the receiver h 
// The loop runs over the receiver elements
func (h *UniformHistogram) AddHist(h1 *UniformHistogram) {
	if len(h.Data) != len(h1.Data) {
		panic("Incompatible histograms")
	}

	for i := range h.Data {
		h.Data[i] = h.Data[i] + h1.Data[i]
	}
}
