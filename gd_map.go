/*
	Description : 固定顺序map
	Author : ManGe
			2912882908@qq.com
			https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"sync"
)

// GDMapApi 固定顺序 Map 接口
type GDMapApi interface {
	Add(key string, value interface{}) *gDMap
	Get(key string) interface{}
	Del(key string) *gDMap
	Len() int
	KeyList() []string
	AddMap(data map[string]interface{}) *gDMap
	Range(f func(k string, v interface{})) *gDMap
	RangeAt(f func(id int, k string, v interface{})) *gDMap
	CheckValue(value interface{}) bool // 检查是否存在某个值
	Reverse()                          //反序
}

// gDMap 固定顺序map
type gDMap struct {
	mux     sync.Mutex
	data    map[string]interface{}
	keyList []string
	size    int
}

// NewGDMap ues: NewGDMap().Add(k,v)
func NewGDMap() *gDMap {
	return &gDMap{
		data:    make(map[string]interface{}),
		keyList: make([]string, 0),
		size:    0,
	}
}

// Add  添加kv
func (m *gDMap) Add(key string, value interface{}) *gDMap {
	m.mux.Lock()
	defer m.mux.Unlock()
	if _, ok := m.data[key]; ok {
		m.data[key] = value
		return m
	}
	m.keyList = append(m.keyList, key)
	m.size++
	m.data[key] = value
	return m
}

// Get 通过key获取值
func (m *gDMap) Get(key string) interface{} {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.data[key]
}

// Del 删除指定key的值
func (m *gDMap) Del(key string) *gDMap {
	m.mux.Lock()
	defer m.mux.Unlock()
	if _, ok := m.data[key]; ok {
		delete(m.data, key)
		for i := 0; i < m.size; i++ {
			if m.keyList[i] == key {
				m.keyList = append(m.keyList[:i], m.keyList[i+1:]...)
				m.size--
				return m
			}
		}
	}
	return m
}

// Len map的长度
func (m *gDMap) Len() int {
	return m.size
}

// KeyList 打印map所有的key
func (m *gDMap) KeyList() []string {
	return m.keyList
}

// AddMap 写入map
func (m *gDMap) AddMap(data map[string]interface{}) *gDMap {
	for k, v := range data {
		m.Add(k, v)
	}
	return m
}

// Range 遍历map
func (m *gDMap) Range(f func(k string, v interface{})) *gDMap {
	for i := 0; i < m.size; i++ {
		f(m.keyList[i], m.data[m.keyList[i]])
	}
	return m
}

// RangeAt Range 遍历map含顺序id
func (m *gDMap) RangeAt(f func(id int, k string, v interface{})) *gDMap {
	for i := 0; i < m.size; i++ {
		f(i, m.keyList[i], m.data[m.keyList[i]])
	}
	return m
}

// CheckValue 查看map是否存在指定的值
func (m *gDMap) CheckValue(value interface{}) bool {
	m.mux.Lock()
	defer m.mux.Unlock()
	for i := 0; i < m.size; i++ {
		if m.data[m.keyList[i]] == value {
			return true
		}
	}
	return false
}

// Reverse map反序
func (m *gDMap) Reverse() {
	for i, j := 0, len(m.keyList)-1; i < j; i, j = i+1, j-1 {
		m.keyList[i], m.keyList[j] = m.keyList[j], m.keyList[i]
	}
}
