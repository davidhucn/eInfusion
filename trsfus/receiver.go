package trsfus

import "time"

// Receiver ：检测器对象
type Receiver struct {
	ID       string
	nativeIP string
	// nativePort    string
	detectAmount  int
	reconnectTime time.Duration
	targetIP      string
	targetPort    string
}

// NewReceiver ：新建一个接收器对象
func NewReceiver(id string, detectAmount int, nativeIP string, reconnectTime time.Duration,
	targetIP string, targetPort string) *Receiver {
	return &Receiver{
		ID:            id,
		detectAmount:  detectAmount,
		nativeIP:      nativeIP,
		reconnectTime: reconnectTime,
		targetIP:      targetIP,
		targetPort:    targetPort,
	}
}

// GetStatus :获取检测器状态
// func (r *Receiver) GetStatus() []byte {

// }
