package slice

import (
	"reflect"
	"testing"
)

func TestReverse(t *testing.T) {
	input := interface{}([]string{"this", "is", "a", "test"})
	expected := interface{}([]string{"test", "a", "is", "this"})
	Reverse(input)
	if !reflect.DeepEqual(input, expected) {
		t.Fatalf("%t != %t", input, expected)
	}
}

func BenchmarkReverse(b *testing.B) {
	input := interface{}([]string{"this", "is", "a", "test"})
	for i := 0; i < b.N; i++ {
		Reverse(input)
	}
}

func BenchmarkReverseString(b *testing.B) {
	input := []string{"this", "is", "a", "test"}
	for i := 0; i < b.N; i++ {
		ReverseString(input)
	}
}

func TestJoin(t *testing.T) {
	input := []string{"this", "is", "a", "test"}
	got := Join(input, ' ')
	expected := "this is a test"
	if got != expected {
		t.Fatalf("Got `%s` expected `%s`", got, expected)
	}
}
