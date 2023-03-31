/*
	Description : Tcp的连接 (Tcp客户端); 应用场景是模拟Tcp客户端;
	Author : ManGe
	Mail : 2912882908@qq.com
	Github : https://github.com/mangenotwork/gathertool

// ===== Use
func main(){
	client := gt.NewTcpClient()
	client.Run("192.168.0.9:29123", f)
}

func f(client *gt.TcpClient){
	go func() {
		// 发送登录请求
		_,err := client.Send([]byte(`{
			"cmd":"Auth",
			"data":{
				"account":"a10",
				"password":"123456",
				"device":"1",
				"source":"windows"
			}
		}`))
		if err != nil {
			log.Println("err = ", err)
		}
	}()
}

*/

package gathertool

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// TcpClient Tcp客户端
type TcpClient struct {
	Connection *net.TCPConn
	HawkServer *net.TCPAddr
	StopChan   chan struct{}
	CmdChan    chan string
	Token      string
	RConn      chan struct{}
}

func NewTcpClient() *TcpClient {
	return new(TcpClient)
}

func (c *TcpClient) Send(b []byte) (int, error) {
	if c.Connection == nil {
		return 0, fmt.Errorf("conn is null")
	}
	return c.Connection.Write(b)
}

func (c *TcpClient) Read(b []byte) (int, error) {
	if c.Connection == nil {
		return 0, fmt.Errorf("conn is null")
	}
	return c.Connection.Read(b)
}

func (c *TcpClient) Addr() string {
	if c.Connection == nil {
		return ""
	}
	return c.Connection.RemoteAddr().String()
}

func (c *TcpClient) Close() {
	if c.Connection == nil {
		return
	}
	_ = c.Connection.Close()
}

func (c *TcpClient) Stop() {
	c.StopChan <- struct{}{}
}

func (c *TcpClient) ReConn() {
	c.RConn <- struct{}{}
}

func (c *TcpClient) Run(serverHost string, r func(c *TcpClient, data []byte), w func(c *TcpClient)) {
	//用于重连
Reconnection:

	hawkServer, err := net.ResolveTCPAddr("tcp", serverHost)
	if err != nil {
		log.Printf("hawk server [%s] resolve error: [%s]", serverHost, err.Error())
		time.Sleep(1 * time.Second)
		goto Reconnection
	}

	//连接服务器
	connection, err := net.DialTCP("tcp", nil, hawkServer)
	if err != nil {
		log.Printf("connect to hawk server error: [%s]", err.Error())
		time.Sleep(1 * time.Second)
		goto Reconnection
	}
	log.Println("[连接成功] 连接服务器成功")

	//创建客户端实例
	c.Connection = connection
	c.HawkServer = hawkServer
	c.StopChan = make(chan struct{})
	c.CmdChan = make(chan string)
	c.RConn = make(chan struct{})

	//启动接收
	go func(conn *TcpClient) {
		for {
			recv := make([]byte, 1024)
			for {
				n, err := conn.Connection.Read(recv)
				if err != nil {
					if err == io.EOF {
						log.Println(conn.Addr(), " 断开了连接!")
					}
					conn.Close()
					c.RConn <- struct{}{}
					return
				}
				if n > 0 && n < 1025 {
					log.Println(string(recv[:n]))
					r(c, recv[:n])
				}

			}
		}
	}(c)

	go w(c)

	for {
		select {
		case a := <-c.RConn:
			log.Println("global.RConn = ", a)
			goto Reconnection
		case <-c.StopChan:
			return

		}
	}
}
