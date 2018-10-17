package datahub

import (
	cm "eInfusion/comm"
	"sync"
)

type orderType struct {
	WebSocket int
	TCP       int
	MQTT      int
}

var OrderType orderType

type OrdersQueue struct {
	Queue map[int]chan *cm.Cmd
	sync.Mutex
}

func init() {
	OrderType.TCP = 0
	OrderType.WebSocket = 1
	OrderType.MQTT = 2
}

func NewOrderQueue() *OrdersQueue {
	return &OrdersQueue{
		Queue: make(map[int]chan *cm.Cmd, 1024),
	}
}

func (oq *OrdersQueue) AddQueue(rOdType int, rOd *cm.Cmd) {
	oq.Lock()
	oq.Queue[rOdType] <- rOd
	oq.Unlock()
}

func (oq *OrdersQueue) GetQueueOrder(rOdType int, rCmdID string) {
	for od, ok := range oq[rOdType] {

	}
}
