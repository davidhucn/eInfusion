package thttp

import (
	cm "eInfusion/comm"
	ep "eInfusion/protocol"
)

type reqData struct {
	ID      string `json:"ID"`
	CmdType string `json:"CmdType"` //指令类型(代码)
	Args    string `json:"Args"`    //相关参数 (例如：ip、port)
	// Action string `json:"-"`
}

var clisMsg []reqData

// GetClisCmd :确定websocket指令
func GetClisCmd(rReq reqData) {
	// TODO: 处理客户端发来的指令
	switch cm.ConvertBasStrToUint(10, rReq.CmdType) {
	case cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.AddDetect):

	case cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.DelDetect):

	case cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.SetRcvNetCfg):

	}

}
