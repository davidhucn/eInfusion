package tqueue

import "sync"

// ReqCmd :指令类型
type ReqCmd struct {
	TargetID chan string
	CmdType  chan uint8  // 指令类型(代码)
	Args     chan string // 相关参数 (例如：ip、port)
}

//定义锁
var (
	cMkMutex  sync.Mutex
	cDelMutex sync.Mutex
)

// sOrders :全局指令map,MAP索引为时间戳
var sdOrders map[string][]byte

var sdIDStream chan string

func init() {
	sdOrders = make(map[string][]byte)
	sdIDStream = make(chan string, 1024)
}
