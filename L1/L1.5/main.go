package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Error: ", err)
	}

	counterChan := make(chan int)

	go func() {
		for i := 1; i <= 100; i++ {
			counterChan <- i
			time.Sleep(time.Second)
		}
		close(counterChan)
	}()

	timout := time.After(time.Duration(n) * time.Second)

	for {
		select {
		case val, ok := <-counterChan:
			if !ok {
				break
			}
			fmt.Println("Прием данных: ", val)
		case stopTime := <-timout:
			fmt.Println("Таймаут! Вышло время - ", stopTime)
			return
		}
	}
}
