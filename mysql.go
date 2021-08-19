/*
	Description : mysql 相关方法
	Author : ManGe
	Version : v0.3
	Date : 2021-08-19
*/


package gathertool

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	SHOW_TABLES = "SHOW TABLES"
	TABLE_NAME_NULL = "table name is null."
)

// 开放的mysql对象
var MysqlDB = &Mysql{}

type Mysql struct {
	host string
	port int
	user string
	password string
	dataBase string
	maxOpenConn int
	maxIdleConn int
	DB *sql.DB
	log bool
	tableTemp  map[string]*tableDescribe //表结构缓存
	once *sync.Once
	allTN *allTableName
}

// NewMysqlDB 给mysql对象进行连接
func NewMysqlDB(host string, port int, user, password, database string) (err error) {
	MysqlDB, err = NewMysql(host,port, user, password, database)
	if err != nil {
		return
	}
	return MysqlDB.Conn()
}

func NewMysql(host string, port int, user, password, database string) (*Mysql, error) {
	if len(host) < 1 {
		return nil, errors.New("Host is Null.")
	}
	if port < 1 {
		port = 3369
	}
	m := &Mysql{
		host : host,
		port : port,
		user : user,
		password : password,
		dataBase : database,
		log: true,
		maxOpenConn: 10,
		maxIdleConn: 10,
		once: &sync.Once{},
	}
	m.once.Do(func() {
		m.tableTemp = make(map[string]*tableDescribe)
	})
	return m, nil
}

// GetMysqlDBConn 获取mysql 连接
func GetMysqlDBConn() (*Mysql,error) {
	err := MysqlDB.Conn()
	return MysqlDB, err
}

// CloseLog 关闭日志
func (m *Mysql) CloseLog(){
	m.log = false
}

// SetMaxOpenConn
func (m *Mysql) SetMaxOpenConn(number int) {
	m.maxOpenConn = number
}

// SetMaxIdleConn
func (m *Mysql) SetMaxIdleConn(number int) {
	m.maxIdleConn = number
}

// Conn 连接mysql
func (m *Mysql) Conn() (err error){
	m.DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s",
		m.user, m.password, "tcp", m.host, m.port, m.dataBase))
	if err != nil {
		if m.log{
			log.Println("[Sql] Conn Fail : " + err.Error())
		}
		return err
	}
	m.DB.SetConnMaxLifetime(time.Hour)  //最大连接周期，超过时间的连接就close
	if m.maxOpenConn < 1{
		m.maxOpenConn = 10
	}
	if m.maxIdleConn < 1{
		m.maxIdleConn = 5
	}
	m.DB.SetMaxOpenConns(m.maxOpenConn)//设置最大连接数
	m.DB.SetMaxIdleConns(m.maxIdleConn) //设置闲置连接数
	return
}

func (m *Mysql) allTableName() (err error) {
	if m.DB == nil{
		_=m.Conn()
	}
	m.allTN = newAllTableName()
	rows,err := m.DB.Query(SHOW_TABLES)
	if err != nil{
		return
	}

	for rows.Next() {
		var result string
		err = rows.Scan(&result)
		log.Println(err, result)
		m.allTN.add(result)
	}

	_=rows.Close()
	return
}

// IsHaveTable 表是否存在
func (m *Mysql) IsHaveTable(name string) bool {
	if m.allTN == nil {
		_=m.allTableName()
	}
	return m.allTN.isHave(name)
}

// 表信息
type TableInfo struct {
	Field    string
	Type     string
	Null     string
	Key      string
	Default  interface{}
	Extra    string
}

type tableDescribe struct {
	Base map[string]string
}

// 记录当前库的所有表名
type allTableName struct{
	mut *sync.Mutex
	tableName map[string]struct{}
}

func newAllTableName() *allTableName {
	return &allTableName{
		mut: &sync.Mutex{},
		tableName: make(map[string]struct{}),
	}
}

func (a *allTableName) add(name string) *allTableName{
	a.mut.Lock()
	a.tableName[name] = struct{}{}
	a.mut.Unlock()
	return a
}

func (a *allTableName) remove(name string) *allTableName{
	a.mut.Lock()
	delete(a.tableName, name)
	a.mut.Unlock()
	return a
}

func (a *allTableName) isHave(name string) bool {
	a.mut.Lock()
	_, ok := a.tableName[name]
	a.mut.Unlock()
	return ok
}

