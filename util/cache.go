package util

import (
	"container/list"
	"sync"
)

type keyValue struct {
	key   string
	value interface{}
}

type LruCache struct {
	itemList *list.List
	itemMap  map[string]*list.Element
	maxSize  int
	lock     sync.Mutex
}

func (cache *LruCache) Get(key string) interface{} {
	elem, ok := cache.itemMap[key]
	if !ok {
		return nil
	}
	cache.itemList.MoveToFront(elem)
	kv := elem.Value.(*keyValue)
	return kv.value
}

func (cache *LruCache) Put(key string, val interface{}, timeout int64) error {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	elem, ok := cache.itemMap[key]

	if ok {
		cache.itemList.MoveToFront(elem) //将元素elem移动到itemList的首部，如果elem不属于itemList，则itemList不改变。
		kv := elem.Value.(*keyValue)
		kv.value = val
	} else {
		elem := cache.itemList.PushFront(&keyValue{key: key, value: val}) //在itemList的首部插入值为v的元素，并返回该元素。
		cache.itemMap[key] = elem
		//超出就要删除
		if cache.itemList.Len() > cache.maxSize {
			delElem := cache.itemList.Back() //获取最后一个
			kv := delElem.Value.(*keyValue)
			cache.itemList.Remove(delElem)
			delete(cache.itemMap, kv.key)
		}
	}

	return nil
}

func (cache *LruCache) Delete(key string) error {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	elem, ok := cache.itemMap[key]
	if ok {
		cache.itemList.Remove(elem)
		delete(cache.itemMap, key)
	}

	return nil
}

func (cache *LruCache) IsExist(key string) bool {
	if _, ok := cache.itemMap[key]; ok {
		return true
	}
	return false
}

func (cache *LruCache) ClearAll() error {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	for k, e := range cache.itemMap {
		cache.itemList.Remove(e)
		delete(cache.itemMap, k)
	}

	return nil
}

func (cache *LruCache) Len() int {
	return cache.itemList.Len()
}

func NewLruCache(maxSize int) *LruCache {
	return &LruCache{
		itemList: list.New(),
		itemMap:  make(map[string]*list.Element),
		maxSize:  maxSize,
	}
}
