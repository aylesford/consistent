package consistent

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

var (
	// Represents how many virtual nodes are replicated on the hash ring in default.
	DefaultReplicas = 128

	// To convert items and keys into index.
	DefaultHashFunc = crc32.ChecksumIEEE
)

// option provides the model function to set the consistent hash circle.
type option func(*Consistent)

// WithReplicas return a function to set the replicas of the consistent hash circle.
func WithReplicas(r int) option {
	return func(c *Consistent) {
		c.replicas = r
	}
}

// WithHashFunc return a function to set the hashFunc of the consistent hash circle.
func WithHashFunc(f func([]byte) uint32) option {
	return func(c *Consistent) {
		c.hashFunc = f
	}
}

// Consistent holds the information about the members of the consistent hash circle.
type Consistent struct {
	replicas int
	nodes    []int
	circle   map[int]string
	hashFunc func([]byte) uint32
	sync.RWMutex
}

// New returns a consistent hash circle according to your ideas.
func New(options ...option) *Consistent {
	c := &Consistent{
		replicas: DefaultReplicas,
		nodes:    make([]int, 0, 128*2),
		circle:   make(map[int]string),
		hashFunc: DefaultHashFunc,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

// Add adds one or more items into the consistent hash circle.
func (c *Consistent) Add(items ...string) {
	c.Lock()
	defer c.Unlock()

	for _, item := range items {
		for i := 0; i < c.replicas; i++ {
			node := int(c.hashFunc([]byte(item + strconv.Itoa(i))))
			c.nodes = append(c.nodes, node)
			c.circle[node] = item
		}
	}

	sort.Ints(c.nodes)
}

// Get returns a item according to the key.
func (c *Consistent) Get(key string) string {
	c.RLock()
	defer c.RUnlock()

	hkey := int(c.hashFunc([]byte(key)))
	index := sort.Search(len(c.nodes), func(i int) bool { return hkey <= c.nodes[i] })

	if index == len(c.nodes) {
		index = 0
	}

	return c.circle[c.nodes[index]]
}
