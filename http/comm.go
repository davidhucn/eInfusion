package http

import (
	cm "eInfusion/comm"
	"strings"
	"sync"

	ws "github.com/gorilla/websocket"
)

type webMsg struct {
	WSReceiveDataError string
	WSConnectError     string
	WSSendDataError    string
	WSSendDataSuccess  string
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

func init() {
	WebMsg.WSConnectError = "错误，websocket连接错误！"
	WebMsg.WSSendDataError = "错误，websocket发送数据错误！"
	WebMsg.WSReceiveDataError = "错误，websocket接收数据失败！"
	WebMsg.WSSendDataSuccess = "提示，完成数据发送！"
}

// WebClients :web连接客户端
type WebClients struct {
	Connections map[string]*ws.Conn
	Orders      chan *cm.Cmd
	sync.Mutex
}

// CStore :全局cookie记录对象
// var CStore = ss.NewCookieStore([]byte(os.Getenv("Session-Key")))

// NewWebClients :创建新的WebClient对象
func NewWebClients() *WebClients {
	return &WebClients{
		Connections: make(map[string]*ws.Conn),
		Orders:      make(chan *cm.Cmd, 1024),
	}
}

// NewWSOrderID :生成新的websocket消息编号
func NewWSOrderID(rWSConnectionID string) string {
	return rWSConnectionID + "#" + GetRandString(8)
}

// DecodeToWSConnID :解析生成websocket连接序号
func DecodeToWSConnID(rWSOrderID string) string {
	return strings.Split(rWSOrderID, "#")[0]
}
