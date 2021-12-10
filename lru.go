package lru

import (
	"container/list"
	"sync"
)

type KeyValue struct {
	Key   []byte
	Value []byte
}

type Cache struct {
	capacity int
	cache    *list.List
	elements map[string]*list.Element
	mutex    sync.RWMutex
}

func New(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		cache:    list.New(),
		elements: make(map[string]*list.Element),
	}
}

func (l *Cache) Get(key []byte) []byte {
	l.mutex.RLock()
	l.mutex.RUnlock()
	if elem, ok := l.elements[string(key)]; ok {
		value := elem.Value.(*list.Element).Value.(KeyValue).Value
		l.cache.MoveToFront(elem)
		return value
	}
	return nil
}

func (l *Cache) Put(key, val []byte) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if elem, ok := l.elements[string(key)]; ok {
		l.cache.MoveToFront(elem)
		elem.Value.(*list.Element).Value = KeyValue{
			Key:   key,
			Value: val,
		}
	} else {
		if l.cache.Len() == l.capacity {
			index := l.cache.Back().Value.(*list.Element).Value.(KeyValue).Key
			delete(l.elements, string(index))
			l.cache.Remove(l.cache.Back())
		}
	}

	e := &list.Element{Value: KeyValue{
		Key:   key,
		Value: val,
	}}

	ptr := l.cache.PushFront(e)
	l.elements[string(key)] = ptr
}

func (l *Cache) Remove(key []byte) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if elem, ok := l.elements[string(key)]; ok {
		delete(l.elements, string(key))
		l.cache.Remove(elem)
	}
}
