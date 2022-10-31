package strings

import (
	"reflect"
	"testing"
)

func TestComputeWeight(t *testing.T) {
	input := "180"
	got := ComputeWeight(input)
	exp := 9

	if exp != got {
		t.Fatalf("(got != exp) %d != %d", got, exp)
	}
}

func TestQuickSort(t *testing.T) {
	input_weight := []int{8, 10, 1, 1200, 4, 9, 400, 300, 500}
	input_string := []string{"8", "10", "1", "1200", "4", "9", "400", "300", "500"}
	got_weight, got_strs := QuickSortTuple(input_weight, input_string, 0, len(input_weight)-1)
	expected_weight := []int{1, 4, 8, 9, 10, 300, 400, 500, 1200}
	expected_string := []string{"1", "4", "8", "9", "10", "300", "400", "500", "1200"}

	if !reflect.DeepEqual(got_weight, expected_weight) {
		t.Fatalf("WEIGHTS not equal:\n\t`%v`\n\t`%v`", got_weight, expected_weight)
	}

	if !reflect.DeepEqual(got_strs, expected_string) {
		t.Fatalf("STRINGS not equal:\n\t`%v`\n\t`%v`", got_strs, expected_string)
	}
}

func TestOrderWeightUniqueValues(t *testing.T) {
	input := "103 123 4444 99 2000"
	got := OrderWeight(input)
	expected := "2000 103 123 4444 99"
	if expected != got {
		t.Fatalf("got %s, expected %s", got, expected)
	}
}

func TestOrderWeight(t *testing.T) {
	input := "2000 10003 1234000 44444444 9999 11 11 22 123"
	got := OrderWeight(input)
	expected := "11 11 2000 10003 22 123 1234000 44444444 9999"
	if expected != got {
		t.Fatalf("got\n\t`%s`\n, expected\n\t`%s`", got, expected)
	}
}
