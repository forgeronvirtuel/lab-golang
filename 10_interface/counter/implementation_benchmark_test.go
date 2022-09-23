package counter

import "testing"

func BenchmarkConcreteIncrement(b *testing.B) {
	cnt := &ConcreteCounter{}
	for i := 0; i < b.N; i++ {
		cnt.Increment()
	}
}

func BenchmarkIncrementInterface(b *testing.B) {
	var cnt Counter
	cnt = &ConcreteCounter{}
	for i := 0; i < b.N; i++ {
		cnt.Increment()
	}
}

//BenchmarkConcreteIncrement-32           1000000000               0.4825 ns/op
//BenchmarkIncrementInterface-32          613636964                1.707 ns/op
