package text

import (
	"github.com/forgeronvirtuel/labgolang/codewars/slice"
	"strings"
)

func Split(txt string, r rune) []string {
	previ := -1
	var strs []string
	for i := 0; i < len(txt); i++ {
		if rune(txt[i]) == r {
			strs = append(strs, txt[previ+1:i])
			previ = i
		}
	}
	strs = append(strs, txt[previ+1:])
	return strs
}

// ReverseWords reverse all words in a text
func ReverseWords(txt string) string {
	splitted := Split(txt, ' ')
	slice.ReverseString(splitted)
	return strings.Join(splitted, " ")
}

// Full implementation
//func ReverseWords(txt string) string {
//	// Part 1. split the string
//	previ := -1
//	var strs []string
//	for i := 0; i < len(txt); i++ {
//		if rune(txt[i]) == ' ' {
//			strs = append(strs, txt[previ+1:i])
//			previ = i
//		}
//	}
//	strs = append(strs, txt[previ+1:len(txt)])
//
//	// Part 2. reverse order
//	length := len(strs)
//	for i, j := 0, length-1; i < length/2; i, j = i+1, j-1 {
//		strs[i], strs[j] = strs[j], strs[i]
//	}
//
//
//	// Part 3. Recreate the string
//	return strings.Join(strs, " ")
//}
