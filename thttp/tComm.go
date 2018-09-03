package thttp

import ws "github.com/gorilla/websocket"

type reqData struct {
	ID      string `json:"ID"`
	CmdType string `json:"CmdType"` //指令类型(代码)
	Args    string `json:"Args"`    //相关参数 (例如：ip、port)
	// Action string `json:"-"`
}

// clisData :客户端请求对象，内部用
var clisData []reqData

// WSConnet :全局ws连接对象
type WSConnet struct {
	// websocket 连接器
	conn *ws.Conn
	// 发送信息的缓冲 channel
	sdData chan []byte
}

var ClisWS []*WSConnet
