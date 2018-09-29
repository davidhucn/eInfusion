package thub

import "sync"

// WebClients :全局ws连接对象
type WebClients struct {
	Connections sync.Map  // [string] *websocket(gorilla)
	Orders      chan *Cmd // 发送信息的缓冲 channel
}

// NewDevices :新建Alarms对象
func NewDevices() *Devices {
	return &Devices{
		Orders: make(chan *Cmd, 1024),
	}
}

// NewWebClients :新建WebClient对象
func NewWebClients() *WebClients {
	return &WebClients{
		Orders: make(chan *Cmd, 1024),
	}
}

// Rcvs :对象
var Rcvs = NewDevices()
