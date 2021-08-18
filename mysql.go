/*
	Description : mysql 相关方法
	Author : ManGe
	Version : v0.3
	Date : 2021-08-18
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
	"time"
	"sync"
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
	tableTemp  map[string]tableDescribe //表结构缓存
	once *sync.Once
}

// 给mysql对象进行连接
func NewMysqlDB(host string,port int, user, password, database string)(err error){
	MysqlDB, err = NewMysql(host,port, user, password, database)
	err = MysqlDB.Conn()
	return
}

func NewMysql(host string,port int, user, password, database string) (*Mysql, error) {
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
		m.tableTemp = make(map[string]tableDescribe)
	})
	return m, nil
}

// 获取mysql 连接
func GetMysqlDBConn() (*Mysql,error) {
	err := MysqlDB.Conn()
	return MysqlDB, err
}

// 关闭日志
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

// 连接mysql
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

// Describe 获取表结构
func (m *Mysql) Describe(table string) (tableDescribe, error){
	if m.DB == nil{
		_=m.Conn()
	}

	if v,ok := m.tableTemp[table]; ok {
		return v, nil
	}

	if table == ""{
		return tableDescribe{}, errors.New("table name is null.")
	}

	rows,err := m.DB.Query("DESCRIBE " + table)
	if err != nil{
		return tableDescribe{}, err
	}
	defer rows.Close()
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

	td := tableDescribe{
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
	defer rows.Close()
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
	//_ = rows.Close()
	return list, nil
}

// 从select语句获取 table name
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
	return nil
}


// TODO: NewTable - 2 创建表  字段顺序需要固定



// Insert 新增数据
func (m *Mysql) Insert(table string, fieldData map[string]interface{}) error {
	var (
		insertSql bytes.Buffer
		fieldSql bytes.Buffer
		valueSql bytes.Buffer
		line = len(fieldData)
		n = 0
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

	// 表结构Describe
	m.Describe(table)
	describe := m.tableTemp[table]
	log.Println("[describe] = ", describe)
	isDescribe := describe.Base != nil
	insertSql.WriteString("insert ")
	insertSql.WriteString(table)
	fieldSql.WriteString(" (")
	valueSql.WriteString(" (")

	for k,v := range fieldData {
		n++
		vType,ok := describe.Base[k];
		if isDescribe && !ok {
			continue
		}
		fieldSql.WriteString(k)
		vStr := StringValueMysql(v)
		_,isStr := v.(string)
		if vType == "string" && !isStr {
			vStr = "'"+vStr+"'"
		}
		valueSql.WriteString(vStr)
		if n < line-1{
			fieldSql.WriteString(", ")
			valueSql.WriteString(", ")
		}
	}

	insertSql.WriteString(fieldSql.String())
	insertSql.WriteString(") VALUES ")
	insertSql.WriteString(valueSql.String())
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

// TODO: 新增数据结构体

// TODO: 新增数据json字符串

// 执行 Update
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

// 执行sql Exec
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

// 执行sql Query
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
	defer rows.Close()
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
	//_ = rows.Close()
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
func dataType2Mysql(typ reflect.Value) string{
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
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

// TODO 删除表

// TODO 判断是否存