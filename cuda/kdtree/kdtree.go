package kdtree

import (
	"github.com/npadmana/go-xi/cuda/cudalib"
	"github.com/npadmana/go-xi/cuda/particle"
	"sort"
	"sync"
)

// Sort along a predin
type particleSorter struct {
	idim int
	arr  particle.ParticleArr
}

func (p *particleSorter) Len() int {
	return len(p.arr)
}

func (p *particleSorter) Less(i, j int) bool {
	if p.arr[i][p.idim] < p.arr[j][p.idim] {
		return true
	}
	return false
}

func (p *particleSorter) Swap(i, j int) {
	p.arr[i], p.arr[j] = p.arr[j], p.arr[i]
}

type Node struct {
	Arr            particle.ParticleArr // Store the slice of the full array
	Npart          int                  // Number of particles
	BoxMin, BoxMax cudalib.Float4       // Really could be a 3 element array, but this is simpler 
	Left, Right    *Node                // Pointers to node
	IsLeaf         bool                 // Is this a leaf?
	Id             int64                // Node id
}

// NewNode returns a new root/leaf node, built from the particle data
func NewNode(arr particle.ParticleArr, id int64) *Node {
	nn := new(Node)

	// Some of these are unnecessary, but it's nice to be explicit
	nn.Arr = arr
	nn.Npart = len(nn.Arr)
	nn.IsLeaf = true
	nn.BoxMin, nn.BoxMax = nn.Arr.MinMax()
	nn.Id = id

	return nn
}

// Grow grows the tree, splitting along the largest dimension. The tree splits are 
// truncated when the new leaves have too few particles (minpart) or the longest box dimension
// is too small. Grow starts jobs concurrently, so we pass in a waitgroup. 
func (n *Node) Grow(minpart int, minbox float32, wg *sync.WaitGroup) {
	defer wg.Done()

	// Check to see if the current node has enough particles to split
	if n.Npart < minpart {
		return
	}

	// Compute the largest box dimension
	var lbox float32
	bigdim := float32(0)
	idim := -1
	for ii := 0; ii < 3; ii++ {
		lbox = n.BoxMax[ii] - n.BoxMin[ii]
		if lbox > bigdim {
			bigdim = lbox
			idim = ii
		}
	}

	// If bigdim too small, don't bother splitting the sample
	if bigdim < minbox {
		return
	}

	// Sort on idim
	ps := particleSorter{idim, n.Arr}
	sort.Sort(&ps)

	// Create left and right nodes
	split := n.Npart/2 + 1
	n.Left = NewNode(n.Arr[0:split], 2*n.Id)
	n.Right = NewNode(n.Arr[split:n.Npart], 2*n.Id+1)

	// Spawn grow on both left and right
	wg.Add(2)
	go n.Left.Grow(minpart, minbox, wg)
	go n.Right.Grow(minpart, minbox, wg)
	// All done
}
