/*
	Description : websocket的连接, 模拟websocket客户端
	Author : ManGe
	Mail : 2912882908@qq.com
	Github : https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"golang.org/x/net/websocket"
)

// WSClient websocket 客户端
type WSClient interface {
	Send(body []byte) error
	Read(data []byte) error
	Close()
}

// WsClient websocket 客户端
func WsClient(host, path string, isSSL bool) (WSClient, error) {
	ws := &webSocketClient{
		Host:  host,
		Path:  path,
		IsSSL: isSSL,
	}
	err := ws.conn()
	return ws, err
}

type webSocketClient struct {
	Host  string
	Path  string
	Ws    *websocket.Conn
	IsSSL bool
}

func (c *webSocketClient) conn() error {
	var err error
	u := c.Host + c.Path
	if c.IsSSL {
		c.Ws, err = websocket.Dial(u, "", "https://"+c.Host+"/")
	} else {
		c.Ws, err = websocket.Dial(u, "", "http://"+c.Host+"/")
	}
	return err
}

func (c *webSocketClient) Send(body []byte) error {
	_, err := c.Ws.Write(body)
	return err
}

func (c *webSocketClient) Close() {
	_ = c.Ws.Close()
}

func (c *webSocketClient) Read(data []byte) error {
	_, err := c.Ws.Read(data)
	return err
}
