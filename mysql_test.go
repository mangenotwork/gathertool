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
	//log.Println(db,err)
	//err = db.Conn()
	log.Println(err)
	log.Println(db.IsHaveTable("task"))
	//data, err := db.Select("select * from task limit 2;")
	//log.Println(err, len(data))
	//for _,v := range data{
	//	log.Println(v)
	//	log.Println("___________")
	//}
	//
	////db.Describe("task")
	//
	//db.NewTable("test", map[string]string{
	//	"name": "varchar(100)",
	//	"age": "int(11)",
	//	"txt": "text",
	//})
	//
	//db.Insert("test", map[string]interface{}{
	//	"name": "mange",
	//	"age": 22,
	//	"txt": "texasdasdasdasdasdasdasdsaddt",
	//})

	//data, err := db.Query("show master status")
	//log.Println(data, err)

	//data, err := db.Query("show binary logs")
	//log.Println(data, err)

	//data, err := db.Query("show binlog events in 'mysql-bin.000033' limit 10")
	//log.Println( err)
	//for _, v := range data{
	//	log.Println(v)
	//	log.Println("Log_name", v["Log_name"])
	//	log.Println("Pos", v["Pos"])
	//	log.Println("Event_type", v["Event_type"])
	//	log.Println("Server_id", v["Server_id"])
	//	log.Println("End_log_pos", v["End_log_pos"])
	//	log.Println("Info", v["Info"])
	//}
}