package gathertool

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"time"
)

type Mysql struct {
	Host string
	Port int
	User string
	Password string
	DataBase string
	MaxOpenConn int
	MaxIdleConn int
	DB *sql.DB
}

func NewMysql(host string,port int, user, password, database string) (*Mysql, error) {
	if len(host) < 1 {
		return nil, errors.New("Host is Null.")
	}
	if port < 1 {
		port = 3369
	}
	return &Mysql{
		Host : host,
		Port : port,
		User : user,
		Password : password,
		DataBase : database,
	}, nil
}

// 连接mysql
func (m *Mysql) Conn() (err error){
	m.DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s",
		m.User, m.Password, "tcp", m.Host, m.Port, m.DataBase))
	if err != nil {
		return err
	}
	m.DB.SetConnMaxLifetime(100*time.Second)  //最大连接周期，超过时间的连接就close
	if m.MaxOpenConn < 1{
		m.MaxOpenConn = 10
	}
	if m.MaxIdleConn < 1{
		m.MaxIdleConn = 5
	}
	m.DB.SetMaxOpenConns(m.MaxOpenConn)//设置最大连接数
	m.DB.SetMaxIdleConns(m.MaxIdleConn) //设置闲置连接数
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

// Describe 获取表结构
func (m *Mysql) Describe(table string) (map[string]string, error){
	if m.DB == nil{
		_=m.Conn()
	}

	if table == ""{
		return nil, errors.New("table name is null.")
	}

	rows,err := m.DB.Query("DESCRIBE " + table)
	if err != nil{
		return nil, err
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

	return fieldMap, nil
}

// Select 查询语句 返回 map
func (m *Mysql) Select(sql string) ([]map[string]string, error) {
	if m.DB == nil{
		_=m.Conn()
	}
	log.Println("DB = ", m.DB)

	rows,err := m.DB.Query(sql)
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
	log.Println(createSql.String())
	res,err :=  m.DB.Exec(createSql.String())
	log.Println(res,err)
	return nil
}

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

	insertSql.WriteString("insert ")
	insertSql.WriteString(table)
	fieldSql.WriteString(" (")
	valueSql.WriteString(" (")
 	for k,v := range fieldData {
		fieldSql.WriteString(k)
		valueSql.WriteString(StringValue(v))
		n++
		if n < line{
			fieldSql.WriteString(", ")
			valueSql.WriteString(", ")
		}
	}

	insertSql.WriteString(fieldSql.String())
	insertSql.WriteString(") VALUES ")
	insertSql.WriteString(valueSql.String())
	insertSql.WriteString(");")
	log.Println(insertSql.String())
	_, err := m.DB.Exec(insertSql.String())
 	return err
}

