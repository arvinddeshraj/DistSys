// Golang Implementation of Sasaki's Time Optimal N-1 Round Algorithm
// Arvind Deshraj

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Block is a wrapper around the list value to track the marked values
type Block struct {
	val      int
	ismarked bool
}

// Process is the main component of the list.
// Each Process contains a single block as well as
// a area variable that is used for solution selection
type Process struct {
	value Block
	area  int
}

// InitList is a function used to initialise the list
// to one of the following cases:
// 1. Worst Case (Reverse Order)
// 2. Random Initialisation
// 3. Best Case (Sorted Order)
// Also initialises the area variable to zero (except for first element)
func InitList(proc []Process, order int) {
	n := len(proc)
	for i := 0; i < len(proc); i++ {
		if order == 1 {
			proc[i].value.val = n - i
		} else if order == 2 {
			proc[i].value.val = rand.Intn(12345)
		} else {
			proc[i].value.val = i
		}
		proc[i].area = 0
		if i != 0 && i != len(proc)-1 {
			proc[i+1].value.val = proc[i].value.val
			proc[i+1].value.ismarked = false
			i += 1
		} else {
			proc[i].value.ismarked = true
		}
		if i == 0 {
			proc[i].area = -1
		}
	}
}

// utility function to display the current state of each process
func DisplayCurrent(proc []Process) {
	for _, v := range proc {
		if v.value.ismarked == true {
			fmt.Print(v.value.val, "*\t")
		} else {
			fmt.Print(v.value.val, "\t")
		}
	}
	fmt.Println()
}

// SendAndReceive is a function used to simulate the send and receive operations
// that would occur in a distributed system. It also changes the area variables based 
// on the movement of the marked numbers.
func SendAndReceive(proc []Process, i int) {
	var temp Block
	if proc[i].value.val > proc[i+1].value.val {
		if proc[i].value.ismarked == true && proc[i+1].value.ismarked == false {
			proc[i+1].area -= 1
		} else if proc[i].value.ismarked == false && proc[i+1].value.ismarked == true {
			proc[i+1].area += 1
		} else if proc[i].value.ismarked == true && proc[i+1].value.ismarked == true {
			proc[i+1].area += 1
			proc[i+1].area -= 1
		}
		temp = proc[i].value 
		proc[i].value = proc[i+1].value
		proc[i+1].value = temp
	}
}

// SortNeighbours is used to order the values in the middle processes
// (processes other than the first and last)
func SortNeighbours(proc []Process, i int) {
	var temp Block
	if proc[i].value.val > proc[i+1].value.val {
		temp = proc[i].value
		proc[i].value = proc[i+1].value
		proc[i+1].value = temp
	}
}

// SolutionSelection is a function to select the left or right value for the final solution
// after n-1 rounds.
func SolutionSelection(proc []Process) {
	fmt.Print(proc[0].value.val, "\t")
	for i := 1; i < len(proc)-1; i += 2 {
		if proc[i].area == -1 {
			fmt.Print(proc[i+1].value.val, "\t")
		} else {
			fmt.Print(proc[i].value.val, "\t")
		}

	}
	fmt.Print(proc[len(proc)-1].value.val, "\t")
	fmt.Println()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var wg sync.WaitGroup // wg is used to wait for go routines
	// It ensures that sorting of intermeditate processes' values
	// happens only after all goroutines finish their swaps
	n := rand.Intn(100) + 1
	processes := make([]Process, 2*n-2)
	InitList(processes, 2)
	DisplayCurrent(processes)
	for i := 0; i < n-1; i++ {
		for j, _ := range processes {
			if j%2 == 1 {
				continue
			} else {
				wg.Add(1)
				go func(processes []Process, j int) {
					SendAndReceive(processes, j)
					wg.Done()
				}(processes, j)
			}
		}
		wg.Wait()
		for j, _ := range processes {
			if j%2 == 1 && j != len(processes)-1 {
				wg.Add(1)
				go func(processes []Process, j int) {
					SortNeighbours(processes, j)
					wg.Done()
				}(processes, j)
			}
		}
		wg.Wait()
		fmt.Println("The list after Round:", i)
		DisplayCurrent(processes)
	}
	fmt.Println("The final solution is:")
	SolutionSelection(processes)
}
