package main

import (
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	tree := NewBTree(2)

	tree.Insert("c", "3")
	tree.Insert("a", "1")
	tree.Insert("b", "2")
	tree.Insert("d", "4")
	tree.Insert("e", "5")
	tree.Insert("f", "6")
	tree.Insert("g", "7")
	tree.Insert("h","-100")
	tree.Insert("z", "-52")

	fmt.Println(tree.Print())
}
