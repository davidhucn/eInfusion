package thub

import (
	"sync"
)

// Cmd : 指令对象
type Cmd struct {
	Cmd   []byte
	CmdID string
}

// WebClis :web连接对象
type WebClis struct {
	Connections sync.Map // [string] *gorilla/websocke
	Orders      chan Cmd //指令
}

// Alarms :设备对象，现主要对接tcp
type Alarms struct {
	Connections sync.Map //[string] *net.connection
	Orders      chan Cmd // 指令
}

// NewAlarms :新建得一个tcp 设备对象
func NewAlarms() *Alarms {
	return &Alarms{
		Orders: make(chan Cmd, 1024),
	}
}

// NewWebClis :新建得一个websocket设备对象
func NewWebClis() *WebClis {
	return &WebClis{
		Orders: make(chan Cmd, 1024),
	}
}
