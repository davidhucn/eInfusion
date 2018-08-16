package tqueue

type TrsfusionCtrl struct {
	rcvID string
	detID string
	cmd   []byte // 命令 、 参数数据
	// ipAdd string
}

var HTTPReqStream chan TrsfusionCtrl

func init() {
	HTTPReqStream = make(chan TrsfusionCtrl, 1024)
}
