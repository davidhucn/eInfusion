package tcp

//@Author:david
//@Date:2018/09/28
//@Purpose: Tcp Socket Server(Epoll模型）

import (
	cm "eInfusion/comm"
	dh "eInfusion/datahub"
	ep "eInfusion/protocol"
	logs "eInfusion/tlogs"
	"net"
	"os"
	"strconv"
	"time"
)

// Broadcast :对所有连接发送广播
func (ts *TCPServer) Broadcast(rOrder cm.Cmd) {
	for _, c := range ts.Connections {
		_, err := c.Write(rOrder.Cmd)
		if cm.CkErr(TCPMsg.SendError, err) {
			continue
		}
	}
}

// SendOrder :发送给设备单独的某条命令、数据
func (ts *TCPServer) SendOrder(rOrder *cm.Cmd, rWebMsg string) error {
	// 获取tcp连接id
	connID := dh.DecodeToTCPConnID(rOrder.CmdID)
	if _, ok := ts.Connections[connID]; !ok {
		return cm.ConvertStrToErr(TCPMsg.CanNotFindConnection)
	}
	time.Sleep(15 * time.Millisecond)
	_, err := ts.Connections[connID].Write(rOrder.Cmd)
	if cm.CkErr(TCPMsg.SendError, err) {
		return cm.ConvertStrToErr(TCPMsg.SendError)
	}
	if rWebMsg != "" {
		od := cm.NewOrder(rOrder.CmdID, []byte(rWebMsg))
		dh.SendMsgToWeb(od)
	}
	// 	// TODO:抽象化log对象
	logs.LogMain.Info("=> IP："+ts.Connections[connID].RemoteAddr().String(), TCPMsg.SendSuccess)
	return nil
}

// LoopingDoOrders :循环发送设备对象内的指令序列
func (ts *TCPServer) LoopingDoOrders() {
	// 循环获取指令
	go func() {
		for dh.DeviceTCPOrderQueue != nil {
			select {
			case od := <-dh.DeviceTCPOrderQueue:
				ts.Orders <- od
			}
		}
	}()
	for od := range ts.Orders {
		if cm.CkErr(TCPMsg.SendError, ts.SendOrder(od, TCPMsg.SendSuccess)) {
			//发送不成功，则延迟发送
			cTicker := time.NewTicker(12 * time.Second) // 定时
			lastCk := time.After(1 * time.Minute)       // 延时
			defer cTicker.Stop()
			for i := 0; i < 3; i++ {
				select {
				case <-cTicker.C:
					if !cm.CkErr("", ts.SendOrder(od, TCPMsg.SendSuccess)) {
						continue
					}
				}
			}
			select {
			case <-lastCk:
				if !cm.CkErr("", ts.SendOrder(od, TCPMsg.SendSuccess)) {
					continue
				}
			}
			dh.SendMsgToWeb(cm.NewOrder(od.CmdID, []byte(TCPMsg.SendFailureForLongTime)))
		}
	}
}

// RetrieveTCPOrdersFromDataHub : 循环检测datahub包内DeviceTCPOrders对象，加入到发送队列
func (ts *TCPServer) RetrieveTCPOrdersFromDataHub() {
	for dh.DeviceTCPOrderQueue != nil {
		select {
		case od := <-dh.DeviceTCPOrderQueue:
			ts.Orders <- od
		}
	}
}

// setReadTimeout:设置读数据超时xtswa
func setReadTimeout(conn *net.TCPConn, t time.Duration) {
	conn.SetReadDeadline(time.Now().Add(t))
}

// setReadTimeout :设定TCP连接接收数据时间
func (ts *TCPServer) setReadTimeout(rConnID string, t time.Duration) {
	if rConnID == "" {
		for _, c := range ts.Connections {
			c.SetReadDeadline(time.Now().Add(t))
		}
		return
	}
	if c, ok := ts.Connections[rConnID]; ok {
		c.SetReadDeadline(time.Now().Add(t))
	}
}

//madeConn :连接初始处理(ed)
func (ts *TCPServer) madeConn(c *net.TCPConn) {
	//初始化连接这个函数被调用
	connID := cm.GetPureIPAddr(c.RemoteAddr().String())
	ts.Lock()
	ts.Connections[connID] = c
	ts.Unlock()
	logs.LogMain.Info("IP:", connID, "上线")
	cm.SepLi(20, "-")
	// ****定时处理(心跳等)
	//	go loopingCall(conn)
}

