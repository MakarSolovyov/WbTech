package cache

import (
	"container/list"
	"l0/internal/model"
	"log"
)

type CacheEntry struct {
	Key   string
	Value interface{}
}

type LRUCache struct {
	Capacity int
	Items    map[string]*list.Element
	list     *list.List
}

func CreateLRU(capacity int) *LRUCache {
	return &LRUCache{
		Capacity: capacity,
		Items:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func FillRecentData(cache *LRUCache) {

	recentOrders, err := model.GetOrders(cache.Capacity)
	if err != nil {
		log.Println(err)
		return
	}

	for _, order := range recentOrders {
		cache.Add(order.Order_uid, order)
	}
}

func (lru *LRUCache) Add(key string, value interface{}) {

	if existingElement, exists := lru.Items[key]; exists {
		lru.list.MoveToFront(existingElement)
		existingElement.Value.(*CacheEntry).Value = value
		return
	}

	entry := &CacheEntry{Key: key, Value: value}
	element := lru.list.PushFront(entry)
	lru.Items[key] = element

	if len(lru.Items) > lru.Capacity {
		lru.removeOldest()
	}
}

func (lru *LRUCache) removeOldest() {

	oldestElement := lru.list.Back()
	if oldestElement != nil {
		delete(lru.Items, oldestElement.Value.(*CacheEntry).Key)
		lru.list.Remove(oldestElement)
	}
}

func (lru *LRUCache) Get(key string) (interface{}, bool) {

	if element, exists := lru.Items[key]; exists {
		lru.list.MoveToFront(element)
		return element.Value.(*CacheEntry).Value, true
	}
	return nil, false
}
