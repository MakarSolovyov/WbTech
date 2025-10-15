package main

import (
	"fmt"
	"sync"
)

func main() {
	numbers := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}

	numbersChan := make(chan int)
	doublesChan := make(chan int)
	var wg sync.WaitGroup

	wg.Add(1)
	go intakeNumbers(numbers, numbersChan)
	go doubleNumbers(numbersChan, doublesChan)
	go displayNumbers(doublesChan, &wg)

	wg.Wait()
	fmt.Println("Done")
}

func intakeNumbers(numbers []int, intakeChan chan<- int) {
	defer close(intakeChan)
	for _, num := range numbers {
		intakeChan <- num
	}
}

func doubleNumbers(intakeChan <-chan int, displayChan chan<- int) {
	defer close(displayChan)
	for num := range intakeChan {
		displayChan <- num * 2
	}
}

func displayNumbers(displayChan <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for result := range displayChan {
		fmt.Println(result)
	}
}
