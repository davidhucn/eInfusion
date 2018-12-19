package ntcp

import "time"

// StartTCPService :启动TCP服务
func StartTCPService() {
	h := NewTCPHeader(3, []byte("0x66"), 1)
	ser := NewTCPServer(":9909", 10*time.Minute, h)

	ser.WhenNewClientConnected(func(c *Client) {
		c.SendData([]byte("nice to meet you!"))
	})
	ser.Listen()
}
