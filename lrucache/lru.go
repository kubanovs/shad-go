//go:build !solution

package lrucache

import "container/list"

type (
	ElementWrapped struct {
		Key int
		Val int
	}

	LRUCache struct {
		lst    *list.List
		mp     map[int]*list.Element
		mxSize int
	}
)

func New(cap int) Cache {
	return &LRUCache{
		lst: list.New(),
		mp:  make(map[int]*list.Element, cap), mxSize: cap,
	}
}

// Get returns value associated with the key.
// The second value is a bool that is true if the key exists in the cache,
// and false if not.
func (cache *LRUCache) Get(key int) (int, bool) {
	el, ok := cache.mp[key]
	if !ok {
		return 0, false
	}
	cache.lst.MoveToFront(el)
	return el.Value.(*ElementWrapped).Val, true
}

// Set updates value associated with the key.
//
// If there is no key in the cache new (key, value) pair is created.
func (cache *LRUCache) Set(key int, value int) {
	if cache.mxSize < 1 {
		return
	}

	el, ok := cache.mp[key]

	if !ok {
		newEl := &ElementWrapped{Key: key, Val: value}
		if cache.lst.Len() == cache.mxSize {
			delete(cache.mp, cache.lst.Back().Value.(*ElementWrapped).Key)
			cache.lst.Remove(cache.lst.Back())
		}
		cache.mp[key] = cache.lst.PushFront(newEl)
	} else {
		el.Value.(*ElementWrapped).Val = value
		cache.lst.MoveToFront(el)
	}
}

// Range calls function f on all elements of the cache
// in increasing access time order.
//
// Stops earlier if f returns false.
func (cache *LRUCache) Range(f func(key, value int) bool) {
	for e := cache.lst.Back(); e != nil; e = e.Prev() {
		wrapped := e.Value.(*ElementWrapped)
		if !f(wrapped.Key, wrapped.Val) {
			break
		}
	}
}

// Clear removes all keys and values from the cache.
func (cache *LRUCache) Clear() {
	clear(cache.mp)
	cache.lst.Init()
}
