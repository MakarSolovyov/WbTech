package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {

	ConditionBoolVar()
	time.Sleep(3 * time.Second)
	ConditionBoolAtomic()
	time.Sleep(3 * time.Second)

	BoolChan()
	time.Sleep(3 * time.Second)
	CloseChan()
	time.Sleep(3 * time.Second)

	ContextCancel()
	time.Sleep(3 * time.Second)
	ContextTimeout()
	time.Sleep(3 * time.Second)
	ContextDeadline()
	time.Sleep(3 * time.Second)

	Timer()
	time.Sleep(4 * time.Second)
	PanicIntercept()
	time.Sleep(3 * time.Second)
	Goexit()
	time.Sleep(3 * time.Second)
	go OSNotifyChan()
	go GracefulShutdown()
	time.Sleep(3 * time.Second)
}

// Сигнал о завершнеии горутины через передачу ref bool переменной
func worker(stopper *bool) {
	fmt.Println("ConditionBoolVar started")
	for !*stopper {
		fmt.Println("ConditionBoolVar working...")
		time.Sleep((time.Second))
	}
	fmt.Println("ConditionBoolVar stopped")
}

func ConditionBoolVar() {
	stopper := false

	go worker(&stopper)
	time.Sleep(3 * time.Second)
	stopper = true
}

// Сигнал о завершнеии горутины через передачу atomic bool переменной
func ConditionBoolAtomic() {
	var atomicBool atomic.Bool
	atomicBool.Store(true)

	go func() {
		fmt.Println("ConditionBoolAtomic started")
		defer atomicBool.Store(false)

		for atomicBool.Load() {
			fmt.Println("ConditionBoolAtomic working...")
			time.Sleep(time.Second)
		}
		fmt.Println("ConditionBoolAtomic stopped")
	}()

	time.Sleep(3 * time.Second)
	atomicBool.Store(false)
}

// Сигнал о завершении горутины через закрытие канала
func CloseChan() {
	stopperChan := make(chan int)

	//Worker
	go func() {
		fmt.Println("CloseChan started ^_^")
		for {
			if _, ok := <-stopperChan; !ok {
				fmt.Println("CloseChan stopped")
				break
			} else {
				fmt.Println("CloseChan working...")
				time.Sleep((time.Second))
			}
		}
	}()

	time.Sleep((3 * time.Second))
	close(stopperChan)
}

// Сигнал о завершении горутины через булевый канал
func BoolChan() {
	done := make(chan bool)

	//Worker
	go func() {
		fmt.Println("BoolChan started ^_^")
		for {
			select {
			case signal := <-done:
				fmt.Println("Worker stopped, signal value: ", signal)
				return
			default:
				fmt.Println("BoolChan working...")
				time.Sleep((time.Second))
			}
		}
	}()

	time.Sleep(3 * time.Second)
	done <- true
}

// Завершение горутины через отмену контекста
func ContextCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Worker
	go func() {
		fmt.Println("ContextCancel started")
		for {
			select {
			case signal := <-ctx.Done():
				fmt.Println("Worker stopped, signal value: ", signal)
				return
			default:
				fmt.Println("ContextCancel working...")
				time.Sleep((time.Second))
			}
		}
	}()

	time.Sleep(3 * time.Second)
	cancel()
}

// Завершение горутины через таймаут контекста
func ContextTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go func() {
		fmt.Println("ContextTimeout started")
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Worker stopped, signal value: ", ctx.Err())
				return
			default:
				fmt.Println("ContextTimeout working...")
				time.Sleep((time.Second))
			}
		}
	}()

	time.Sleep(3 * time.Second)
	<-ctx.Done()
}

// Завершение горутины через дедлайн контекста
func ContextDeadline() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	go func() {
		fmt.Println("ContextDeadline started")
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Worker stopped, signal value: ", ctx.Err())
				return
			default:
				fmt.Println("ContextDeadline working...")
				time.Sleep((time.Second))
			}
		}
	}()

	time.Sleep(3 * time.Second)
	<-ctx.Done()
}

// Завершение горутины по таймеру без контекста
func Timer() {
	timer := time.After(3 * time.Second)
	go func() {
		fmt.Println("Timer started")
		for {
			select {
			case <-timer:
				fmt.Println("Worker stopped, timer ran out")
				return
			default:
				fmt.Println("Timer working...")
				time.Sleep((time.Second))
			}
		}
	}()
}

// Перехват паники
func PanicIntercept() {
	go func() {
		fmt.Println("Panic worker started")
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Worker stopped due to AtrificialPanic: ", r)
			}
		}()
		panic("Fake panic")
	}()
}

// runtime.Goexit()
func Goexit() {
	go func() {
		defer fmt.Println("Worker stopped due to Goexit")

		go func() {
			time.Sleep(3 * time.Second)
			runtime.Goexit() // Закрывает текущего worker'a
		}()
		fmt.Println("Goexit working...")
	}()
}

// Перехват сигнала на закрытие программы
func OSNotifyChan() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("SIGINT started")
		for {
			select {
			case sig := <-sigChan:
				fmt.Println("Worker stopped, received signal: ", sig)
				return
			default:
				fmt.Println("SIGINT working...")
				time.Sleep((time.Second))
			}
		}
	}()

	time.Sleep(3 * time.Second)
	sigChan <- syscall.SIGINT
}

// Graceful shutdown
func gracefulWorker(ctx context.Context, wg *sync.WaitGroup, id int) {
	defer wg.Done()
	defer fmt.Println("gracefulWorker", id, "stopped")
	fmt.Println("gracefulWorker", id, "started")

	for {
		select {
		case signal := <-ctx.Done():
			fmt.Println("Worker stopped, signal value: ", signal)
			return
		default:
			fmt.Println("gracefulWorker working...")
			time.Sleep((time.Second))
		}
	}

}

func GracefulShutdown() {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			gracefulWorker(ctx, &wg, i)
		}(i)
	}

	<-sigChan
	fmt.Println("Initiating graceful shutdown...")
	cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Ждем завершения или таймаута
	select {
	case <-done:
		fmt.Println("All workers stopped gracefully")
	case <-time.After(3 * time.Second):
		fmt.Println("Timeout: forcing shutdown")
	}

	fmt.Println("Graceful shutdown completed")
}
