package Golang_Advanced_Camp_Language

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// The linked list struct.
type IntList struct {
	head   *Node
	length int64
}

// Each node in the linked list.
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

// Get the next node's pointer using atomic.Load.
// Call by this way may prevent concurrent problems.
func (n *Node) nextNodePointer() (val unsafe.Pointer) {
	return atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&n.next)))
}

// Return true if insert a new Node x into the list successfully.
func (l *IntList) Insert(value int) bool {
	a := l.head
	b := (*Node)(a.nextNodePointer())
	for {
		// Step 1. Find the positions of a and b.
		for b != nil && b.value < value {
			a = b
			// Using atomic operation to replace lock.
			b = (*Node)(a.nextNodePointer())
		}
		// Check if the node is exist.
		if b != nil && b.value == value {
			return false
		}
		// Step 2. Lock a so that we can add a new node x.
		a.mu.Lock()
		if a.next != b {
			// Node a was changed.
			a.mu.Unlock()
			a = l.head
			b = (*Node)(a.nextNodePointer())
			continue
		}
		break
	}
	// Step 3. Link the new node x between a and b.
	x := newIntNode(value)
	x.next = b
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&a.next)), unsafe.Pointer(x))

	_ = atomic.AddInt64(&l.length, 1)
	// Step 4. Unlock a.
	a.mu.Unlock()

	return true
}

// Return true if delete the Node with specific value from the list successfully.
func (l *IntList) Delete(value int) bool {
	a := l.head
	b := (*Node)(a.nextNodePointer())
	for {
		// Step 1. Find the positions of a and b.
		for b != nil && b.value < value {
			a = b
			b = (*Node)(a.nextNodePointer())
		}
		// Check if b is not exists
		if b == nil || b.value != value {
			return false
		}
		// Step 2. Check if b has been deleted.
		b.mu.Lock()
		if b.marked == uint32(1) {
			b.mu.Unlock()
			a = l.head
			b = (*Node)(a.nextNodePointer())
			continue
		}
		// Step 3. Check if a is edited or available.
		a.mu.Lock()
		if a.next != b || a.marked == uint32(1) {
			a.mu.Unlock()
			b.mu.Unlock()
			a = l.head
			b = (*Node)(a.nextNodePointer())
			continue
		}
		break
	}

	// Step 4. Remove b.
	atomic.StoreUint32(&b.marked, uint32(1))
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&a.next)), unsafe.Pointer(b.next))
	_ = atomic.AddInt64(&l.length, -1)
	// Step 5. Unlock.
	a.mu.Unlock()
	b.mu.Unlock()
	return true
}

// Return true if the list contains a node with specific value.
func (l *IntList) Contains(value int) bool {
	a := l.head
	x := (*Node)(a.nextNodePointer())
	// Find the node with that value.
	for x != nil && x.value < value {
		a = x
		x = (*Node)(a.nextNodePointer())
	}
	if x == nil {
		return false
	}
	// Check if the value is valid.
	return x.value == value && atomic.LoadUint32(&x.marked) == uint32(0)
}

// Traverse all nodes in the list and stop the traverse if function f returns false.
func (l *IntList) Range(f func(value int) bool) {
	x := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&l.head.next))))
	for x != nil {
		// Check if the node is deleted.
		if atomic.LoadUint32(&x.marked) == uint32(1) {
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&x)), x.nextNodePointer())
			continue
		}
		// Apply to the function f.
		if !f(x.value) {
			break
		}
		// Point to the next node.
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&x)), x.nextNodePointer())
	}
}

// Return the length of the list.
func (l *IntList) Len() int {
	return int(atomic.LoadInt64(&l.length))
}
