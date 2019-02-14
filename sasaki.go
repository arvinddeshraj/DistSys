// Multithreaded Golang Implementation of Sasaki's Time Optimal N-1 Round Algorithm
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
// Each Process contains a left block, a right block as well as
// a area variable that is used for solution selection
type Process struct {
	leftvalue Block
	rightvalue Block
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
	for i := 0; i < n; i++ {
		if order == 1 {
			proc[i].leftvalue.val = n - i
			proc[i].rightvalue.val = n - i
		} else if order == 2 {
			proc[i].leftvalue.val = rand.Intn(12345)
			proc[i].rightvalue.val =proc[i].leftvalue.val
		} else {
			proc[i].leftvalue.val = i
			proc[i].rightvalue.val = i
		}
		proc[i].area = 0
		if i == 0 || i == n -1 {
			proc[i].leftvalue.ismarked = true
			proc[i].rightvalue.ismarked = true
		}
	}
	proc[0].area = -1
}
// utility function to display the current state of each process
func DisplayCurrent(proc []Process) {
	for _, v := range proc {
		if v.leftvalue.ismarked == true {
			fmt.Print(v.leftvalue.val,"*\t")
		} else {
			fmt.Print(v.leftvalue.val,"\t")
		}
		if v.rightvalue.ismarked == true {
			fmt.Print(v.rightvalue.val,"*\t")
		} else {
			fmt.Print(v.rightvalue.val,"\t")
		}
	}
	fmt.Println()
}
// SendAndReceive is a function used to simulate the send and receive operations
// that would occur in a distributed system. It also changes the area variables based 
// on the movement of the marked numbers.
func SendAndReceive(proc []Process, i int) {
	var temp Block
	if proc[i].rightvalue.val > proc[i+1].leftvalue.val {
		if proc[i].rightvalue.ismarked == true && proc[i+1].leftvalue.ismarked == false	{
			proc[i+1].area -= 1
		} else if proc[i].rightvalue.ismarked == false && proc[i+1].leftvalue.ismarked == true {
			proc[i+1].area += 1
		} else if proc[i].rightvalue.ismarked == true && proc[i+1].leftvalue.ismarked == true {
			proc[i+1].area += 1
			proc[i+1].area -= 1
		}
		temp = proc[i].rightvalue
		proc[i].rightvalue = proc[i+1].leftvalue
		proc[i+1].leftvalue = temp
	}
}
// SortNeighbours is used to order the values in the middle processes
// (processes other than the first and last)
func SortNeighbours(proc []Process, i int) {
	var temp Block
	if i == 0 {
		proc[i].leftvalue = proc[i].rightvalue
	}
	if i + 1 == len(proc) - 1 {
		proc[i+1].rightvalue = proc[i+1].leftvalue
	}
	if proc[i].leftvalue.val > proc[i].rightvalue.val {
		temp = proc[i].rightvalue
		proc[i].rightvalue = proc[i].leftvalue
		proc[i].leftvalue = temp
	}
}

// SolutionSelection is a function to select the left or right value for the final solution
// after n-1 rounds.
func SolutionSelection(proc []Process) {
	fmt.Print(proc[0].rightvalue.val, "\t")
	for i := 1; i < len(proc)-1; i++ {
		if proc[i].area == -1 {
			fmt.Print(proc[i].rightvalue.val, "\t")
		} else {
			fmt.Print(proc[i].leftvalue.val, "\t")
		}

	}
	fmt.Print(proc[len(proc)-1].leftvalue.val, "\t")
	fmt.Println()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var wg sync.WaitGroup // wg is used to wait for go routines
	// It ensures that sorting of intermeditate processes' values
	// happens only after all goroutines finish their swaps
	n := rand.Intn(100) + 1
	processes := make([]Process, n)
	InitList(processes, 2)
	DisplayCurrent(processes)
	for i := 0; i < n-1; i++ {
		for j, _ := range processes {
			if j != n -1 {
				wg.Add(1)
				go func(processes []Process, j int) {
					SendAndReceive(processes, j)
					wg.Done()
				}(processes, j)
			}
		}
		wg.Wait()
		for j, _ := range processes {
			if j != len(processes)-1 {
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