// Describe 获取表结构
func (m *Mysql) Describe(table string) (*tableDescribe, error){
	if m.DB == nil{
		_=m.Conn()
	}

	if v,ok := m.tableTemp[table]; ok {
		return v, nil
	}

	if table == ""{
		return &tableDescribe{}, errors.New(TABLE_NAME_NULL)
	}

	rows,err := m.DB.Query("DESCRIBE " + table)
	if err != nil{
		return &tableDescribe{}, err
	}

	fieldMap := make(map[string]string,0)
	for rows.Next() {
		result := &TableInfo{}
		err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
		log.Println(err, result)
		fiedlType := "null"
		if strings.Contains(result.Type, "int"){
			fiedlType = "int"
		}
		if strings.Contains(result.Type, "varchar") || strings.Contains(result.Type, "text"){
			fiedlType = "string"
		}
		if strings.Contains(result.Type, "float") || strings.Contains(result.Type, "doble") {
			fiedlType = "float"
		}
		if strings.Contains(result.Type, "blob")  {
			fiedlType = "[]byte"
		}
		if strings.Contains(result.Type, "date") || strings.Contains(result.Type, "time") {
			fiedlType = "time"
		}
		fieldMap[result.Field] = fiedlType
	}
	_=rows.Close()

	td := &tableDescribe{
		Base : fieldMap,
	}

	// 缓存
	m.tableTemp[table] = td
	return td, nil
}

