package tcpOperate

//@Author:david
//@Date:2018/06/29
//@Purpose: Tcp Socket Server(Epoll模型）
import (
	"eInfusion/comm"
	"eInfusion/logs"
	//	ep "eInfusion/protocol"
	"net"
	"time"
)

//func mkClisConn(key string, conn *net.TCPConn) {
//	connMkMutex.Lock()
//	defer connMkMutex.Unlock()
//	clisConnMap[key] = conn
//}

///*
//   删除socket conn 映射
//*/
//func delClisConn(key string) {
//	connDelMutex.Lock()
//	defer connDelMutex.Unlock()
//	delete(clisConnMap, key)
//}

//echo server Goroutine
func EchoFunc(c EndPointer) {
	//	FIXME: 修改这里
	defer c.Conn.Close()
	buf := make([]byte, 1024)
	for {
		_, err := c.Conn.Read(buf)
		if err != nil {
			//println("Error reading:", err.Error())
			return
		}
		//send reply
		_, err = c.Conn.Write(buf)
		if err != nil {
			//println("Error send reply:", err.Error())
			return
		}
	}
}

//initial listener and run
func StartTcpServer() {
	listener, err := net.Listen("tcp", ":"+c_TcpServer_Port)
	defer listener.Close()
	if err != nil {
		logs.LogMain.Critical("监听TCP出错", err)
		panic(err)
	}
	logs.LogMain.Info(c_Msg_Info_ServerStart + "（" + comm.GetCurrentDate() + "）")
	comm.SepLi(60)
	comm.Msg("TCP Port:" + c_TcpServer_Port)
	//并发，在线数量
	var intConcurrentNum int = 0

	//	var c_stream chan Clienter
	conn_stream := make(chan net.Conn)

	intConnCounter_stream := make(chan int)
	//////////////////////////////////////////////////////////////
	go func() {
		for intConnTemp := range intConnCounter_stream {
			intConcurrentNum += intConnTemp
		}
	}()
	go func() {
		for _ = range time.Tick(1e8) {
			comm.Msg("cur conn num: %f\n", intConcurrentNum)
		}
	}()
	////////////////////process per connection///////////////////////////////
	for i := 0; i < c_MaxConnectionAmount; i++ {
		go func() {
			for conn := range conn_stream {
				intConnCounter_stream <- 1
				EchoFunc(conn)
				intConnCounter_stream <- -1

			}
		}()
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
			logs.LogMain.Error(c_Msg_Err_AcceptConnection)
		}
		conn_stream <- conn
	}
}
