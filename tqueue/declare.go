package tqueue

// ReqCmd :指令类型
type ReqCmd struct {
	TargetID chan string
	CmdType  chan uint8  // 指令类型(代码)
	Args     chan string // 相关参数 (例如：ip、port)
}

// sendOrder :全局指令map,MAP索引为时间戳
var sendOrders chan map[string][]byte

// var sOrders chan []byte

func init() {
	sendOrders = make(chan map[string][]byte, 1024)

	// sOrders = make(chan []byte, 1024)

}
