package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrNoWorkerProvided    = errors.New("no workers provided")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Обработка граничных случаев
	if len(tasks) == 0 {
		return nil
	}
	if n <= 0 {
		return ErrNoWorkerProvided
	}
	var errorCount int64
	taskChannel := make(chan Task)
	wg := sync.WaitGroup{}
	checkError := m > 0

	// Создаем горутины-воркеры
	for i := 0; i < min(n, len(tasks)); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChannel {
				if err := task(); err != nil {
					atomic.AddInt64(&errorCount, 1)
				}
			}
		}()
	}

	isError := func() bool {
		return checkError && atomic.LoadInt64(&errorCount) >= int64(m)
	}
	// Отправляем задачи на выполнение
	for _, task := range tasks {
		if isError() {
			break
		}
		taskChannel <- task
	}

	close(taskChannel)

	wg.Wait()

	if isError() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
