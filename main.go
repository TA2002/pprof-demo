package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"pprof-demo/sorting"
)

func generateRandomArray(size int) []int {
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = rand.Intn(size)
	}
	return arr
}

func generateSortedArray(size int) []int {
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = i
	}
	return arr
}

// Function to perform and time a sorting operation
func performSort(sortName string, sortFunc func([]int) []int, arr []int, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	fmt.Printf("Sorting with %s...\n", sortName)
	sortedArr := sortFunc(arr) // Capture returned sorted array
	fmt.Printf("%s time taken: %v\n", sortName, time.Since(start))

	// Optional: Log sorted array (or handle as needed)
	_ = sortedArr // Prevents unused variable warning if not using sortedArr
}

func main() {
	// Set memory profiling rate - 1 means track all allocations
	runtime.MemProfileRate = 1

	// Open files to save CPU and memory profiles
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("Could not create CPU profile:", err)
		return
	}
	defer cpuFile.Close()

	memFile, err := os.Create("mem.prof")
	if err != nil {
		fmt.Println("Could not create memory profile:", err)
		return
	}
	defer memFile.Close()

	// Start CPU profiling
	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		fmt.Println("Could not start CPU profile:", err)
		return
	}
	defer pprof.StopCPUProfile()

	runtime.SetBlockProfileRate(1)

	// Generate large dataset
	// Test case 1: Random array
	// Comment next line and uncomment the next text case to test with a sorted array
	arr := generateRandomArray(100000)
	
	// Test case 2: Sorted array (worst case scenario for QuickSort)
	// Uncomment next line to test with a sorted array. Don't forget to comment the previous line
	// arr := generateSortedArray(100000)

	// Prepare different arrays for each sorting function
	sortTasks := []struct {
		name     string
		sortFunc func([]int) []int
		arr      []int
	}{
		{"Bubble Sort", sorting.BubbleSort, make([]int, len(arr))},
		{"Count Sort", sorting.CountSort, make([]int, len(arr))},
		{"Quick Sort", sorting.QuickSortStart, make([]int, len(arr))},
		{"Builtin Sort", sorting.BuiltinSort, make([]int, len(arr))},
		{"Merge Sort", sorting.MergeSort, make([]int, len(arr))},
	}

	for _, task := range sortTasks {
		copy(task.arr, arr)
	}

	// Concurrently execute each sorting function
	var wg sync.WaitGroup
	for _, task := range sortTasks {
		wg.Add(1)
		go performSort(task.name, task.sortFunc, task.arr, &wg)
	}

	wg.Wait()
	pprof.WriteHeapProfile(memFile)
}
