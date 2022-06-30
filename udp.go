/*
	Description : Udp的连接 (Udp客户端); 应用场景是模拟Udp客户端;
	Author : ManGe
			2912882908@qq.com
			https://github.com/mangenotwork/gathertool
 */

package gathertool

import (
	"fmt"
	"log"
	"net"
)

// UdpClient Udp客户端
type UdpClient struct {
	SrcAddr *net.UDPAddr
	DstAddr *net.UDPAddr
	Conn *net.UDPConn
	Token string
}

func (u *UdpClient) Run(hostServer string, port int, r func(u *UdpClient, data []byte), w func(u *UdpClient)) {
	var err error
	sip := net.ParseIP(hostServer)
	u.SrcAddr = &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	u.DstAddr = &net.UDPAddr{IP: sip, Port: port}
	u.Conn, err = net.DialUDP("udp", u.SrcAddr, u.DstAddr)
	if err != nil {
		log.Println(err)
	}
	log.Println("连接成功; c = ", u.Conn)
	defer u.Conn.Close()

	go func() {
		data := make([]byte, 1024)
		for{
			n, remoteAddr, err := u.Conn.ReadFromUDP(data)
			if err != nil {
				log.Printf("error during read: %s", err)
			}
			log.Printf("<%s> %s\n", remoteAddr, data[:n])
			r(u, data[:n])
		}
	}()

	go w(u)

	select {}
}

func NewUdpClient() *UdpClient {
	return new(UdpClient)
}

func (u *UdpClient) Send(b []byte) (int, error) {
	if u.Conn == nil {
		return 0, fmt.Errorf("conn is null")
	}
	return u.Conn.Write(b)
}

func (u *UdpClient) Read(b []byte) (int, *net.UDPAddr, error) {
	if u.Conn == nil {
		return 0, u.DstAddr, fmt.Errorf("conn is null")
	}
	return u.Conn.ReadFromUDP(b)
}

func (u *UdpClient) Addr() string {
	return u.DstAddr.String()
}

func (u *UdpClient) Close(){
	if u.Conn == nil {
		return
	}
	_ = u.Conn.Close()
}
