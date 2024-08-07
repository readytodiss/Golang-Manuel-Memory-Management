package allocater

import (
	"sync"
	"unsafe"
)

type Allocator struct {
	memory []byte
	offset uintptr
	mu     sync.Mutex
}

func NewAllocator(size int) *Allocator {
	return &Allocator{
		memory: make([]byte, size),
		offset: 0,
	}
}

func (a *Allocator) Allocate(size uintptr) unsafe.Pointer {
	a.mu.Lock()
	defer a.mu.Unlock()
	if int(a.offset+size) > len(a.memory) {
		panic("out of memory")
	}
	ptr := unsafe.Pointer(&a.memory[a.offset])
	a.offset += size
	return ptr
}

func (a *Allocator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.offset = 0
}

func (a *Allocator) AllocateInt(value int) *int {
	ptr := a.Allocate(unsafe.Sizeof(value))
	typedPtr := (*int)(ptr)
	*typedPtr = value
	return typedPtr
}

func (a *Allocator) AllocateIntSlice(length int) []int {
	elemSize := unsafe.Sizeof(int(0))
	size := uintptr(length) * elemSize
	ptr := a.Allocate(size)
	return unsafe.Slice((*int)(ptr), length)
}
