package strings

import "strings"

func ComputeWeight(s string) int {
	var sum int
	for _, r := range s {
		sum += int(r - '0')
	}
	return sum
}

func OrderWeight(strng string) string {
	// Part 1. split the string
	previ := -1
	var strs []string
	for i := 0; i < len(strng); i++ {
		if rune(strng[i]) == ' ' {
			strs = append(strs, strng[previ+1:i])
			previ = i
		}
	}
	strs = append(strs, strng[previ+1:])

	// Part 2. convert in numbers
	var weights = make([]int, len(strs))
	for i, s := range strs {
		var sum int
		for _, r := range s {
			sum += int(r - '0')
		}
		weights[i] = sum
	}

	// Part 3. sort the two arrays
	QuickSortTuple(weights, strs, 0, len(weights)-1)

	// Part 4. search for duplicate values
	//prev := strs[0]
	//prevIdx := 0
	//treatInterval := false
	//for i := 1; i < len(strs); i++ {
	//	if strs[i] == prev {
	//		treatInterval = true
	//		continue
	//	}
	//	if treatInterval {
	//		QuickSort(weights, strs, 0, len(weights)-1)
	//	}
	//	prev = strs[i]
	//	prevIdx = i + 1
	//}

	return strings.Join(strs, " ")
}

func QuickSort(arr []int, low, high int) []int {
	if low < high {
		var p int
		arr, p = partition(arr, low, high)
		arr = QuickSort(arr, low, p-1)
		arr = QuickSort(arr, p+1, high)
	}
	return arr
}

func partition(arr []int, low, high int) ([]int, int) {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return arr, i
}

func QuickSortTuple(arr []int, strs []string, low, high int) ([]int, []string) {
	if low < high {
		var p int
		arr, strs, p = partitionTuple(arr, strs, low, high)
		arr, strs = QuickSortTuple(arr, strs, low, p-1)
		arr, strs = QuickSortTuple(arr, strs, p+1, high)
	}
	return arr, strs
}

func partitionTuple(arr []int, strs []string, low, high int) ([]int, []string, int) {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			strs[i], strs[j] = strs[j], strs[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	strs[i], strs[high] = strs[high], strs[i]
	return arr, strs, i
}
