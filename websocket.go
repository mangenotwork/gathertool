/*
	Description : websocket的连接
	Author : ManGe
	Version : v0.1
	Date : 2021-04-24
*/

package gathertool

import (
	"golang.org/x/net/websocket"
)

type WSClient interface {
	Send(body []byte) error
	Read(data []byte) error
	Close()
}

func WsClient(host, path string) (WSClient, error) {
	ws := &webSocketClient{
		Host: host,
		Path: path,
	}
	err := ws.conn()
	return ws, err
}


type webSocketClient struct {
	Host string
	Path string
	Ws *websocket.Conn
}

func (c *webSocketClient) conn() error {
	var err error
	u := c.Host + c.Path
	c.Ws, err = websocket.Dial(u, "", "https://"+c.Host+"/")
	return err
}

func (c *webSocketClient) Send(body []byte) error {
	_, err := c.Ws.Write(body)
	return err
}

func (c *webSocketClient) Close() {
	c.Ws.Close()
}

func (c *webSocketClient) Read(data []byte) error{
	_, err := c.Ws.Read(data)
	return err
}

