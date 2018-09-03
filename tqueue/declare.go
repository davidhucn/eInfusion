package tqueue

import (
	"sync"
)

// ReqCmd :指令类型
// type ReqCmd struct {
// 	TargetID chan string
// 	CmdType  chan uint8  // 指令类型(代码)
// 	Args     chan string // 相关参数 (例如：ip、port)
// }

//定义锁
var (
	cMkMutex  sync.Mutex
	cDelMutex sync.Mutex
)

// sOrders :全局发送指令map数组,MAP索引为时间戳
var sdOrders map[string][]byte

// sdIDStream :发送指令标识，触发用
var sdIDStream chan string

// RetMsg :指令操作返回数据映射
var RetMsg map[string][]byte

func init() {
	sdOrders = make(map[string][]byte)
	sdIDStream = make(chan string, 1024)
	RetMsg = make(map[string][]byte)
}
