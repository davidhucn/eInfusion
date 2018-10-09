package http

import (
	"os"
	"sync"

	ss "github.com/gorilla/sessions"
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

// WebClients :web连接客户端
type WebClients struct {
	Connections map[string]*ws.Conn
	Orders      chan *Cmd
	sync.Mutex
}

// CStore :全局cookie记录对象
var CStore = ss.NewCookieStore([]byte(os.Getenv("Session-Key")))

// NewOrder :生成新的命令对象
func NewOrder(rID string, rCnt []byte) *Cmd {
	return &Cmd{
		Cmd:   rCnt,
		CmdID: rID,
	}
}

// NewWebClients :创建新的WebClient对象
func NewWebClients() *WebClients {
	return &WebClients{
		Connections: make(map[string]*ws.Conn),
		Orders:      make(chan *Cmd, 1024),
	}
}
