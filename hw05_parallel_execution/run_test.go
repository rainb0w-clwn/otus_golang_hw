package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestRunAdditional(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("error in last task", func(t *testing.T) {
		tasksCount := 100
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := tasksCount - 1
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
	})

	t.Run("error if len(tasks) == m", func(t *testing.T) {
		tasksCount := 5
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				return err
			})
		}

		workersCount := 1
		maxErrorsCount := tasksCount
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
	})

	t.Run("no tasks", func(t *testing.T) {
		tasks := make([]Task, 0)

		err := Run(tasks, 1, 1)

		require.NoError(t, err)
	})

	t.Run("no workers", func(t *testing.T) {
		err := Run([]Task{func() error { return nil }}, 0, 1)

		require.Truef(t, errors.Is(err, ErrNoWorkerProvided), "actual err - %v", err)
	})

	t.Run("no errors when m <= 0", func(t *testing.T) {
		tasksCount := 1
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return ErrErrorsLimitExceeded
			})
		}

		workersCount := 1
		err := Run(tasks, workersCount, 0)
		require.NoError(t, err)

		err = Run(tasks, workersCount, -1)
		require.NoError(t, err)
	})

	t.Run("n > len(tasks)", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 50
		maxErrorsCount := 1

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
	})

	t.Run("all goroutines freed", func(t *testing.T) {
		goroutinesCountBefore := runtime.NumGoroutine()
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		err := Run(tasks, workersCount, maxErrorsCount)
		goroutinesCountAfter := runtime.NumGoroutine()

		require.NoError(t, err)

		require.Equal(t, goroutinesCountBefore, goroutinesCountAfter, "goroutines not freed")
	})
}

func TestRunAdditionalPlus(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("* tasks without errors", func(t *testing.T) {
		tasksCount := 100
		workersCount := 100
		maxErrorsCount := 1
		var sumTime time.Duration
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep
			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				return nil
			})
		}

		require.Eventually(
			t,
			func() bool {
				err := Run(tasks, workersCount, maxErrorsCount)
				return err == nil
			},
			sumTime/2,
			time.Millisecond,
		)
	})
}
