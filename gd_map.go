/*
	Description : 固定顺序map
	Author : ManGe
	Version : v0.1
	Date : 2021-10-11
*/

package gathertool

import "sync"

// 固定顺序 Map 接口
type GDMapApi interface {
	Add(key string, value interface{}) *gDMap
	Del(key string) *gDMap
	Get(key string) interface{}
	Range(f func(k string, v interface{})) *gDMap
	Len() int
	CheckValue(value interface{}) bool // 检查是否存在某个值
	KeyList() []string
}

// 固定顺序map
type gDMap struct {
	mux sync.RWMutex
	data map[string]interface{}
	keyList []string
	size int
}

// ues  NewGDMap().Add(k,v)
func NewGDMap() GDMapApi {
	return &gDMap{
		data:make(map[string]interface{}),
		keyList: make([]string, 0),
		size: 0,
	}
}

func (m *gDMap) Add(key string, value interface{}) *gDMap {
	m.mux.Lock()
	defer m.mux.Unlock()
	// 存在key 就更新值
	if _,ok := m.data[key]; ok {
		m.data[key] = value
		return m
	}
	m.keyList = append(m.keyList, key)
	m.size ++
	m.data[key] = value
	return m
}

func (m *gDMap) Del(key string) *gDMap {
	m.mux.Lock()
	defer m.mux.Unlock()
	if _,ok := m.data[key]; ok {
		delete(m.data, key)
		for i := 0; i < m.size; i++ {
			if m.keyList[i] == key {
				m.keyList = append(m.keyList[:i], m.keyList[i+1:]...)
				m.size --
				return m
			}
		}
	}
	return m
}

func (m *gDMap) Get(key string) interface{} {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return m.data[key]
}

func (m *gDMap) Range(f func(k string, v interface{})) *gDMap {
	m.mux.RLock()
	defer m.mux.RUnlock()
	for i := 0; i < m.size; i++ {
		f(m.keyList[i], m.data[m.keyList[i]])
	}
	return m
}

func (m *gDMap) Len() int {
	return m.size
}

func (m *gDMap) CheckValue(value interface{}) bool {
	m.mux.RLock()
	defer m.mux.RUnlock()
	for i := 0; i < m.size; i++ {
		if m.data[m.keyList[i]] == value {
			return true
		}
	}
	return false
}

func (m *gDMap) KeyList() []string {
	return m.keyList
}
