/*
	Description : 固定顺序map
	Author : ManGe
	Version : v0.1
	Date : 2021-10-11
*/

package gathertool

import (
	"reflect"
	"sync"
)

// 固定顺序 Map 接口
type GDMapApi interface {
	Add(key string, value interface{}) *gDMap
	Del(key string) *gDMap
	Get(key string) interface{}
	Range(f func(k string, v interface{})) *gDMap
	Len() int
	CheckValue(value interface{}) bool // 检查是否存在某个值
	KeyList() []string
	AddMap(data map[string]interface{}) *gDMap
}

// 固定顺序map
type gDMap struct {
	mux sync.RWMutex
	data map[string]interface{}
	keyList []string
	size int
}

// ues  NewGDMap().Add(k,v)
func NewGDMap() *gDMap {
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

func (m *gDMap) AddMap(data map[string]interface{}) *gDMap {
	for k,v := range data {
		m.Add(k,v)
	}
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


func MapCopy(data map[string]interface{}) (copy map[string]interface{}) {
	copy = make(map[string]interface{}, len(data))
	for k, v := range data {
		copy[k] = v
	}
	return
}

// MapMergeCopy
func MapMergeCopy(src ...map[string]interface{}) (copy map[string]interface{}) {
	copy = make(map[string]interface{})
	for _, m := range src {
		for k, v := range m {
			copy[k] = v
		}
	}
	return
}

// Map2Slice Eg: {"K1": "v1", "K2": "v2"} => ["K1", "v1", "K2", "v2"]
func Map2Slice(data interface{}) []interface{} {
	var (
		reflectValue = reflect.ValueOf(data)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Map:
		array := make([]interface{}, 0)
		for _, key := range reflectValue.MapKeys() {
			array = append(array, key.Interface())
			array = append(array, reflectValue.MapIndex(key).Interface())
		}
		return array
	}
	return nil
}

func SliceCopy(data []interface{}) []interface{} {
	newData := make([]interface{}, len(data))
	copy(newData, data)
	return newData
}

// Slice2Map ["K1", "v1", "K2", "v2"] => {"K1": "v1", "K2": "v2"}
// ["K1", "v1", "K2"]       => nil
func Slice2Map(slice interface{}) map[string]interface{} {
	var (
		reflectValue = reflect.ValueOf(slice)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		length := reflectValue.Len()
		if length%2 != 0 {
			return nil
		}
		data := make(map[string]interface{})
		for i := 0; i < reflectValue.Len(); i += 2 {
			data[Any2String(reflectValue.Index(i).Interface())] = reflectValue.Index(i + 1).Interface()
		}
		return data
	}
	return nil
}




