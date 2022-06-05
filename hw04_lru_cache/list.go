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

type ListItem struct { //узел
	Value interface{} // any type
	Next  *ListItem   // указатель на след значение
	Prev  *ListItem   // указатель на пред значение
}

type list struct {
	len   int // количество эл-тов листа
	first *ListItem
	last  *ListItem // элемент
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
		l.first.Prev = newItem // поменяли ссылку предыдущего первого элемента (до этого был nil)
		l.first = newItem      // поменяли сам первый элемент (значение и указатели)
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
		l.last.Next = newItem // поменяли указатель бывшего последнего элемента (был nil на след значение)
		l.last = newItem      // поменяли сам последний элемент
	}
	l.len++
	return l.last
}

func (l *list) Remove(i *ListItem) {
	if i == l.first {
		l.first = l.first.Next
		l.first.Prev = nil
	} else if i == l.last {
		l.last = l.last.Prev
		l.last.Next = nil
	} else {
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
