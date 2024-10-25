package sorting

import "sort"

// BuiltinSort uses Go's built-in sort.Ints to sort the array.
func BuiltinSort(arr []int) []int {
	sort.Ints(arr)
	return arr
}
