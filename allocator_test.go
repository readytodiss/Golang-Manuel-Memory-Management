package allocator

import (
	"reflect"
	"testing"
)

func TestAllocator(t *testing.T) {
	a := NewAllocator(1024)

	// Allocate int
	intVal := 42
	intPtr := a.AllocateValue(intVal).(*int)
	if *intPtr != 42 {
		t.Fatalf("expected 42, got %d", *intPtr)
	}

	// Allocate float64
	floatVal := 3.14
	floatPtr := a.AllocateValue(floatVal).(*float64)
	if *floatPtr != 3.14 {
		t.Fatalf("expected 3.14, got %f", *floatPtr)
	}

	// Allocate int slice
	intSliceType := reflect.TypeOf([]int{})
	intSlice := a.AllocateSlice(intSliceType, 10).([]int)
	for i := range intSlice {
		intSlice[i] = i * 2
	}
	for i := range intSlice {
		if intSlice[i] != i*2 {
			t.Fatalf("expected %d, got %d", i*2, intSlice[i])
		}
	}

	// Allocate float64 slice
	floatSliceType := reflect.TypeOf([]float64{})
	floatSlice := a.AllocateSlice(floatSliceType, 10).([]float64)
	for i := range floatSlice {
		floatSlice[i] = float64(i) * 2.5
	}
	for i := range floatSlice {
		if floatSlice[i] != float64(i)*2.5 {
			t.Fatalf("expected %f, got %f", float64(i)*2.5, floatSlice[i])
		}
	}

	// Reset the allocator and allocate again
	a.Reset()
	intVal2 := 24
	intPtr2 := a.AllocateValue(intVal2).(*int)
	if *intPtr2 != 24 {
		t.Fatalf("expected 24, got %d", *intPtr2)
	}
}
