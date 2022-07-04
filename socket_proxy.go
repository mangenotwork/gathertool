/*
	Description : 启动一个socket5代理
	Author : ManGe
			2912882908@qq.com
			https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"io"
	"log"
	"net"
	"strconv"
)

// SockerProxy 启动一个socket5代理
func SockerProxy(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleClientRequest2(client)
	}
}

func handleClientRequest2(client net.Conn) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	log.Println("socket 请求 : ", client.RemoteAddr(), " --> ", client.LocalAddr())

	if client == nil {
		return

	}
	defer client.Close()
	var b [1024*100]byte

	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}

	if b[0] == 0x05 { //只处理Socks5协议
		//客户端回应：Socks服务端不需要验证方式
		client.Write([]byte{0x05, 0x00})
		n, err = client.Read(b[:])
		var host, port string

		switch b[3] {
		case 0x01: //IP V4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()

		case 0x03: //域名
			host = string(b[5 : n-2]) //b[4]表示域名的长度

		case 0x04: //IP V6
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}

		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

		server, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			log.Println(err)
			return

		}
		defer server.Close()
		_,_= client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功

		//进行转发
		go io.Copy(server, client)
		io.Copy(client, server)
	}
}
