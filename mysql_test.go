package gathertool

import (
	"log"
	"testing"
)

func TestConn(t *testing.T){
	var (
		host = "192.168.0.197"
		port = 3306
		user = "root"
		password = "root123"
		database = "spider"
	)
	db,err := NewMysql(host, port, user, password, database)
	log.Println(db,err)
}