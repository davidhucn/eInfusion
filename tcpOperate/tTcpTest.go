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
		comm.ShowScreen("client address ", c.Conn().RemoteAddr())
	})

	server.OnNewMessage(func(c *tcp_server.Client, m string) {
		//		comm.ShowScreen("test")
		var bb []byte
		n, _ := c.Conn().Read(bb)
		comm.ShowScreen(n)
		//		comm.ShowScreen("receive numer: ", n)
		//		comm.ShowScreen(bb[0])
		//		comm.ShowScreen("from client: ", m)
	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		// connection with client lost
		comm.ShowScreen(c.Conn().RemoteAddr(), "lost...")
	})

	server.Listen()
}
