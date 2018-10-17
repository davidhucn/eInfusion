package datahub

import (
	cm "eInfusion/comm"
	"sync"
)

// type orderType string

type OrderType struct {
	WebSocket string
	TCP       string
	MQTT      string
}

type OrdersQueue struct {
	Queue map[string]chan *cm.Cmd
	sync.Mutex
}

func NewOrderQueue() *OrdersQueue {
	return &OrdersQueue{
		// TODO:完成新建对象
		Queue : make(map[string])
	}
}
