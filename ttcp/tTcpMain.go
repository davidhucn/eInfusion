package ttcp

// //@Author:david
// //@Date:2018/06/29
// //@Purpose: Tcp Socket Server(Epoll模型）

// import (
// 	"eInfusion/comm"
// 	ep "eInfusion/protocol"
// 	logs "eInfusion/tlogs"
// 	"net"
// 	"os"
// 	"strconv"
// 	"time"
// )

// func mkClisConn(key string, conn *net.TCPConn) {
// 	connMkMutex.Lock()
// 	defer connMkMutex.Unlock()
// 	ClisConnMap[key] = conn
// 	//TODO:获取等待指令，一联线就通知消息平台，由消息平台发送指令

// }

// //  删除socket conn 映射
// func delClisConn(key string) {
// 	connDelMutex.Lock()
// 	defer connDelMutex.Unlock()
// 	delete(ClisConnMap, key)
// }

// //  初始化socket conn 映射
// func initClisConnMap() {
// 	ClisConnMap = make(map[string]*net.TCPConn)
// }

// // SendData :发送命令和数据
// func SendData(conn *net.TCPConn, data []byte) (err error) {
// 	ip := comm.GetPureIPAddr(conn)
// 	//考虑网络延迟、断网问题，另外发送两个数据须间隔15毫秒(millionseconds)
// 	time.Sleep(15 * time.Millisecond)
// 	_, err = conn.Write(data)
// 	if err == nil {
// 		logs.LogMain.Info("=>"+ip, "完成数据发送")
// 		// 如果发送没有错误，即表示成功
// 	}
// 	return
// }

// //Broadcast :广播数据
// func Broadcast(data []byte) {
// 	for _, conn := range ClisConnMap {
// 		SendData(conn, data)
// 	}
// }

// //   定时处理&延时处理
// //func loopingCall(conn *net.TCPConn) {
// //	pingTicker := time.NewTicker(30 * time.Second) // 定时
// //	testAfter := time.After(5 * time.Second)       // 延时

// //	for {
// //		select {
// //		case <-pingTicker.C:
// //			//发送心跳
// //			_, err := SendData(conn, []byte("PING"))
// //			if err != nil {
// //				pingTicker.Stop()
// //				return
// //			}
// //		case <-testAfter:
// //			//	doLog("testAfter:")
// //			//TODO:日志记录
// //		}
// //	}
// //}

// //连接初始处理(ed)
// func madeConn(conn *net.TCPConn) {
// 	//初始化连接这个函数被调用
// 	mkClisConn(comm.GetPureIPAddr(conn), conn)
// 	logs.LogMain.Info("IP:", comm.GetPureIPAddr(conn), "上线")
// 	comm.SepLi(20, "-")
// 	// ****定时处理(心跳等)
// 	//	go loopingCall(conn)
// }

// // 连接断开
// func lostConn(conn *net.TCPConn) {
// 	//连接断开这个函数被调用

// 	ip := comm.GetPureIPAddr(conn)
// 	delClisConn(ip) // 删除关闭的连接对应的clisMap项
// 	logs.LogMain.Info("IP:", ip, "下线")
// 	comm.SepLi(30, "*")
// }

// //echo server Goroutine
// func receiveData(c *net.TCPConn) {

// 	{ // TODO:测试包
// 		// SendData(c, ep.CmdGetRcvStatus(comm.ConvertPerTwoOxCharOfStrToBytes("A0000000")))
// 		// time.Sleep(10 * time.Millisecond)
// 		// dtID := comm.ConvertStrToBytesByPerTwoChar("B0000000")
// 		// rvID := comm.ConvertStrToBytesByPerTwoChar("A0000000")

// 		// // bIP := comm.ConvertStrIPToBytes("192.168.121.12")
// 		// // bPort := comm.ConvertDecToBytes(7778)
// 		// orders := ep.CmdOperateDetect(ep.G_TsCmd.DelDetect, rvID, 1, dtID)
// 		// // orders := ep.CmdSetRcvReconTime(rvID, comm.ConvertDecToBytes(900))
// 		// // orders := ep.CmdSetRcvCfg(rvID, bIP, bPort)
// 		// SendData(c, orders)
// 	}
// 	defer c.Close()
// 	for {
// 		setReadTimeout(c, 5*time.Minute)
// 		//	指定接收数据包头的帧长
// 		recDataHeader := make([]byte, ep.TrsDefin.HeaderLength)
// 		_, err := c.Read(recDataHeader)
// 		if err != nil {
// 			break
// 		}
// 		// 数据包数据内容长度记录变量
// 		var intPckContentLength int
// 		// 判断包头是否正确，如果正确，获取长度
// 		if !ep.DecodeHeader(recDataHeader, &intPckContentLength) {
// 			logs.LogMain.Info(comm.GetPureIPAddr(c), "的数据包头异常,强制下线!")
// 			//	如果包头不正确，断开连接
// 			break
// 			//	continue 退出本次
// 		}
// 		// 如果包头接收正确
// 		r := make([]byte, intPckContentLength)
// 		n, er := c.Read(r)
// 		if !comm.CkErr("接收报文出错", er) {
// 			//	实际长度和包头内规定长度不一致
// 			if n == intPckContentLength {
// 				// 处理报文数据内容
// 				ep.DecodeRcvData(r, comm.GetPureIPAddr(c))
// 			} else {
// 				/*如果长度不一致，退出*/
// 				//	continue
// 				/*断开连接*/
// 				break
// 			}
// 		} else {
// 			continue /*如果接收出错，退出*/
// 			//	break /*断开连接*/
// 		}
// 	}
// }

// // StartTCPServer :TCP启动服务器入口
// func StartTCPServer(port int) {
// 	host := ":" + strconv.Itoa(port)
// 	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
// 	if comm.CkErr("TCP资源错误", err) {
// 		os.Exit(1)
// 	}
// 	listener, err := net.ListenTCP("tcp", tcpAddr)
// 	defer listener.Close()
// 	if comm.CkErr("TCP资源错误", err) {
// 		os.Exit(1)
// 	}
// 	comm.SepLi(60, "")
// 	logs.LogMain.Info(c_Msg_Info_ServerStart + "（" + comm.GetCurrentDate() + "）")
// 	comm.Msg("Transfusion System Server Port", host)
// 	comm.SepLi(60, "")
// 	connStream := make(chan *net.TCPConn)
// 	initClisConnMap()
// 	//打开N个Goroutine等待连接，Epoll模式
// 	for i := 0; i < c_MaxConnectionAmount; i++ {
// 		go func() {
// 			for cs := range connStream {
// 				madeConn(cs)
// 				//	接收数据
// 				receiveData(cs)
// 				lostConn(cs)
// 			}
// 		}()
// 	}
// 	for {
// 		lc, err := listener.AcceptTCP()
// 		if err != nil {
// 			logs.LogMain.Error(c_Msg_Err_AcceptConnection)
// 			continue
// 		}
// 		connStream <- lc
// 	}
// }

// //   设置读数据超时
// func setReadTimeout(conn *net.TCPConn, t time.Duration) {
// 	conn.SetReadDeadline(time.Now().Add(t))
// }
