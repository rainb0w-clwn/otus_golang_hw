package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func listToSlice[T any](l List) []T {
	elems := make([]T, 0, l.Len())
	for i := l.Front(); i != nil; i = i.Next {
		elems = append(elems, i.Value.(T))
	}
	return elems
}

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("complex_custom", func(t *testing.T) {
		l := NewList()

		l.PushBack(1) // [1]
		require.Equal(t, l.Front(), l.Back())
		l.Remove(l.Front())
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())

		l.PushBack(1)  // [1]
		l.PushFront(2) // [2, 1]
		l.PushFront(3) // [3, 2, 1]
		l.PushBack(4)  // [3, 2, 1, 4]
		l.PushBack(5)  // [3, 2, 1, 5]
		require.Equal(t, []int{3, 2, 1, 4, 5}, listToSlice[int](l))

		l.Remove(l.Front())
		require.Equal(t, []int{2, 1, 4, 5}, listToSlice[int](l))
		l.Remove(l.Back())
		require.Equal(t, []int{2, 1, 4}, listToSlice[int](l))

		middle := l.Front().Next // 1
		l.MoveToFront(middle)
		require.Equal(t, []int{1, 2, 4}, listToSlice[int](l))

		l.Remove(l.Front()) // [2, 4]
		require.Equal(t, []int{2, 4}, listToSlice[int](l))
		l.MoveToFront(l.Back()) // [4, 2]
		l.MoveToFront(l.Back()) // [2, 4]
		require.Equal(t, []int{2, 4}, listToSlice[int](l))
	})
}
