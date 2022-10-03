package slice

import (
	"reflect"
	"strings"
)

func Reverse(arr interface{}) {
	length := reflect.ValueOf(arr).Len()
	swapFunc := reflect.Swapper(arr)
	for i := 0; i < length/2; i++ {
		swapFunc(i, length-i-1)
	}
}

func ReverseString(arr []string) {
	length := len(arr)
	for i, j := 0, length-1; i < length/2; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

// Join join a set of string separated by sep.
func Join(strs []string, sep rune) string {
	var b strings.Builder
	for _, s := range strs {
		b.WriteString(s)
		b.WriteRune(sep)
	}
	return b.String()
}
