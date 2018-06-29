package tcpOperate

//@Author:david
//@Date:2018/06/29
//@Purpose: Tcp Socket Server(Epoll模型）
import (
	"eInfusion/comm"
	"eInfusion/logs"
	//	ep "eInfusion/protocol"
	"net"
	"strconv"
	"time"
)

<<<<<<< HEAD
func mkClisConn(key string, conn *net.TCPConn) {
	connMkMutex.Lock()
	defer connMkMutex.Unlock()
	clisConnMap[key] = conn
}

//  删除socket conn 映射
func delClisConn(key string) {
	connDelMutex.Lock()
	defer connDelMutex.Unlock()
	delete(clisConnMap, key)
}

//  初始化socket conn 映射
func initClisConnMap() {
	clisConnMap = make(map[string]*net.TCPConn)
}

// 连接断开
func connectionLost(conn *net.TCPConn) {
	//连接断开这个函数被调用
	addr := conn.RemoteAddr().String()
	ip := strings.Split(addr, ":")[0]

	delClisConn(ip) // 删除关闭的连接对应的clisMap项
	//TODO:	记录日志
	//	doLog("connectionLost:", addr)
}

//   发送数据
func sendData(conn *net.TCPConn, data []byte) (n int, err error) {
	addr := conn.RemoteAddr().String()
	n, err = conn.Write(data)
	if err == nil {
		//TODO:记录日志
	}
	return
}

func ckError(err error, title string, exit bool) {
	if err != nil {
		if exit == true {
			logs.LogMain.Critical(title, err.Error())
			os.Exit(1)
		} else {
			logs.LogMain.Error(title, err.Error())
=======
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
>>>>>>> 7d5589a2bc1ca06c570d39bf6c057c1a6ce15959
		}
	}
}

//   定时处理&延时处理
func loopingCall(conn *net.TCPConn) {
	pingTicker := time.NewTicker(30 * time.Second) // 定时
	testAfter := time.After(5 * time.Second)       // 延时

	for {
		select {
		case <-pingTicker.C:
			//发送心跳
			_, err := sendData(conn, []byte("PING"))
			if err != nil {
				pingTicker.Stop()
				return
			}
		case <-testAfter:
			//	doLog("testAfter:")
			//TODO:日志记录
		}
	}
}

//连接初始处理(ed)
func connectionMade(conn *net.TCPConn) {
	//初始化连接这个函数被调用

	// ****建立conn映射
	addr := conn.RemoteAddr().String()
	ip := strings.Split(addr, ":")[0]
	mkClisConn(ip, conn)

	doLog("connectionMade:", addr)

	// ****定时处理(心跳等)
	//	go loopingCall(conn)
}

//initial listener and run
func StartTcpServer(port int) {
	host := ":" + strconv.Itoa(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	ckError(err, "TCP资源错误", true)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	defer listener.Close()
	ckError(err, "TCP监听错误", true)
	logs.LogMain.Info(c_Msg_Info_ServerStart + "（" + comm.GetCurrentDate() + "）")
	comm.SepLi(60)
<<<<<<< HEAD
	comm.Msg("TCP Server Info:" + tcpAddr.IP.String() + host)

	conStream := make(chan net.TCPConn)
	initClisConnMap()
	//打开N个Goroutine等待连接，Epoll模式
	for i := 0; i < c_MaxConnectionAmount; i++ {
		go func() {
			for cs := range conStream {
				connectionMade(cs)
				//TODO:[优先]接收数据
				//EchoFunc(cs)
=======
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

>>>>>>> 7d5589a2bc1ca06c570d39bf6c057c1a6ce15959
			}
		}()
	}

	for {
		lc, err := listener.AcceptTCP()
		if err != nil {
			logs.LogMain.Error(c_Msg_Err_AcceptConnection)
			continue
		}
		conStream <- lc
	}
}
