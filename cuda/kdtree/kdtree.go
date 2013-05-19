package kdtree

import (
	"github.com/npadmana/go-xi/cuda/cudalib"
	"github.com/npadmana/go-xi/cuda/particle"
	"math"
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
	Arr            particle.ParticleArr // Store the full array
	Lo, Hi         int                  // Store the slice limits -- usual inclusive/exclusive interval [lo,hi)
	Npart          int                  // Number of particles
	BoxMin, BoxMax cudalib.Float4       // Really could be a 3 element array, but this is simpler 
	Left, Right    *Node                // Pointers to node
	IsLeaf         bool                 // Is this a leaf?
	Id             int64                // Node id
}

// NewNode returns a new root/leaf node, built from the particle data
func NewNode(arr particle.ParticleArr, lo, hi int, id int64) *Node {
	nn := new(Node)

	// Some of these are unnecessary, but it's nice to be explicit
	nn.Arr = arr
	nn.Lo = lo
	nn.Hi = hi
	nn.Npart = hi - lo
	nn.IsLeaf = true
	nn.BoxMin, nn.BoxMax = (nn.Arr[lo:hi]).MinMax()
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
	ps := particleSorter{idim, n.Arr[n.Lo:n.Hi]}
	sort.Sort(&ps)

	// Create left and right nodes
	split := n.Npart/2 + 1
	n.Left = NewNode(n.Arr, n.Lo, n.Lo+split, 2*n.Id)
	n.Right = NewNode(n.Arr, n.Lo+split, n.Hi, 2*n.Id+1)
	n.IsLeaf = false

	// Spawn grow on both left and right
	wg.Add(2)
	go n.Left.Grow(minpart, minbox, wg)
	go n.Right.Grow(minpart, minbox, wg)
	// All done
}

// NodeDist computes the minimum and maximum distances between nodes
func (n *Node) NodeDist(n1 *Node) (mindist, maxdist float32) {
	var x1, x2, dx1, dx2 float32
	for ii := 0; ii < 3; ii++ {
		x1 = (n.BoxMax[ii] + n.BoxMin[ii]) / 2
		dx1 = (n.BoxMax[ii] - n.BoxMin[ii]) / 2
		x2 = (n1.BoxMax[ii] + n1.BoxMin[ii]) / 2
		dx2 = (n1.BoxMax[ii] - n1.BoxMin[ii]) / 2
		x1 = x1 - x2
		if x1 < 0 {
			x1 = -x1
		}
		dx1 += dx2
		dx2 = x1 - dx1
		if dx2 > 0 {
			mindist += dx2 * dx2
		}
		dx2 = x1 + dx1
		maxdist += dx2 * dx2
	}
	mindist = float32(math.Sqrt(float64(mindist)))
	maxdist = float32(math.Sqrt(float64(maxdist)))
	return
}

// NodeList allows one to send and receive groups of nodes.
type NodeList []*Node

// MapDecision determines decision values
type TreeDecision int

const (
	PRUNE    TreeDecision = iota // Prune the tree walk here
	CONTINUE                     // Continue the tree walk 
	EVALUATE                     // Don't continue the tree walk, but evaluate these nodes
)

// This makes decisions while walking the treee
type TreeDecider func(NodeList) TreeDecision

// TreeMap walks a tree, checking to see what nodes match the criterion 
// set by the decider. These are then put onto the channel for future processing
// Again, this is done concurrently, so we send in a WaitGroup as well.
func TreeMap(n1 *Node, ff TreeDecider, out chan NodeList, wg *sync.WaitGroup) {
	defer wg.Done()

	switch ff(NodeList{n1}) {
	case PRUNE:
		return
	case CONTINUE:
		{
			// if n1 is a leaf
			if n1.IsLeaf {
				out <- NodeList{n1}
				return
			}

			wg.Add(2)
			go TreeMap(n1.Left, ff, out, wg)
			go TreeMap(n1.Right, ff, out, wg)
		}
	case EVALUATE:
		out <- NodeList{n1}
	default:
		panic("Unknown decision type")
	}

}

// DualTreeMap walks two trees, checking to see what nodes match the criterion 
// set by the decider. These are then put onto the channel for future processing
// Again, this is done concurrently, so we send in a WaitGroup as well.
func DualTreeMap(n1, n2 *Node, ff TreeDecider, out chan NodeList, wg *sync.WaitGroup) {
	defer wg.Done()

	switch ff(NodeList{n1, n2}) {
	case PRUNE:
		return
	case CONTINUE:
		{
			// Both are leaves, process
			if n1.IsLeaf && n2.IsLeaf {
				out <- NodeList{n1, n2}
				return
			}

			// In all of these cases, we spawn two goroutines
			wg.Add(2)

			// n1 is a leaf
			if n1.IsLeaf {
				go DualTreeMap(n1, n2.Left, ff, out, wg)
				go DualTreeMap(n1, n2.Right, ff, out, wg)
				return
			}

			// n2 is a leaf
			if n2.IsLeaf {
				go DualTreeMap(n1.Left, n2, ff, out, wg)
				go DualTreeMap(n1.Right, n2, ff, out, wg)
				return
			}

			// Neither are leaves
			if n1.Npart > n2.Npart {
				go DualTreeMap(n1.Left, n2, ff, out, wg)
				go DualTreeMap(n1.Right, n2, ff, out, wg)
			} else {
				go DualTreeMap(n1, n2.Left, ff, out, wg)
				go DualTreeMap(n1, n2.Right, ff, out, wg)
			}
		}
	case EVALUATE:
		out <- NodeList{n1, n2}
	default:
		panic("Unknown decision type")
	}

}

type RInterval struct {
	Lo, Hi float32
}

func (r RInterval) DualNodeTest(nn NodeList) TreeDecision {
	min, max := nn[0].NodeDist(nn[1])

	// Prune decisions
	if min > r.Hi {
		return PRUNE
	}
	if max < r.Lo {
		return PRUNE
	}

	return CONTINUE
}
