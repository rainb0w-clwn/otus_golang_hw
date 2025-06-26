package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Наличие структуры необходимо для удаления из словаря последнего элемента.
type listValue struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if element, exists := c.items[key]; exists {
		c.queue.MoveToFront(element)
		// если сделать как ниже, то тест будет ругаться на разность типов
		// element.Value.(*listValue).value = value
		if value != element.Value.(listValue).value {
			element.Value = listValue{key, value}
		}
		return true
	}
	c.queue.PushFront(listValue{key, value})
	c.items[key] = c.queue.Front()
	if len(c.items) > c.capacity {
		delete(c.items, c.queue.Back().Value.(listValue).key)
		c.queue.Remove(c.queue.Back())
	}
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if element, exists := c.items[key]; exists {
		c.queue.MoveToFront(element)
		return element.Value.(listValue).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	// Под очисткой понимается обнуление НО не затирание списка и словаря
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
