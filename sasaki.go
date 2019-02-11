package main

import (
	"fmt"
	"math/rand"
	"time"
	"sync"
)

type Block struct {
	val int
	ismarked bool	
}

type Process struct {
	value Block
	area int
}

func init_list(proc []Process, order int) {
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
		if i != 0 && i != len(proc) -1 {
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

func display_current(proc []Process) {
	for _, v := range proc {
		if v.value.ismarked == true {
			fmt.Print(v.value.val,"*\t")
		} else {
			fmt.Print(v.value.val,"\t")
		}
	}
	fmt.Println()
}

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
		temp = proc[i].value//.value.val
		proc[i].value = proc[i+1].value
		proc[i+1].value = temp
	}
}

func sortNeighbours(proc []Process, i int) {
	var temp Block
	if proc[i].value.val > proc[i+1].value.val {
		temp = proc[i].value
		proc[i].value = proc[i+1].value
		proc[i+1].value =temp
	}
}
func solutionSelection(proc []Process) {
	fmt.Print(proc[0].value.val,"\t")
	for i :=1; i < len(proc) - 1; i+=2 {
		if proc[i].area == -1 {
				fmt.Print(proc[i+1].value.val, "\t")
		} else {
				fmt.Print(proc[i].value.val, "\t")
			}
		
	} 
	fmt.Print(proc[len(proc) - 1].value.val,"\t")
	fmt.Println()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())	
	var wg sync.WaitGroup
	n := rand.Intn(100) + 1;
	processes := make([]Process, 2*n - 2)
	init_list(processes, 2)
	display_current(processes)
	for i := 0; i < n - 1 ; i++ {
		for j,_ := range processes {
			if j % 2 == 1 {
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
		for j,_ := range processes {
			if j % 2  == 1 && j != len(processes) - 1 {
				wg.Add(1)
				go func(processes []Process, j int) {
					sortNeighbours(processes, j)
					wg.Done()					
				}(processes, j)
			}
		}
		wg.Wait()
		fmt.Println("The list after Round:", i)	
		display_current(processes)
	}
	fmt.Println("The final solution is:")
	solutionSelection(processes)
}