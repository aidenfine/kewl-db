package main

import (
	"encoding/binary"
	"os"
)

const PageSize = 4096

type Pager struct {
	file     *os.File
	nextPage uint32
}

func NewPager(filename string) (*Pager, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}

	nextPage := uint32(info.Size() / PageSize)

	return &Pager{file: f, nextPage: nextPage}, nil
}

func (p *Pager) Allocate() uint32 {
	id := p.nextPage
	p.nextPage++
	return id
}

func (p *Pager) WritePage(pageID uint32, data []byte) error {
	buf := make([]byte, PageSize)
	copy(buf, data)
	offset := int64(pageID) * PageSize
	_, err := p.file.WriteAt(buf, offset)
	return err
}

func (p *Pager) ReadPage(pageID uint32) ([]byte, error) {
	buf := make([]byte, PageSize)
	offset := int64(pageID) * PageSize
	_, err := p.file.ReadAt(buf, offset)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (p *Pager) Close() error {
	return p.file.Close()
}

func serializeNode(n *node) []byte {
	buf := make([]byte, PageSize)
	pos := 0

	if len(n.children) == 0 {
		buf[pos] = 1
	} else {
		buf[pos] = 0
	}
	pos++

	binary.LittleEndian.PutUint16(buf[pos:], uint16(len(n.items)))
	pos += 2

	for _, item := range n.items {
		binary.LittleEndian.PutUint16(buf[pos:], uint16(len(item.Key)))
		pos += 2
		copy(buf[pos:], item.Key)
		pos += len(item.Key)

		binary.LittleEndian.PutUint16(buf[pos:], uint16(len(item.Value)))
		pos += 2
		copy(buf[pos:], item.Value)
		pos += len(item.Value)
	}

	binary.LittleEndian.PutUint16(buf[pos:], uint16(len(n.children)))
	pos += 2

	for _, childID := range n.children {
		binary.LittleEndian.PutUint32(buf[pos:], childID)
		pos += 4
	}

	return buf
}

func deserializeNode(data []byte, pageID uint32) *node {
	n := &node{pageID: pageID}
	pos := 0

	isLeaf := data[pos]
	_ = isLeaf
	pos++

	numItems := int(binary.LittleEndian.Uint16(data[pos:]))
	pos += 2

	n.items = make([]Item, numItems)
	for i := range numItems {
		keyLen := int(binary.LittleEndian.Uint16(data[pos:]))
		pos += 2
		key := string(data[pos : pos+keyLen])
		pos += keyLen

		valLen := int(binary.LittleEndian.Uint16(data[pos:]))
		pos += 2
		val := string(data[pos : pos+valLen])
		pos += valLen

		n.items[i] = Item{Key: key, Value: val}
	}

	numChildren := int(binary.LittleEndian.Uint16(data[pos:]))
	pos += 2

	n.children = make([]uint32, numChildren)
	for i := range numChildren {
		n.children[i] = binary.LittleEndian.Uint32(data[pos:])
		pos += 4
	}

	return n
}
