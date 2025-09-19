package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// Возьмем за основу задание L1.3

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Ошибка: неверное количество аргументов, должно быть два")
		os.Exit(1)
	}

	workerNumber, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Ошибка: невозможно конвертировать введенное значение в int")
		os.Exit(1)
	}
	if workerNumber <= 0 {
		fmt.Println("Ошибка: введенное значение должно быть больше 0")
		os.Exit(1)
	}

	dataChannel := make(chan string)
	done := make(chan struct{})
	var wg sync.WaitGroup

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigchan
		fmt.Println("\nПолучен сигнал завершения, останавливаем воркеры...")
		close(done)
	}()

	// Создание воркеров для принятия данных
	for i := 0; i < workerNumber; i++ {
		wg.Add(1)
		go worker(i, dataChannel, &wg, done)
	}

	// Главная горутина: постоянное создание данные
	counter := 0

	for {
		select {
		case <-done:
			// Получен сигнал завершения
			fmt.Println("Останавливаем запись данных...")
			close(dataChannel)
			wg.Wait()
			fmt.Println("Программа завершена")
			return
		default:

			counter++

			data := fmt.Sprintf("Данные #%d", counter)

			// Отправляем данные с проверкой, не закрыт ли уже канал
			select {
			case dataChannel <- data:
			case <-done:
				continue
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}

func worker(id int, dataChannel <-chan string, wg *sync.WaitGroup, done <-chan struct{}) {
	defer wg.Done()

	// Постоянное считывание данных
	for {
		select {
		case <-done:
			fmt.Printf("Воркер %d: получен сигнал завершения\n", id)
			return
		case data, ok := <-dataChannel:

			if !ok {
				fmt.Printf("Воркер %d: канал закрыт, завершаю работу\n", id)
				return
			}

			fmt.Printf("Воркер %d: %s\n", id, data)
		}
	}
}
