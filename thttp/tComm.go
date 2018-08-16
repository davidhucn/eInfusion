package thttp

import (
	cm "eInfusion/comm"
	ep "eInfusion/protocol"
)

type reqMsg struct {
	ID      string `json:"ID"`
	CmdType string `json:"CmdType"` //指令类型(代码)
	Args    string `json:"Args"`    //相关参数 (例如：ip、port)
	// Action string `json:"-"`
}

var clisMsg []reqMsg

// cvtDec :用于当下转换
func cvtDec(rData uint8) uint8 {
	return cm.ConvertBasStrToUint(10, cm.ConvertBasNumberToStr(16, rData))
}

// verifyReqWS :判断websocket数据执行标志是否为真
func va(rReq reqMsg) bool {
	// TODO: 处理客户端发来的指令
	switch cm.ConvertBasStrToUint(10, rReq.CmdType) {
	case cvtDec(ep.TrsCmdType.AddDetect):

	case cvtDec(ep.TrsCmdType.DelDetect):
	}

	return false
	// }
	// return false
}
