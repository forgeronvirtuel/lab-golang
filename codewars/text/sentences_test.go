package text

import (
	"reflect"
	"testing"
)

func TestReverseWords(t *testing.T) {
	input := "C'est un petit bateau qui naviguait sur l'eau."
	got := ReverseWords(input)
	expected := "l'eau. sur naviguait qui bateau petit un C'est"

	if expected != got {
		t.Fatalf("Got: \n\t`%s`\n, expected: \n\t`%s`\n", got, expected)
	}
}

func TestSplit(t *testing.T) {
	input := "C'est un petit bateau"
	got := Split(input, ' ')
	expected := []string{"C'est", "un", "petit", "bateau"}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Got: \n\t`%s`\n, expected: \n\t`%s`\n", got, expected)
	}
}
