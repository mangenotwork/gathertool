package gathertool

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

// RedisSSHConn SSH普通连接
// addr : SSH主机地址, 如: 127.0.0.1:22
// user : SSH用户
// pass : SSH密码
// ip : redis服务地址
// port:  Redis 服务端口
// password  Redis 服务密码
// 返回redis连接
func RedisSSHConn(user string, pass string, addr string, ip string, port int, password string) (redis.Conn, error) {
	sshClient, err := SSHClient(user, pass, addr)
	if nil != err {
		fmt.Println(err)
		return nil, err
	}

	host := fmt.Sprintf("%s:%d", ip, port)
	conn, err := sshClient.Dial("tcp", host)
	if nil != err {
		return nil, err
	}

	redisConn := redis.NewConn(conn, -1, -1)

	//TODO: 针对无密码的连接
	if password == "" {
		return redisConn, nil
	}

	if _, authErr := redisConn.Do("AUTH", password); authErr != nil {
		return nil, fmt.Errorf("redis auth password error: %s", authErr)
	}

	return redisConn, nil
}