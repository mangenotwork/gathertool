package main

import (
	gt "github.com/mangenotwork/gathertool"
	"log"
)

func main() {
	DelKeys()
}

func DelKeys() {
	log.Println("run DelKeys")
	var (
		ssh_user       = ""
		ssh_password   = ""
		ssh_addr       = ""
		redis_host     = "127.0.0.1"
		redis_port     = "6379"
		redis_password = "123"
		dbnumber       = 0
	)
	conn, err := gt.NewRedis(redis_host, redis_port, redis_password, dbnumber, gt.NewSSHInfo(ssh_user, ssh_password, ssh_addr))
	log.Println(conn, err)

	//rds := gt.NewRedisPool(redis_host, redis_port, redis_password, dbnumber, 5,10,10,
	//	gt.NewSSHInfo(ssh_addr, ssh_user, ssh_password))
	//
	//gt.RedisDELKeys(rds, "in:*", 100)

	for _, v := range []string{"1", "2", "3", "4", "5", "6", "7"} {
		err = conn.MqProducer("mqtest", v)
		log.Println(err)
	}

	data, err := conn.ToString(conn.MqConsumer("mqtest"))
	log.Println(data, err)

	//for conn.MqLen("mqtest") > 0 {
	//	data, err := conn.MqConsumer("mqtest")
	//	log.Println(data, err)
	//}

}
