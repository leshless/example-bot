package syncmap

import (
	"sync"

	om "github.com/wk8/go-ordered-map"
)

// this library implements custom lock-free ordered map type (sync.Map is not very convinient to use)

type SyncMap[K comparable, V any] struct {
	data  *om.OrderedMap
	mutex sync.Mutex
}

func (m *SyncMap[K, V]) Init(){
	m.data = om.New()
}

func (m *SyncMap[K, V]) Get(key K) (V, bool){
	m.mutex.Lock()
	value, ok := m.data.Get(key)
	m.mutex.Unlock()
	if ok{
		return value.(V), ok
	}else{
		var nilvalue V
		return nilvalue, ok
	}
}

func (m *SyncMap[K, V]) First() (V, bool){
	m.mutex.Lock()
	pair := m.data.Oldest()
	m.mutex.Unlock()
	if pair != nil{
		return pair.Value.(V), true
	}else{
		var nilvalue V
		return nilvalue, false
	}
}

func (m *SyncMap[K, V]) GetAll() []V{
	list := []V{}

	m.mutex.Lock()
	for pair := m.data.Oldest(); pair != nil; pair = pair.Next() {
		list = append(list, pair.Value.(V))
	}
	m.mutex.Unlock()

	return list
}

func (m *SyncMap[K, V]) GetKeys() []K{
	list := []K{}

	m.mutex.Lock()
	for pair := m.data.Oldest(); pair != nil; pair = pair.Next() {
		list = append(list, pair.Key.(K))
	}
	m.mutex.Unlock()

	return list
}

func (m *SyncMap[K, V]) Set(key K, value V){
	m.mutex.Lock()
	m.data.Set(key, value)
	m.mutex.Unlock()
}

func (m *SyncMap[K, V]) Delete(key K){
	m.mutex.Lock()
	m.data.Delete(key)
	m.mutex.Unlock()
}

func (m *SyncMap[K, V]) Clear(){
	m.mutex.Lock()
	m.data = om.New()
	m.mutex.Unlock()
}

func (m *SyncMap[K, V]) Len() int{
	m.mutex.Lock()
	l := m.data.Len()
	m.mutex.Unlock()
	return l
}
