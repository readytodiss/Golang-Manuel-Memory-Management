package allocator

import (
	"fmt"
	"reflect"
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

func (a *Allocator) AllocateValue(value interface{}) interface{} {
	ptr := a.Allocate(unsafe.Sizeof(reflect.Indirect(reflect.ValueOf(value)).Interface()))
	rv := reflect.NewAt(reflect.TypeOf(value), ptr)
	rv.Elem().Set(reflect.ValueOf(value))
	return rv.Interface()
}

func (a *Allocator) AllocateSlice(sliceType reflect.Type, length int) interface{} {
	elemSize := sliceType.Elem().Size()
	size := uintptr(length) * elemSize
	ptr := a.Allocate(size)
	sliceHeader := reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  length,
		Cap:  length,
	}
	return reflect.NewAt(sliceType, unsafe.Pointer(&sliceHeader)).Elem().Interface()
}

func (a *Allocator) String() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return fmt.Sprintf("Allocator{offset: %d, capacity: %d}", a.offset, len(a.memory))
}
