// Golang Implementation of an Alternate N-1 Round Algorithm for distributed sorting on a line network
// Arvind Deshraj 
package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// Process is the main component in the list.
// Each Process contains the value stored in that process,
// a modIndex (index % 3).
// Also contains space for two more values - the values sent to
// the process from the neighbouring processes
type Process struct {
	procVal      int
	rightRecvVal int
	leftRecvVal  int
	modIndex     int
}

// InitList is a function used to initialise the list
// to one of the following cases:
// 1. Worst Case (Reverse Order)
// 2. Random Initialisation
// 3. Best Case (Sorted Order)
func InitList(proc []Process, order int) {
	for i, _ := range proc {
		proc[i].modIndex = i % 3
		if order == 1 {
			proc[i].procVal = len(proc) - i // worst case
		} else if order == 2 {
			proc[i].procVal = rand.Intn(12345) // random numbers
		} else {
			proc[i].procVal = i // best case
		}
		proc[i].leftRecvVal = 0
		proc[i].rightRecvVal = 0
	}
}

// IncrementIndex is a function used to increment the modIndex
// of all the processes at the end of a round.
func IncrementIndex(proc []Process) {
	for i, _ := range proc {
		proc[i].modIndex = (proc[i].modIndex + 2) % 3
	}
}

// Utility function to find min of three numbers
func MinInt(a, b, c int) int {
	if a < b && a < c {
		return a
	} else if b < a && b < c {
		return b
	} else {
		return c
	}
}

// Utility function to find max of three numbers
func MaxInt(a, b, c int) int {
	if a > b && a > c {
		return a
	} else if b > a && b > c {
		return b
	} else {
		return c
	}
}

// AltSort is a function that receives the values from its neighbours
// and then sends the respective max and min values to the respective neghbours
func AltSort(proc []Process, i int) {
	var left_avail bool
	var right_avail bool
	var max_val, min_val int

	if i-1 >= 0 {
		left_avail = true
		proc[i].leftRecvVal = proc[i-1].procVal
	}

	if i+1 < len(proc) {
		right_avail = true
		proc[i].rightRecvVal = proc[i+1].procVal
	}

	if left_avail == true && right_avail == true {
		max_val = MaxInt(proc[i].leftRecvVal, proc[i].rightRecvVal, proc[i].procVal)
		min_val = MinInt(proc[i].leftRecvVal, proc[i].rightRecvVal, proc[i].procVal)
		proc[i].procVal = proc[i].leftRecvVal + proc[i].rightRecvVal + proc[i].procVal - max_val - min_val
		proc[i-1].procVal = min_val
		proc[i+1].procVal = max_val
	}

	if left_avail == true && right_avail == false {
		max_val = MaxInt(proc[i].leftRecvVal, proc[i].procVal, math.MinInt64)
		min_val = MinInt(proc[i].leftRecvVal, proc[i].procVal, math.MaxInt64)
		proc[i-1].procVal = min_val
		proc[i].procVal = max_val
	}

	if left_avail == false && right_avail == true {
		max_val = MaxInt(proc[i].rightRecvVal, proc[i].procVal, math.MinInt64)
		min_val = MinInt(proc[i].rightRecvVal, proc[i].procVal, math.MaxInt64)
		proc[i].procVal = min_val
		proc[i+1].procVal = max_val
	}
}

// utility function to display the current state of each process
func DisplayCurrent(proc []Process) {
	for _, v := range proc {
		fmt.Print(v.procVal, "\t")
	}
	fmt.Println()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var wg sync.WaitGroup
	n := rand.Intn(100) + 1
	processes := make([]Process, n)

	InitList(processes, 1)
	DisplayCurrent(processes)

	for i := 0; i < n-1; i++ {
		for j, _ := range processes {
			if processes[j].modIndex == 0 {
				continue
			} else if processes[j].modIndex == 2 {
				continue
			} else {
				wg.Add(1)
				go func(processes []Process, j int) {
					AltSort(processes, j)
					wg.Done()
				}(processes, j)
			}
		}
		wg.Wait()
		IncrementIndex(processes)
		fmt.Println("The list after round ", i)
		DisplayCurrent(processes)
	}
}
