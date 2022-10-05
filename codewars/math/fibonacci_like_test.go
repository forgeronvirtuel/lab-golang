package math

import (
	"reflect"
	"testing"
)

func TestTribonacci(t *testing.T) {
	input_arr := [3]float64{1, 1, 1}
	input_n := 10
	got := Tribonacci(input_arr, input_n)
	expected := []float64{1, 1, 1, 3, 5, 9, 17, 31, 57, 105}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("not equal: %v != %v", got, expected)
	}
}
