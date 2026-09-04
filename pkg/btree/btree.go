package btree

import (
	"fmt"
	"strings"
)

type Item struct {
	Key   string
	Value string
}

type Node struct {
	items    []Item
	children []uint32
	pageID   uint32
}

type BTree struct {
	root   *Node
	pager  *Pager
	degree int
}

func NewBTree(degree int, filename string) (*BTree, error) {
	if degree <= 1 {
		panic("degree must be > 1")
	}

	p, err := NewPager(filename)
	if err != nil {
		return nil, err
	}

	t := &BTree{degree: degree, pager: p}

	if p.nextPage > 0 {
		t.root = t.getNode(0)
	}

	return t, nil
}

func (t *BTree) Close() error {
	return t.pager.Close()
}

func (t *BTree) maxItems() int {
	return t.degree*2 - 1
}

func (t *BTree) getNode(pageID uint32) *Node {
	data, err := t.pager.ReadPage(pageID)
	if err != nil {
		panic(err)
	}
	return deserializeNode(data, pageID)
}

func (t *BTree) writeNode(n *Node) {
	data := serializeNode(n)
	err := t.pager.WritePage(n.pageID, data)
	if err != nil {
		panic(err)
	}
}

func (t *BTree) Insert(key, value string) {
	item := Item{Key: key, Value: value}

	if t.root == nil {
		pageID := t.pager.Allocate()
		t.root = &Node{items: []Item{item}, pageID: pageID}
		t.writeNode(t.root)
		return
	}

	if len(t.root.items) >= t.maxItems() {
		oldRoot := t.root
		newRootPage := t.pager.Allocate()
		t.root = &Node{children: []uint32{oldRoot.pageID}, pageID: newRootPage}
		t.splitChild(t.root, 0)
	}

	t.insertNonFull(t.root, item)
}

func (t *BTree) insertNonFull(n *Node, item Item) {
	i := len(n.items) - 1

	if len(n.children) == 0 {
		n.items = append(n.items, Item{})
		for i >= 0 && item.Key < n.items[i].Key {
			n.items[i+1] = n.items[i]
			i--
		}
		n.items[i+1] = item
		t.writeNode(n)
		return
	}

	for i >= 0 && item.Key < n.items[i].Key {
		i--
	}
	i++

	child := t.getNode(n.children[i])
	if len(child.items) >= t.maxItems() {
		t.splitChild(n, i)
		if item.Key > n.items[i].Key {
			i++
		}
		child = t.getNode(n.children[i])
	}

	t.insertNonFull(child, item)
}

func (t *BTree) GetItemByKey(key string) (Item, bool) {
	n := t.root
	for n != nil {
		i := 0
		for i < len(n.items) && key > n.items[i].Key {
			i++
		}

		if i < len(n.items) && key == n.items[i].Key {
			return n.items[i], true
		}

		if len(n.children) == 0 {
			return Item{}, false
		}

		n = t.getNode(n.children[i])
	}

	return Item{}, false
}

func (t *BTree) splitChild(parent *Node, i int) {
	full := t.getNode(parent.children[i])
	mid := t.degree - 1

	rightPage := t.pager.Allocate()
	right := &Node{
		items:  make([]Item, len(full.items[mid+1:])),
		pageID: rightPage,
	}
	copy(right.items, full.items[mid+1:])

	if len(full.children) > 0 {
		right.children = make([]uint32, len(full.children[mid+1:]))
		copy(right.children, full.children[mid+1:])
		full.children = full.children[:mid+1]
	}

	midItem := full.items[mid]
	full.items = full.items[:mid]

	parent.items = append(parent.items, Item{})
	copy(parent.items[i+1:], parent.items[i:])
	parent.items[i] = midItem

	parent.children = append(parent.children, 0)
	copy(parent.children[i+2:], parent.children[i+1:])
	parent.children[i+1] = rightPage

	t.writeNode(full)
	t.writeNode(right)
	t.writeNode(parent)
}

func (t *BTree) Print() string {
	if t.root == nil {
		return "(empty tree)"
	}

	var sb strings.Builder
	queue := []uint32{t.root.pageID}
	level := 0

	for len(queue) > 0 {
		var next []uint32
		sb.WriteString(fmt.Sprintf("Level %d: ", level))

		for i, pageID := range queue {
			if i > 0 {
				sb.WriteString(" ")
			}
			n := t.getNode(pageID)
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
