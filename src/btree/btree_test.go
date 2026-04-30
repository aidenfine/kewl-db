package btree_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aidenfine/kewl-db/src/btree"
)

func TestInsert(t *testing.T) {
	os.Remove("test.db")
	defer os.Remove("test.db")

	tree, err := btree.NewBTree(2, "test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer tree.Close()

	tree.Insert("c", "3")
	tree.Insert("a", "1")
	tree.Insert("b", "2")
	tree.Insert("d", "4")
	tree.Insert("e", "5")
	tree.Insert("f", "6")
	tree.Insert("g", "7")
	tree.Insert("h", "-100")
	tree.Insert("z", "-52")

	item, found := tree.GetItemByKey("g")
	if !found {
		t.Fatal("expected to find key 'g'")
	}
	fmt.Println(tree.Print())
	fmt.Println(item)
}

func TestPersistence(t *testing.T) {
	os.Remove("persist_test.db")
	defer os.Remove("persist_test.db")

	tree, err := btree.NewBTree(2, "persist_test.db")
	if err != nil {
		t.Fatal(err)
	}

	tree.Insert("hello", "world")
	tree.Insert("foo", "bar")
	tree.Insert("abc", "123")
	tree.Close()

	tree2, err := btree.NewBTree(2, "persist_test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer tree2.Close()

	item, found := tree2.GetItemByKey("foo")
	if !found {
		t.Fatal("expected 'foo' to not be deleted")
	}
	if item.Value != "bar" {
		t.Fatalf("expected value 'bar', got '%s'", item.Value)
	}

	fmt.Println("after launch:", tree2.Print())
}
