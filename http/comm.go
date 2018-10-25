package http

import (
	cm "eInfusion/comm"
	"sync"

	ws "github.com/gorilla/websocket"
)

type webMsg struct {
	WSReceiveDataError        string
	WSConnectError            string
	WSSendDataError           string
	WSSendDataFailureTryLater string
	WSSendDataSuccess         string
}

// WebMsg :Web信息对象
var WebMsg webMsg

type reqData struct {
	TargetID string `json:"TargetID"`
	CmdType  string `json:"CmdType"` //指令类型(代码)
	Args     string `json:"Args"`    //相关参数 (例如：ip、port)
	// Action string `json:"-"`
}

// clisData :客户端请求对象，内部用
var clisData []reqData

func init() {
	WebMsg.WSConnectError = "错误：websocket连接错误！"
	WebMsg.WSSendDataError = "错误：websocket发送数据错误！"
	WebMsg.WSReceiveDataError = "错误：websocket接收数据失败！"
	WebMsg.WSSendDataFailureTryLater = "提示：由于设备未连线等原因，发送指令至设备失败，系统稍候后将自动尝试！"
	WebMsg.WSSendDataSuccess = "提示：完成数据发送！"
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
