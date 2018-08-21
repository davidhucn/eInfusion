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

// cvtDec :用于简便,转换成decimal
func cvtDec(rData uint8) uint8 {
	return cm.ConvertBasStrToUint(10, cm.ConvertBasNumberToStr(16, rData))
}

// GetClisCmd :确定websocket指令
func GetClisCmd(rReq reqMsg) bool {
	// TODO: 处理客户端发来的指令
	switch cm.ConvertBasStrToUint(10, rReq.CmdType) {
	case cvtDec(ep.TrsCmdType.AddDetect):

		return true
	case cvtDec(ep.TrsCmdType.DelDetect):

		return true
	case cvtDec(ep.TrsCmdType.SetRcvNetCfg):

		return true
	}

	return false
	// }
	// return false
}