// lostConn :连接断开
func (ts *TCPServer) lostConn(c *net.TCPConn) {
	connID := cm.GetPureIPAddr(c.RemoteAddr().String())
	ts.Lock()
	delete(ts.Connections, connID)
	ts.Unlock()
	logs.LogMain.Info("IP:", connID, "下线")
	cm.SepLi(30, "*")
}

//echo server Goroutine
func receiveData(c *net.TCPConn) {

	{ // TODO:测试包
		// SendData(c, ep.CmdGetRcvStatus(cm.ConvertPerTwoOxCharOfStrToBytes("A0000000")))
		// time.Sleep(10 * time.Millisecond)
		// dtID := cm.ConvertStrToBytesByPerTwoChar("B0000000")
		// rvID := cm.ConvertStrToBytesByPerTwoChar("A0000000")

		// // bIP := cm.ConvertStrIPToBytes("192.168.121.12")
		// // bPort := cm.ConvertDecToBytes(7778)
		// orders := ep.CmdOperateDetect(ep.G_TsCmd.DelDetect, rvID, 1, dtID)
		// // orders := ep.CmtsetRcvReconTime(rvID, cm.ConvertDecToBytes(900))
		// // orders := ep.CmtsetRcvCfg(rvID, bIP, bPort)
		// SendData(c, orders)
	}
	defer c.Close()
	for {
		setReadTimeout(c, 5*time.Minute)
		//	指定接收数据包头的帧长
		recDataHeader := make([]byte, ep.TrsDefin.HeaderLength)
		_, err := c.Read(recDataHeader)
		if err != nil {
			break
		}
		// 数据包数据内容长度记录变量
		var intPckContentLength int
		// 判断包头是否正确，如果正确，获取长度
		if !ep.DecodeHeader(recDataHeader, &intPckContentLength) {
			logs.LogMain.Info(cm.GetPureIPAddr(c.RemoteAddr().String()), TCPMsg.HeaderDataError)
			//	如果包头不正确，断开连接
			break
			//	continue 退出本次
		}
		// 如果包头接收正确
		r := make([]byte, intPckContentLength)
		n, er := c.Read(r)
		if !cm.CkErr(TCPMsg.ReceiveError, er) {
			//	实际长度和包头内规定长度不一致
			if n == intPckContentLength {
				// 处理报文数据内容
				ep.DecodeRcvData(r, cm.GetPureIPAddr(c.RemoteAddr().String()))
			} else {
				/*如果长度不一致，退出*/
				//	continue
				/*断开连接*/
				break
			}
		} else {
			continue /*如果接收出错，退出*/
			//	break /*断开连接*/
		}
	}
}

// RunTCPJob :Device对象启动TCP任务
func RunTCPJob(ts *TCPServer, port int) {
	host := ":" + strconv.Itoa(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if cm.CkErr(TCPMsg.SourceError, err) {
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	defer listener.Close()
	if cm.CkErr(TCPMsg.SourceError, err) {
		// os.Exit(1)
		return
	}
	cm.SepLi(60, "")
	logs.LogMain.Info(TCPMsg.StartServiceMsg + "（" + cm.GetCurrentDate() + "）")
	cm.Msg("Transfusion System Server Port", host)
	cm.SepLi(60, "")
	// 循环发送指令
	{
		// go ts.LoopingSendOrders()
		// go ts.RetrieveTCPOrdersFromDataHub()
		go ts.LoopingDoOrders()
	}
	// var connStream chan *net.TCPConn
	connStream := make(chan *net.TCPConn)
	//打开N个Goroutine等待连接，Epoll模式
	for i := 0; i < ts.MaxTCPConnect; i++ {
		go func() {
			for cs := range connStream {
				ts.madeConn(cs)
				//	接收数据
				receiveData(cs)
				ts.lostConn(cs)
			}
		}()
	}
	for {
		lc, err := listener.AcceptTCP()
		if cm.CkErr(TCPMsg.ReceiveError, err) {
			continue
		}
		connStream <- lc
	}
}
