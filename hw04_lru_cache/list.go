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
	len   int
	first *ListItem
	last  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int { return l.len }

func (l *list) Front() *ListItem {
	if l.len == 0 {
		return nil
	}
	return l.first
}

func (l *list) Back() *ListItem {
	if l.len == 0 {
		return nil
	}
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v, Next: l.first, Prev: nil}
	if l.len == 0 {
		l.first = newItem
		l.last = newItem
	} else {
		l.first.Prev = newItem
		l.first = newItem
	}
	l.len++
	return l.first
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v, Next: nil, Prev: l.last}
	if l.len == 0 {
		l.first = newItem
		l.last = newItem
	} else {
		l.last.Next = newItem
		l.last = newItem
	}
	l.len++
	return l.last
}

func (l *list) Remove(i *ListItem) {
	switch i {
	case l.first:
		l.first = l.first.Next
		l.first.Prev = nil
	case l.last:
		l.last = l.last.Prev
		l.last.Next = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.first {
		return
	} else if i == l.last {
		l.last = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	i.Next = l.first
	l.first.Prev = i
	l.first = i
	i.Prev = nil
}
