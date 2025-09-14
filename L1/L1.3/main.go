package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Ошибка: неверное количество аргументов, должно быть два")
		return
	}

	workerNumber, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Ошибка: невозможно конвертировать введенное значение в int")
		return
	}
	if workerNumber <= 0 {
		fmt.Println("Ошибка: введенное значение должно быть больше 0")
		return
	}

	dataChannel := make(chan string)
	var wg sync.WaitGroup

	// Создание воркеров для принятия данных
	for i := 0; i < workerNumber; i++ {
		wg.Add(1)
		go worker(i, dataChannel, &wg)
	}

	// Главная горутина: постоянное создание данные
	counter := 0
	for {
		counter++

		data := fmt.Sprintf("Данные %d", counter)
		dataChannel <- data

		time.Sleep(500 * time.Millisecond)
	}

	// Недостижимый код, т.к. идет постоянное создание данных в бесконечном цикле
	close(dataChannel)
	wg.Wait()
}

func worker(id int, dataChannel <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Постоянное считывание данных
	for data := range dataChannel {
		fmt.Printf("Воркер %d: %s\n", id, data)
	}
}
