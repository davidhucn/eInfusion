package http

import (
	"sync"

	ws "github.com/gorilla/websocket"
)

type webMsg struct {
	WSReceiveDataError string
	WSConnectError     string
	WSSendDataError    string
}

// WebMsg :Web信息对象
var WebMsg webMsg

type reqData struct {
	ID      string `json:"ID"`
	CmdType string `json:"CmdType"` //指令类型(代码)
	Args    string `json:"Args"`    //相关参数 (例如：ip、port)
	// Action string `json:"-"`
}

// clisData :客户端请求对象，内部用
var clisData []reqData

// Cmd :指令对象
type Cmd struct {
	Cmd   []byte
	CmdID string
}

func init() {
	WebMsg.WSConnectError = "错误，websocket连接错误！"
	WebMsg.WSSendDataError = "错误，websocket发送数据错误！"
	WebMsg.WSReceiveDataError = "错误，websocket接收数据失败！"
}

// websocketConnectionns :web客户端对象
type webSocketConnections struct {
	WsConnection *ws.Conn
	Orders       chan *Cmd
}

// WebClients :web连接客户端
type WebClients struct {
	Connections map[string]*webSocketConnections
	sync.Mutex
}

// func NewWebClients()*WebClients{
// 	return &WebClients{
// 		Connections : make(map[string]*webSocketConnections)
// 	}
// }
