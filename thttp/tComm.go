package thttp

// import (
// 	"sync"

// 	ws "github.com/gorilla/websocket"
// )

// type reqData struct {
// 	ID      string `json:"ID"`
// 	CmdType string `json:"CmdType"` //指令类型(代码)
// 	Args    string `json:"Args"`    //相关参数 (例如：ip、port)
// 	// Action string `json:"-"`
// }

// // clisData :客户端请求对象，内部用
// var clisData []reqData

// // WSConnet :全局ws连接对象
// type WSConnet struct {
// 	conn   *ws.Conn    // websocket 连接器
// 	sdData chan []byte // 发送信息的缓冲 channel
// }

// // WsClis :全局ws连接对象集
// var WsClis map[string]*WSConnet

// // 定义锁
// var (
// 	cMkMutex  sync.Mutex
// 	cDelMutex sync.Mutex
// )

// // sOrders :全局发送指令map数组,MAP索引为时间戳
// var sdOrders map[string][]byte

// // sdIDStream :发送指令标识,涵盖IP地址及唯一标识字符
// var sdIDStream chan string

// func init() {
// 	sdOrders = make(map[string][]byte)
// 	sdIDStream = make(chan string, 1024)
// 	WsClis = make(map[string]*WSConnet)
// }
