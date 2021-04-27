package gathertool

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)


var MConn *sql.DB

type Mysql struct {
	Host string
	Port int
	User string
	Password string
	DataBase string
	MaxOpenConn int
	MaxIdleConn int
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

func (m *Mysql) Conn() (err error){
	MConn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s",
		m.User, m.Password, "TCP", m.Host, m.Port, m.DataBase))
	if err != nil {
		return err
	}
	MConn.SetConnMaxLifetime(100*time.Second)  //最大连接周期，超过时间的连接就close
	if m.MaxOpenConn < 1{
		m.MaxOpenConn = 10
	}
	if m.MaxIdleConn < 1{
		m.MaxIdleConn = 5
	}
	MConn.SetMaxOpenConns(m.MaxOpenConn)//设置最大连接数
	MConn.SetMaxIdleConns(m.MaxIdleConn) //设置闲置连接数
	return
}

func (m *Mysql) First(sql string) (map[string]interface{}, error) {
	if MConn == nil{
		_=m.Conn()
	}
	row,err := MConn.Query(sql)
	if err != nil {
		return nil, err
	}
	columns, err := row.Columns()
	log.Println(columns,err)
	return nil, nil
}