package thttp

type reqMsg struct {
	ID     string `json:"ID"`
	Action string `json:"Action"` //执行标志
	// Action string `json:"-"`
}

var cliMsg []reqMsg

// verifyReqWS :判断websocket数据执行标志是否为真
func verifyReqAction(rReq reqMsg) bool {
	if rReq.Action == "1" {
		return true
	}
	return false
}
