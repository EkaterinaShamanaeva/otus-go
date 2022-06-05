package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache // Remove me after realization.

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
		c.queue.MoveToFront(element)                                         // перемещаем вперед по очереди
		element.Value.(*cacheItem).value = cacheItem{key: key, value: value} //value // обновляем значение //мб тут словарь cacheItem
		return true
	} else {
		// если количество элементов превышено //нет удаления ?
		if c.queue.Len() == c.capacity {
			//c.queue.Remove(c.queue.Back())
			idx := c.queue.Back().Value.(*cacheItem).key
			delete(c.items, idx)
			c.queue.Remove(c.queue.Back())
			//
		}
		item := &cacheItem{
			key:   key,
			value: value,
		}
		element := c.queue.PushFront(item)
		c.items[key] = element //item.key
		return false
	}
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	element, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(element)
	return element.Value.(*cacheItem).value, true
}

func (c *lruCache) Clear() {

}
