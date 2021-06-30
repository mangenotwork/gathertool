package gathertool

import (
	"golang.org/x/net/websocket"
)

type WSClient interface {

}

type WebSocketClient struct {
	Host string
	Path string
	Ws *websocket.Conn
	Err error
}

func (c *WebSocketClient) Conn() *WebSocketClient {
	u := c.Host + c.Path
	c.Ws, c.Err = websocket.Dial(u, "", "https://"+c.Host+"/")
	//log.Println(ws)
	//c.Ws =
	//if err != nil || ws == nil {
	//	log.Println(err)
	//	return
	//}
	//defer ws.Close() //关闭连接
	//
	//_, err = ws.Write(body)
	//if err != nil {
	//	log.Println(err)
	//}
	return c
}

func (c *WebSocketClient) Send(body []byte) *WebSocketClient{
	_, c.Err = c.Ws.Write(body)
	return c
}

func (c *WebSocketClient) Close() {
	c.Ws.Close()
}