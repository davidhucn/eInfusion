package tcp

import (
	"eInfusion/comm"

	"github.com/firstrow/tcp_server"
)

func StartTcpServer() {
	server := tcp_server.New(":7778")
	comm.ShowScreen("["+comm.GetCurrentTime()+"]", "Transfusion数据平台开始运行...")

	server.OnNewClient(func(c *tcp_server.Client) {

		c.Send("hi,david")
		comm.ShowScreen("local address ", c.Conn().LocalAddr())
	})
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		// new message received
	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		// connection with client lost
		comm.ShowScreen("new clienk lost")
	})

	server.Listen()
}