// Select 查询语句 返回 map
func (m *Mysql) Select(sql string) ([]map[string]string, error) {
	if m.DB == nil{
		_=m.Conn()
	}

	rows,err := m.DB.Query(sql)
	if m.log{
		log.Println("[Sql] Exec : " + sql)
		if err != nil{
			log.Println("[Sql] Error : " + err.Error())
		}
	}
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil{
		return nil,err
	}

	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache { //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}

	var list []map[string]string //返回的切片
	for rows.Next() {
		_ = rows.Scan(cache...)
		item := make(map[string]string)
		for i, data := range cache {
			d := *data.(*interface{})
			if d == nil {
				item[columns[i]] = ""
				continue
			}
			item[columns[i]] = string(d.([]byte)) //取实际类型
		}
		list = append(list, item)
	}
	_ = rows.Close()
	return list, nil
}

// selectGetTable 从select语句获取 table name
func (m *Mysql) selectGetTable(sql string) string{
	tList := strings.Split(sql,"from ")
	if len(tList) > 1{
		tList2 := strings.Split(tList[1]," ")
		if len(tList2) > 1{
			return tList2[0]
		}
	}
	return ""
}

// NewTable 创建表
// 字段顺序不固定
// fields  字段:类型； name:varchar(10);
func (m *Mysql) NewTable(table string, fields map[string]string) error {
	var (
		createSql bytes.Buffer
		line = len(fields)
	)

	if table == ""{
		return errors.New("table is null")
	}
	if line < 1{
		return errors.New("fiedls len is 0")
	}
	if m.DB == nil{
		_=m.Conn()
	}

	createSql.WriteString("CREATE TABLE ")
	createSql.WriteString(table)
	createSql.WriteString(" ( id int(11) NOT NULL AUTO_INCREMENT, ")
	for k,v := range fields{
		createSql.WriteString(k)
		createSql.WriteString(" ")
		createSql.WriteString(v)
		createSql.WriteString(", ")
	}
	createSql.WriteString("PRIMARY KEY (id) ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	_,err :=  m.DB.Exec(createSql.String())
	if m.log{
		loger("[Sql] Exec : " + createSql.String())
		if err != nil{
			loger("[Sql] Error : " + err.Error())
		}
	}
	if m.allTN == nil {
		_=m.allTableName()
	}
	m.allTN.add(table)
	return nil
}


// TODO: NewTable - 2 创建表  字段顺序需要固定


func (m *Mysql) insert(table string, fieldData map[string]interface{}) error {
	var insertSql bytes.Buffer
	_,_=m.Describe(table)
	describe := m.tableTemp[table]
	isDescribe := describe.Base != nil
	insertSql.WriteString("insert ")
	insertSql.WriteString(table)
	fieldStr := " ("
	valueStr := " ("

	for k,v := range fieldData {
		vType,ok := describe.Base[k];
		if isDescribe && !ok {
			continue
		}
		fieldStr += k+", "
		vStr := StringValueMysql(v)
		_,isStr := v.(string)
		if vType == "string" && !isStr {
			vStr = "'"+vStr+"'"
		}
		valueStr += vStr + ", "
	}

	if fieldStr[len(fieldStr)-2:] == ", " {
		fieldStr = fieldStr[:len(fieldStr)-2]
	}
	if valueStr[len(valueStr)-2:] == ", " {
		valueStr = valueStr[:len(valueStr)-2]
	}

	insertSql.WriteString(fieldStr)
	insertSql.WriteString(") VALUES ")
	insertSql.WriteString(valueStr)
	insertSql.WriteString(");")
	_, err := m.DB.Exec(insertSql.String())

	if m.log{
		loger("[Sql] Exec : " + insertSql.String())
		if err != nil{
			loger("[Sql] Error : " + err.Error())
		}
	}
	return err
}


// Insert 新增数据
func (m *Mysql) Insert(table string, fieldData map[string]interface{}) error {
	var (
		line = len(fieldData)
	)
	if table == ""{
		return errors.New("table is null")
	}
	if line < 1{
		return errors.New("fiedls len is 0")
	}
	if m.DB == nil{
		_=m.Conn()
	}
	return m.insert(table, fieldData)
}

// InsertAt 新增数据 如果没有表则先创建表
func (m *Mysql) InsertAt(table string, fieldData map[string]interface{}) error{
	var (
		line = len(fieldData)
	)
	if table == ""{
		return errors.New("table is null")
	}
	if line < 1{
		return errors.New("fiedls len is 0")
	}
	if m.DB == nil{
		_=m.Conn()
	}
	if m.allTN == nil {
		_=m.allTableName()
	}
	if !m.allTN.isHave(table) {
		newField := make(map[string]string)
		for k,v := range fieldData {
			newField[k] = dataType2Mysql(v)
		}
		err := m.NewTable(table, newField)
		if err != nil {
			return err
		}
	}
	return m.insert(table, fieldData)
}

// TODO: 新增数据结构体

// TODO: 新增数据json字符串

// Update
func (m *Mysql) Update(sql string) error {
	if m.DB == nil{
		_=m.Conn()
	}
	_, err := m.DB.Exec(sql)
	if m.log{
		loger("[Sql] Exec : " + sql)
		if err != nil{
			loger("[Sql] Error : " + err.Error())
		}
	}
	return err
}

// Exec
func (m *Mysql) Exec(sql string) error {
	if m.DB == nil{
		_=m.Conn()
	}
	_, err := m.DB.Exec(sql)
	if m.log{
		loger("[Sql] Exec : " + sql)
		if err != nil{
			loger("[Sql] Error : " + err.Error())
		}
	}
	return err
}

// Query
func (m *Mysql) Query(sql string) ([]map[string]string, error) {
	if m.DB == nil{
		_=m.Conn()
	}

	rows,err := m.DB.Query(sql)
	if m.log{
		log.Println("[Sql] Exec : " + sql)
		if err != nil{
			log.Println("[Sql] Error : " + err.Error())
		}
	}
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil{
		return nil,err
	}

	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache { //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}

	var list []map[string]string //返回的切片
	for rows.Next() {
		_ = rows.Scan(cache...)
		item := make(map[string]string)
		for i, data := range cache {
			d := *data.(*interface{})
			if d == nil {
				item[columns[i]] = ""
				continue
			}
			item[columns[i]] = string(d.([]byte)) //取实际类型
		}
		list = append(list, item)
	}
	_ = rows.Close()
	return list, nil
}

// Delete
func (m *Mysql) Delete(sql string) error {
	_, err := m.DB.Exec(sql)
	if m.log{
		loger("[Sql] Exec : " + sql)
		if err != nil{
			loger("[Sql] Error : " + err.Error())
		}
	}
	return err
}

// ToVarChar  写入mysql 的字符类型
func (m *Mysql) ToVarChar(data interface{}) string {
	var txt bytes.Buffer
	txt.WriteString(`"`)
	txt.WriteString(Any2String(data))
	txt.WriteString(`"`)
	return txt.String()
}

// dataType2Mysql
func dataType2Mysql(value interface{}) string{
	typ := reflect.ValueOf(value)
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	loger(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
	return "text"
}

// TODO 删除表

// TODO 判断是否存