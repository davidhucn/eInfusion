package tcpOperate

import (
	"eInfusion/comm"

	"github.com/firstrow/tcp_server"
)

func StartTcpTest() {
	server := tcp_server.New(":7778")
	comm.ShowScreen("["+comm.GetCurrentTime()+"]", "Transfusion数据平台开始运行...")

	server.OnNewClient(func(c *tcp_server.Client) {

		c.Send("3-11,hi")
		comm.ShowScreen("client address: ", c.Conn().RemoteAddr())
	})

	server.OnNewMessage(func(c *tcp_server.Client,message) {
		//		comm.ShowScreen("test")
		//		var bb []byte
		comm.ShowScreen(string(message))
	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		// connection with client lost
		comm.ShowScreen(c.Conn().RemoteAddr(), "lost...")
	})

	server.Listen()
}
