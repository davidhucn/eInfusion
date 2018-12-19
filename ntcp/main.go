package ntcp

import (
	"fmt"
	"time"
)

// StartTCPService :启动TCP服务
func StartTCPService() {
	h := NewTCPHeader(3, []byte("0x66"), 1)
	ser := NewTCPServer(":9909", 10*time.Minute, h)

	ser.WhenNewClientConnected(func(c *Client) {
		if c.SendData([]byte("nice to meet you!")) == nil {
			fmt.Println("try")
		}
	})

	ser.WhenNewDataReceived(func(c *Client, p []byte) {

	})

	ser.WhenClientConnectionClosed(func(c *Client, err error) {

	})
	ser.Listen()
}
