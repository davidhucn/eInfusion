package tqueue

// ReqCmd :指令类型
type ReqCmd struct {
	TargetID chan string
	CmdType  chan uint8  // 指令类型(代码)
	Args     chan string // 相关参数 (例如：ip、port)
}

// sendOrder :全局指令map
var sendOrders map[string][]byte

func init() {
	sendOrders = make(map[string][]byte)
}
