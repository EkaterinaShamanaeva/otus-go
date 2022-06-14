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

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	// элемент есть в словаре
	if element, ok := c.items[key]; ok {
		c.queue.MoveToFront(element) // перемещаем вперед по очереди
		item := element.Value.(*cacheItem)
		item.value = value
		return true
	}
	if c.queue.Len() == c.capacity {
		idx := c.queue.Back().Value.(*cacheItem).key
		delete(c.items, idx)
		c.queue.Remove(c.queue.Back())
	}
	item := &cacheItem{
		key:   key,
		value: value,
	}
	element := c.queue.PushFront(item)
	c.items[key] = element
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	var cValue interface{}

	element, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(element)
	item := element.Value.(*cacheItem)
	cValue = item.value
	return cValue, true
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
