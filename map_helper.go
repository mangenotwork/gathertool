/*
	Description : map相关的方法
	Author : ManGe
	Mail : 2912882908@qq.com
*/

package gathertool

import (
	"bytes"
	"fmt"
	"reflect"
	"sync"
)

// orderMap 固定顺序map
type orderMap[K, V comparable] struct {
	mux     sync.Mutex // TODO 使用读写锁
	data    map[K]V
	keyList []K
	size    int
}

// OrderMap ues: OrderMap[K, V]()
func OrderMap[K, V comparable]() *orderMap[K, V] {
	obj := &orderMap[K, V]{
		mux:     sync.Mutex{},
		data:    make(map[K]V),
		keyList: make([]K, 0),
		size:    0,
	}
	return obj
}

func (m *orderMap[K, V]) Add(key K, value V) *orderMap[K, V] {
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

func (m *orderMap[K, V]) Get(key K) V {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.data[key]
}

func (m *orderMap[K, V]) Del(key K) *orderMap[K, V] {
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

func (m *orderMap[K, V]) Len() int {
	return m.size
}

func (m *orderMap[K, V]) KeyList() []K {
	return m.keyList
}

func (m *orderMap[K, V]) AddMap(data map[K]V) *orderMap[K, V] {
	for k, v := range data {
		m.Add(k, v)
	}
	return m
}

func (m *orderMap[K, V]) Range(f func(k K, v V)) *orderMap[K, V] {
	for i := 0; i < m.size; i++ {
		f(m.keyList[i], m.data[m.keyList[i]])
	}
	return m
}

// RangeAt Range 遍历map含顺序id
func (m *orderMap[K, V]) RangeAt(f func(id int, k K, v V)) *orderMap[K, V] {
	for i := 0; i < m.size; i++ {
		f(i, m.keyList[i], m.data[m.keyList[i]])
	}
	return m
}

// CheckValue 查看map是否存在指定的值
func (m *orderMap[K, V]) CheckValue(value V) bool {
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
func (m *orderMap[K, V]) Reverse() *orderMap[K, V] {
	for i, j := 0, len(m.keyList)-1; i < j; i, j = i+1, j-1 {
		m.keyList[i], m.keyList[j] = m.keyList[j], m.keyList[i]
	}
	return m
}

// TODO 转json

// TODO debug打印

// TODO 插入值指定位置

// TODO 值移动指定位置操作

// TODO 取指定位置的值

// TODO 首位读取并移除

// TODO 末尾读取并移除

// TODO 洗牌

// MysqlNewTable 给mysql提供创建新的固定map顺序为字段的表
func (m *orderMap[K, V]) MysqlNewTable(db Mysql, table string) error {
	var (
		createSql bytes.Buffer
		line      = m.Len()
	)
	if table == "" {
		return TABLE_IS_NULL
	}
	if line < 1 {
		return fmt.Errorf("fiedls len is 0")
	}
	if db.DB == nil {
		_ = db.Conn()
	}
	createSql.WriteString("CREATE TABLE ")
	createSql.WriteString(table)
	createSql.WriteString(" ( temp_id int(11) NOT NULL AUTO_INCREMENT, ")
	m.Range(func(k K, v V) {
		createSql.WriteString(Any2String(k))
		createSql.WriteString(" ")
		createSql.WriteString(dataType2Mysql(v))
		createSql.WriteString(", ")
	})
	createSql.WriteString("PRIMARY KEY (temp_id) ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	_, err := db.DB.Exec(createSql.String())
	if db.log {
		Info("[Sql] Exec : " + createSql.String())
		if err != nil {
			Error("[Sql] Error : " + err.Error())
		}
	}
	if db.allTN == nil {
		_ = db.allTableName()
	}
	db.allTN.add(table)
	return nil
}

// MysqlInsert  给mysql提供固定顺序map写入
func (m *orderMap[K, V]) MysqlInsert(db Mysql, table string) error {
	var (
		line         = m.Len()
		fieldDataMap = make(map[string]interface{})
	)
	if table == "" {
		return TABLE_IS_NULL
	}
	if line < 1 {
		return fmt.Errorf("fiedls len is 0")
	}
	if db.DB == nil {
		_ = db.Conn()
	}
	if db.allTN == nil {
		_ = db.allTableName()
	}
	if !db.allTN.isHave(table) {
		err := m.MysqlNewTable(db, table)
		if err != nil {
			return err
		}
	}
	m.Range(func(k K, v V) {
		fieldDataMap[Any2String(k)] = v
	})
	return db.insert(table, fieldDataMap)
}

// TODO 支持泛型
type Set map[any]struct{}

func NewSet() Set {
	return make(Set)
}

func (s Set) Has(key any) bool {
	_, ok := s[key]
	return ok
}

func (s Set) Add(key any) {
	s[key] = struct{}{}
}

func (s Set) Delete(key any) {
	delete(s, key)
}

// TODO 支持泛型
type Stack struct {
	data map[any]any
}

func NewStack() *Stack {
	s := new(Stack)
	s.data = make(map[any]any)
	return s
}

func (s *Stack) Push(data interface{}) {
	s.data[len(s.data)] = data
}

func (s *Stack) Pop() {
	delete(s.data, len(s.data)-1)
}

func (s *Stack) String() string {
	info := ""
	for i := 0; i < len(s.data); i++ {
		info = info + "[" + StringValue(s.data[i]) + "]"
	}
	return info
}

func MapCopy(data map[any]any) (copy map[any]any) {
	copy = make(map[any]any, len(data))
	for k, v := range data {
		copy[k] = v
	}
	return
}

func MapMergeCopy(src ...map[any]any) (copy map[any]any) {
	copy = make(map[any]any)
	for _, m := range src {
		for k, v := range m {
			copy[k] = v
		}
	}
	return
}

// Map2Slice Eg: {"K1": "v1", "K2": "v2"} => ["K1", "v1", "K2", "v2"]
func Map2Slice(data any) []any {
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
		array := make([]any, 0)
		for _, key := range reflectValue.MapKeys() {
			array = append(array, key.Interface())
			array = append(array, reflectValue.MapIndex(key).Interface())
		}
		return array
	}
	return nil
}

// Slice2Map ["K1", "v1", "K2", "v2"] => {"K1": "v1", "K2": "v2"}
// ["K1", "v1", "K2"]       => nil
func Slice2Map(slice any) map[any]any {
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
		data := make(map[any]any)
		for i := 0; i < reflectValue.Len(); i += 2 {
			data[reflectValue.Index(i).Interface()] = reflectValue.Index(i + 1).Interface()
		}
		return data
	}
	return nil
}
