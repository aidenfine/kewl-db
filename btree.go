package main

import (
	"fmt"
	"strings"
)

type Item struct {
	Key   string
	Value string
}
type node struct {
	items    []Item // key, values (sorted)
	children []*node
}

type BTree struct {
	root   *node
	degree int // each node holds at least 2t - 1 items
}

// NewTree (create tree)
func NewBTree(degree int) *BTree {
	if degree <= 1 {
		panic("degree must must be > 1")
	}
	return &BTree{degree: degree}
}

func (t *BTree) maxItems() int {
	return t.degree*2 - 1
}

func (t *BTree) Insert(key, value string) {
	item := Item{Key: key, Value: value}

	// create root if empty
	if t.root == nil {
		t.root = &node{items: []Item{item}}
		return
	}

	// root is full, split
	if len(t.root.items) >= t.maxItems() {
		oldRoot := t.root
		t.root = &node{children: []*node{oldRoot}}
		t.splitChild(t.root, 0)
	}
	t.insertNonFull(t.root, item)
}

func (t *BTree) insertNonFull(n *node, item Item) {
	i := len(n.items) - 1

	if len(n.children) == 0 {
		// leaf — insert item in sorted position
		n.items = append(n.items, Item{})
		for i >= 0 && item.Key < n.items[i].Key {
			n.items[i+1] = n.items[i]
			i--
		}
		n.items[i+1] = item
		return
	}

	// internal node
	for i >= 0 && item.Key < n.items[i].Key {
		i--
	}
	i++

	// if that child is full, split it first
	if len(n.children[i].items) >= t.maxItems() {
		t.splitChild(n, i)
		// after split, decide which of the two children to descend into
		if item.Key > n.items[i].Key {
			i++
		}
	}

	t.insertNonFull(n.children[i], item)
}

// node full needs to be split to avoid walking back upwards, turning one full node into two half full nodes.
func (t *BTree) splitChild(parent *node, i int) {
	full := parent.children[i]
	mid := t.degree - 1

	// new node gets the right half
	right := &node{
		items: make([]Item, len(full.items[mid+1:])),
	}
	copy(right.items, full.items[mid+1:])

	// if not a leaf, move right half of children too
	if len(full.children) > 0 {
		right.children = make([]*node, len(full.children[mid+1:]))
		copy(right.children, full.children[mid+1:])
		full.children = full.children[:mid+1]
	}

	// middle item goes up to parent
	midItem := full.items[mid]
	full.items = full.items[:mid]

	// insert midItem and right child into parent
	parent.items = append(parent.items, Item{})
	copy(parent.items[i+1:], parent.items[i:])
	parent.items[i] = midItem

	parent.children = append(parent.children, nil)
	copy(parent.children[i+2:], parent.children[i+1:])
	parent.children[i+1] = right
}


// print tree (maybe delete later???)
// ty claude <3
func (t *BTree) Print() string {
	if t.root == nil {
		return "(empty tree)"
	}

	var sb strings.Builder
	queue := []*node{t.root}
	level := 0

	for len(queue) > 0 {
		next := []*node{}
		sb.WriteString("Level " + strings.Repeat("", 0) + fmt.Sprintf("%d: ", level))

		for i, n := range queue {
			if i > 0 {
				sb.WriteString(" ")
			}
			keys := make([]string, len(n.items))
			for j, item := range n.items {
				keys[j] = item.Key
			}
			sb.WriteString("[" + strings.Join(keys, ", ") + "]")
			next = append(next, n.children...)
		}

		sb.WriteString("\n")
		queue = next
		level++
	}

	return sb.String()
}
