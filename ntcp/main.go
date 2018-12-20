package ntcp

import (
	cm "eInfusion/comm"
	"time"
)

// StartTCPService :启动TCP服务
func StartTCPService() {
	h := NewTCPHeader(3, []byte("0x66"), 1)
	sv := NewTCPServer(":9909", 10*time.Minute, 6*time.Hour, h)

	sv.WhenNewClientConnected(func(c *Client) {
		if c.VerifyLegal() {
			// 遍历待发队列，发送
			for i, od := range sv.waitQueue {
				if od.ID == cm.GetPureIPAddr(c.conn) {
					c.SendData(od.Data)
					// 册除待发队列相应项
					sv.waitQueue = append(sv.waitQueue[:i], sv.waitQueue[i+1])
				}
			}
		} else {
			// 非法客户端 TODO:

		}
	})

	sv.WhenNewDataReceived(func(c *Client, p []byte) {
		// TODO: 解析，落到到具体业务

	})

	sv.WhenClientConnectionClosed(func(c *Client, err error) {

	})
	sv.Listen()
}
