package counter

import (
	"fmt"
	"testing"
)

func BenchmarkConcreteIncrement(b *testing.B) {
	cnt := &Counter{}
	for i := 0; i < b.N; i++ {
		cnt.Increment()
	}
}

func BenchmarkIncrementInterface(b *testing.B) {
	var cnt IncrementalCounter
	cnt = &Counter{}
	for i := 0; i < b.N; i++ {
		cnt.Increment()
	}
}

//BenchmarkConcreteIncrement-32           1000000000               0.4825 ns/op
//BenchmarkIncrementInterface-32          613636964                1.707 ns/op

func TestIncrementOnActionWithConcreteType(t *testing.T) {
	var counter *Counter // = nil

	// will work
	counter = &Counter{}
	doStuffAndIncrement(counter /* = typeDescriptor{ type: *ConcreteType, value: &Counter{} } */)

	// will panic
	counter = nil
	doStuffAndIncrement(counter /* = typeDescriptor{ type: *ConcreteType, value: nil } */)
}

func TestIncrementOnActionWithInterface(t *testing.T) {
	var counter IncrementalCounter // = nil

	// will work
	doStuffAndIncrement(counter /* = nil */)

	// will work
	counter = &Counter{} // = typeDescriptor{ type: *ConcreteType, value: &ConcreteType{} }
	doStuffAndIncrement(counter /* = typeDescriptor{ type: *ConcreteType, value: &ConcreteType{} } */)
}

func TestConvertCounter(t *testing.T) {
	var counter IncrementalCounter // = nil
	counter = &Counter{}

	if timer, ok := counter.(DecrementalCounter); ok {
		timer.Decrement()
	}

	fmt.Printf("%s\n", counter)
}
