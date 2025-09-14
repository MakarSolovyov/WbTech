package main

import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	array := [5]int{2, 4, 6, 8, 10}
	wg.Add(len(array))

	doStuff := func(number int) {
		result := number * number
		fmt.Printf("Квадрат числа %d равен %d\n", number, result)
		wg.Done()
	}

	for _, value := range array {

		go doStuff(value)
	}

	wg.Wait()
	fmt.Println("Возведение чисел в степень завершено")
}
