/*
	Description : mysql 相关方法
	Author : ManGe
			2912882908@qq.com
			https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize/v2"
)

var (
	SHOW_TABLES = "SHOW TABLES"
	TABLE_NAME_NULL = fmt.Errorf("table name is null.")
	TABLE_IS_NULL = fmt.Errorf("table is null.")
)

var MysqlDB = &Mysql{}

// Mysql 客户端对象
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
	allTN *allTableName // 所有表名
}

// NewMysqlDB 给mysql对象进行连接
func NewMysqlDB(host string, port int, user, password, database string) (err error) {
	MysqlDB, err = NewMysql(host,port, user, password, database)
	if err != nil {
		return
	}
	return MysqlDB.Conn()
}

// NewMysql 创建一个mysql对象
func NewMysql(host string, port int, user, password, database string) (*Mysql, error) {
	if len(host) < 1 {
		return nil, fmt.Errorf("Host is Null.")
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
			Error("[Sql] Conn Fail : " + err.Error())
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

// allTableName 所有表名
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
		//log.Println(err, result)
		m.allTN.add(result)
	}

	_=rows.Close()
	return
}

// IsHaveTable 表是否存在
func (m *Mysql) IsHaveTable(table string) bool {
	if m.allTN == nil {
		_=m.allTableName()
	}
	return m.allTN.isHave(table)
}

// TableInfo 表信息
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

// allTableName 记录当前库的所有表名
type allTableName struct{
	mut *sync.Mutex
	tableName map[string]struct{}
}

// newAllTableName
func newAllTableName() *allTableName {
	return &allTableName{
		mut: &sync.Mutex{},
		tableName: make(map[string]struct{}),
	}
}

// add 添加
func (a *allTableName) add(name string) *allTableName{
	a.mut.Lock()
	a.tableName[name] = struct{}{}
	a.mut.Unlock()
	return a
}

// remove 移除
func (a *allTableName) remove(name string) *allTableName{
	a.mut.Lock()
	delete(a.tableName, name)
	a.mut.Unlock()
	return a
}

// isHave 是否存在
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
		return &tableDescribe{}, TABLE_NAME_NULL
	}

	rows,err := m.DB.Query("DESCRIBE " + table)
	if err != nil{
		return &tableDescribe{}, err
	}

	fieldMap := make(map[string]string,0)
	for rows.Next() {
		result := &TableInfo{}
		err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
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
		Info("[Sql] Exec : " + sql)
		if err != nil{
			Error("[Sql] Error : " + err.Error())
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
		return TABLE_IS_NULL
	}

	if line < 1{
		return fmt.Errorf("fiedls len is 0")
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
		Info("[Sql] Exec : " + createSql.String())
		if err != nil{
			Error("[Sql] Error : " + err.Error())
		}
	}

	if m.allTN == nil {
		_=m.allTableName()
	}

	m.allTN.add(table)
	return nil
}

// NewTableGd 创建新的固定map顺序为字段的表
func (m *Mysql) NewTableGd(table string, fields *gDMap) error {
	var (
		createSql bytes.Buffer
		line = fields.Len()
	)

	if table == ""{
		return TABLE_IS_NULL
	}

	if line < 1{
		return fmt.Errorf("fiedls len is 0")
	}

	if m.DB == nil{
		_=m.Conn()
	}

	createSql.WriteString("CREATE TABLE ")
	createSql.WriteString(table)
	createSql.WriteString(" ( id int(11) NOT NULL AUTO_INCREMENT, ")
	fields.Range(func(k string, v interface{}){
		createSql.WriteString(k)
		createSql.WriteString(" ")
		createSql.WriteString(dataType2Mysql(v))
		createSql.WriteString(", ")
	})

	createSql.WriteString("PRIMARY KEY (id) ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	_,err :=  m.DB.Exec(createSql.String())
	if m.log{
		Info("[Sql] Exec : " + createSql.String())
		if err != nil{
			Error("[Sql] Error : " + err.Error())
		}
	}

	if m.allTN == nil {
		_=m.allTableName()
	}

	m.allTN.add(table)
	return nil
}

// insert 插入操作
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
		Info("[Sql] Exec : " + insertSql.String())
		if err != nil{
			Error("[Sql] Error : " + err.Error())
		}
	}

	return err
}

// Insert 新增数据
func (m *Mysql) Insert(table string, fieldData map[string]interface{}) error {
	var line = len(fieldData)

	if table == ""{
		return TABLE_IS_NULL
	}

	if line < 1{
		return fmt.Errorf("fiedls len is 0")
	}

	if m.DB == nil{
		_=m.Conn()
	}

	if m.allTN == nil {
		_=m.allTableName()
	}

	if !m.allTN.isHave(table) {
		return fmt.Errorf("[Insert Err] Table not fond.")
	}
	return m.insert(table, fieldData)
}

// InsertAt 新增数据 如果没有表则先创建表
func (m *Mysql) InsertAt(table string, fieldData map[string]interface{}) error{
	var line = len(fieldData)

	if table == ""{
		return TABLE_IS_NULL
	}

	if line < 1{
		return fmt.Errorf("fiedls len is 0")
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

// InsertAtGd  固定顺序map写入
func (m *Mysql) InsertAtGd(table string, fieldData *gDMap) error{
	var (
		line = fieldData.Len()
		fieldDataMap = make(map[string]interface{})
	)

	if table == ""{
		return TABLE_IS_NULL
	}

	if line < 1{
		return fmt.Errorf("fiedls len is 0")
	}

	if m.DB == nil{
		_=m.Conn()
	}

	if m.allTN == nil {
		_=m.allTableName()
	}

	if !m.allTN.isHave(table) {
		err := m.NewTableGd(table, fieldData)
		if err != nil {
			return err
		}
	}

	fieldData.Range(func(k string, v interface{}){
		fieldDataMap[k] = v
	})

	return m.insert(table, fieldDataMap)
}

// InsertAtJson json字符串存入数据库
func (m *Mysql) InsertAtJson(table, jsonStr string) error {
	data := Any2Map(jsonStr)
	return m.InsertAt(table, data)
}

// Update 更新sql
func (m *Mysql) Update(sql string) error {
	if m.DB == nil{
		_=m.Conn()
	}

	_, err := m.DB.Exec(sql)
	if m.log{
		Info("[Sql] Exec : " + sql)
		if err != nil{
			Error("[Sql] Error : " + err.Error())
		}
	}
	return err
}

// Exec 执行sql
func (m *Mysql) Exec(sql string) error {
	if m.DB == nil{
		_=m.Conn()
	}

	_, err := m.DB.Exec(sql)
	if m.log{
		Info("[Sql] Exec : " + sql)
		if err != nil{
			Error("[Sql] Error : " + err.Error())
		}
	}
	return err
}

// Query 执行selete sql
func (m *Mysql) Query(sql string) ([]map[string]string, error) {
	if m.DB == nil{
		_=m.Conn()
	}

	rows,err := m.DB.Query(sql)
	if m.log{
		Info("[Sql] Exec : " + sql)
		if err != nil{
			Error("[Sql] Error : " + err.Error())
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

// Delete 执行delete sql
func (m *Mysql) Delete(sql string) error {
	if strings.Index(sql, "DELETE") != -1 || strings.Index(sql, "delete") != -1 {
		return m.Exec(sql)
	}
	return fmt.Errorf("请检查sql正确性")
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
	default:
		return "text"
	}
	Info(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
	return "text"
}

// DeleteTable 删除表
func (m *Mysql) DeleteTable(tableName string) error {
	err := m.Exec("DROP TABLE "+tableName)
	if err != nil {
		return err
	}
	m.allTN.remove(tableName)
	return nil
}

// HasTable 判断表是否存
func (m *Mysql) HasTable(tableName string) bool {
	return m.allTN.isHave(tableName)
}

// GetFieldList 获取表字段
func (m *Mysql) GetFieldList(table string) (fieldList []string) {
	fieldList = make([]string,0)

	if table == ""{
		return
	}

	rows,err := m.DB.Query("DESCRIBE " + table)
	if err != nil{
		return
	}

	for rows.Next() {
		result := &TableInfo{}
		err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
		fieldList = append(fieldList, result.Field)
	}
	_=rows.Close()

	return
}

// ToXls 数据库查询输出到excel
func (m *Mysql) ToXls(sql, outPath string) {
	data, err := m.Select(sql)
	if err != nil {
		Error(err)
		return
	}
	if len(data) < 1 {
		Error("查询数据为空")
		return
	}
	f := excelize.NewFile()
	count := len(data)
	var bar Bar
	bar.NewOption(0, int64(count-1))

	fields := make([]string,0)
	n := 1
	for k,_ := range data[0] {
		fields = append(fields, k)
		_=f.SetCellValue("Sheet1", toNumberSystem26(n)+"1", k)
		n++
	}
	// 写入数据
	for i:=0; i<count; i++ {
		n := 1
		for _, v := range fields {
			_=f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", toNumberSystem26(n), i+2 ), data[i][v])
			n++
		}
		bar.Play(int64(i))
	}
	bar.Finish()

	if err := f.SaveAs(outPath); err != nil {
		Error("[err] 导出失败: ", err)
		return
	}
	workPath, _ := os.Getwd()
	Info("[导出成功] 文件位置: ", workPath+"/"+outPath)
}

func toNumberSystem26(n int) string {
	s := ""
	for ;n>0; {
		m := n%26
		if m == 0 {
			m = 26
		}
		s = s + string(rune(m+64))
		n = (n - m )/26
	}
	return s
}