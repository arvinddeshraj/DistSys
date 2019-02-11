package main

import (
	"fmt"
	"math/rand"
	"math"
	"time"
	"sync"
	// "bufio"
	// "os"
	// "strconv"
	// "strings"
)

type Process struct {
	procVal int
	rightRecvVal int
	leftRecvVal int
	modIndex int
}

func init_list(proc []Process, order int) {
	for i, _ := range proc {
		proc[i].modIndex = i % 3
		if order == 1 { 
			proc[i].procVal = len(proc) - i // worst case
		} else if order == 2{
			proc[i].procVal = rand.Intn(12345) // random numbers 
		} else {
			proc[i].procVal = i // best case
		}
		proc[i].leftRecvVal = 0
		proc[i].rightRecvVal = 0
	}	
}

func increment_index(proc []Process) {
	for i, _ := range proc {
		proc[i].modIndex = (proc[i].modIndex + 2) % 3
	}	 
}
func minInt(a, b, c int) int {
	if a < b && a < c {
		return a
	} else if b < a && b < c {
		return b
	} else {
		return c
	}
}

func maxInt(a, b, c int) int {
	if a > b && a > c {
		return a
	} else if b > a && b > c {
		return b
	} else {
		return c
	}
}

func alt_sort(proc []Process, i int) {
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
		max_val = maxInt(proc[i].leftRecvVal, proc[i].rightRecvVal, proc[i].procVal)
		min_val = minInt(proc[i].leftRecvVal, proc[i].rightRecvVal, proc[i].procVal)
		proc[i].procVal = proc[i].leftRecvVal + proc[i].rightRecvVal + proc[i].procVal - max_val - min_val
		proc[i-1].procVal = min_val
		proc[i+1].procVal = max_val  
	}

	if left_avail == true && right_avail == false{
		max_val = maxInt(proc[i].leftRecvVal, proc[i].procVal, math.MinInt64)
		min_val = minInt(proc[i].leftRecvVal, proc[i].procVal, math.MaxInt64)
		proc[i-1].procVal = min_val
		proc[i].procVal = max_val
	}

	if left_avail == false && right_avail == true {
		max_val = maxInt(proc[i].rightRecvVal, proc[i].procVal, math.MinInt64)
		min_val = minInt(proc[i].rightRecvVal, proc[i].procVal, math.MaxInt64)
		proc[i].procVal = min_val
		proc[i+1].procVal = max_val
	}
}

func display_current(proc []Process) {
	for _, v := range proc {
		fmt.Print(v.procVal,"\t")
	}
	fmt.Println()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())	
	var wg sync.WaitGroup
	n := rand.Intn(100) + 1;
	processes :=  make([]Process, n)	
	for i,_ := range processes {
		processes[i].modIndex = i
		processes[i].leftRecvVal = 0
		processes[i].rightRecvVal = 0
	}

	init_list(processes, 1)
	display_current(processes)

	for i := 0; i < n-1; i++ {
		for j, _ := range processes {
			if processes[j].modIndex == 0 {
				continue
			} else if processes[j].modIndex == 2 {
				continue
			} else {
				wg.Add(1)
				go func(processes []Process, j int) {
					alt_sort(processes, j)
					wg.Done()
				}(processes, j)
			}
		}
		wg.Wait()
		increment_index(processes)
		fmt.Println("The list after round ", i )
		display_current(processes)
	} 
}