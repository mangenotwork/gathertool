package gathertool

import (
	"google.golang.org/grpc"
	"net"
)

func SimpleServer(addr string, g *grpc.Server){
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic("监听错误:" + err.Error())
	}
	err = g.Serve(lis)
	if err != nil {
		panic("启动错误:" + err.Error())
	}
}