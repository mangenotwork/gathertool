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
	err = db.Conn()
	log.Println(err)
	data, err := db.Select("select * from task limit 2;")
	log.Println(err, len(data))
	for _,v := range data{
		log.Println(v)
		log.Println("___________")
	}

	//db.Describe("task")

	db.NewTable("test", map[string]string{
		"name": "varchar(100)",
		"age": "int(11)",
		"txt": "text",
	})

	db.Insert("test", map[string]interface{}{
		"name": "mange",
		"age": 22,
		"txt": "texasdasdasdasdasdasdasdsaddt",
	})

}