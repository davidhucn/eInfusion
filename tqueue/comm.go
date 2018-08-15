package tqueue

type sTrsfusionCtrl struct {
	rcvID string
	detID string
	// 命令 、 参数数据
	cmd []byte
}

var G_TsReq sTrsfusionCtrl
