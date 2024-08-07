package allocator

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

func (a *Allocator) AllocateValue[T any](value T) *T {
	ptr := a.Allocate(unsafe.Sizeof(value))
	typedPtr := (*T)(ptr)
	*typedPtr = value
	return typedPtr
}

func (a *Allocator) AllocateSlice[T any](length int) []T {
	elemSize := unsafe.Sizeof(*new(T))
	size := uintptr(length) * elemSize
	ptr := a.Allocate(size)
	return unsafe.Slice((*T)(ptr), length)
}
