package cache

type ListNode struct {
	key   int
	value int
	next  *ListNode
	prev  *ListNode
}

type LRUCache struct {
	hashMap map[int]int
	head    *ListNode
	tail    *ListNode
	size    int
}

func NewListNode(key, value int, next, prev *ListNode) *ListNode {
	return &ListNode{
		key:   key,
		value: value,
		next:  next,
		prev:  prev,
	}
}

func NewLRUCache(size int) *LRUCache {
	return &LRUCache{
		hashMap: make(map[int]int),
		head:    nil,
		tail:    nil,
		size:    size,
	}
}
