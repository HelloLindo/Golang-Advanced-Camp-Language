package Golang_Advanced_Camp_Language

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type IntList struct {
	head   *Node
	length int64
}

type Node struct {
	value  int
	next   *Node
	mu     sync.RWMutex
	marked uint32
}

func newIntNode(value int) *Node {
	return &Node{value: value}
}

func NewInt() *IntList {
	return &IntList{head: newIntNode(0)}
}

/**
Insert a new Node x into the list.
*/
func (l *IntList) Insert(value int) bool {
	headP := unsafe.Pointer(l.head)
	a := (*Node)(atomic.LoadPointer(&headP))
	bP := unsafe.Pointer(&a.next)
	b := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(bP)))
	for {
		// Find the b's position.
		for b != nil && b.value < value {
			a = b
			// Using atomic operation to replace lock.
			b = (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.next))))
		}
		// Check if the node is exist.
		if b != nil && b.value == value {
			return false
		}
		// Lock a so that we can add a new node x.
		a.mu.Lock()
		if a.next != b {
			// Node a was changed.
			a.mu.Unlock()
			a = (*Node)(atomic.LoadPointer(&headP))
			bP = unsafe.Pointer(&a.next)
			b = (*Node)(atomic.LoadPointer((*unsafe.Pointer)(bP)))
			continue
		}
		break
	}
	// Link the new node x between a and b.
	x := newIntNode(value)
	x.next = b
	aNextP := unsafe.Pointer(&a.next)
	atomic.StorePointer((*unsafe.Pointer)(aNextP), unsafe.Pointer(x))
	l.length++
	a.mu.Unlock()

	return true
}

/**
Delete the Node with specific value from the list.
*/
func (l *IntList) Delete(value int) bool {
	headP := unsafe.Pointer(l.head)
	a := (*Node)(atomic.LoadPointer(&headP))
	bP := unsafe.Pointer(&a.next)
	b := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(bP)))
	for {
		// Find the b's position.
		for b != nil && b.value < value {
			a = b
			b = (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.next))))
		}
		// Check if b is not exists
		if b == nil || b.value != value {
			return false
		}
		// b exists.
		b.mu.Lock()
		if b.marked == uint32(1) {
			b.mu.Unlock()
			a = (*Node)(atomic.LoadPointer(&headP))
			bP = unsafe.Pointer(&a.next)
			b = (*Node)(atomic.LoadPointer((*unsafe.Pointer)(bP)))
			continue
		}
		// a is available.
		a.mu.Lock()
		if a.next != b || a.marked == uint32(1) {
			a.mu.Unlock()
			b.mu.Unlock()
			a = (*Node)(atomic.LoadPointer(&headP))
			bP = unsafe.Pointer(&a.next)
			b = (*Node)(atomic.LoadPointer((*unsafe.Pointer)(bP)))
			continue
		}
		break
	}
	// Remove b.
	atomic.StoreUint32(&b.marked, uint32(1))
	aNextP := unsafe.Pointer(&a.next)
	atomic.StorePointer((*unsafe.Pointer)(aNextP), unsafe.Pointer(b.next))
	l.length--
	a.mu.Unlock()
	b.mu.Unlock()
	return true
}

/**
Return true if the list contains a node with specific value.
*/
func (l *IntList) Contains(value int) bool {
	headP := unsafe.Pointer(l.head)
	a := (*Node)(atomic.LoadPointer(&headP))
	xP := unsafe.Pointer(a.next)
	x := (*Node)(atomic.LoadPointer(&xP))
	// Find the node with that value.
	for x != nil && x.value < value {
		a = x
		x = (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.next))))
	}
	if x == nil {
		return false
	}
	// Check if the value is valid.
	return x.value == value && x.marked == uint32(0)
}

/**
Apply function f to all nodes in the list.
*/
func (l *IntList) Range(f func(value int) bool) {
	headP := unsafe.Pointer(l.head)
	a := (*Node)(atomic.LoadPointer(&headP))
	xP := unsafe.Pointer(a.next)
	x := (*Node)(atomic.LoadPointer(&xP))
	for x != nil {
		if !f(x.value) {
			break
		}
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&x)), unsafe.Pointer(x.next))
	}
}

/**
Return the length of the list.
*/
func (l *IntList) Len() int {
	return int(l.length)
}
