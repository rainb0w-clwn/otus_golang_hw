package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length int
	first  *ListItem
	back   *ListItem
}

func NewList() List {
	return new(list)
}

func (list *list) Len() int {
	return list.length
}

func (list *list) Front() *ListItem {
	return list.first
}

func (list *list) Back() *ListItem {
	return list.back
}

func (list *list) PushFront(v interface{}) *ListItem {
	list.first = &ListItem{Value: v, Next: list.first}
	if list.length == 0 {
		list.back = list.first
	} else {
		list.first.Next.Prev = list.first
	}
	list.length++
	return list.first
}

func (list *list) PushBack(v interface{}) *ListItem {
	list.back = &ListItem{Value: v, Prev: list.back}
	if list.length == 0 {
		list.first = list.back
	} else {
		list.back.Prev.Next = list.back
	}
	list.length++
	return list.back
}

func (list *list) Remove(i *ListItem) {
	switch {
	case i != list.first:
		i.Prev.Next = i.Next
		if i == list.back {
			list.back = i.Prev
		} else {
			i.Next.Prev = i.Prev
		}
	case list.length == 1:
		list.first, list.back = nil, nil
	default:
		i.Next.Prev = nil
		list.first = i.Next
	}
	list.length--
}

func (list *list) MoveToFront(i *ListItem) {
	if i != list.first {
		i.Prev.Next = i.Next
		if i == list.back {
			list.back = i.Prev
		} else {
			i.Next.Prev = i.Prev
		}
		i.Next = list.first
		list.first.Prev = i
		list.first = i
		i.Prev = nil
	}
}
