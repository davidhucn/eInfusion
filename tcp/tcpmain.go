package tcp

//@Author:david
//@Date:2018/09/28
//@Purpose: Tcp Socket Server(Epoll模型）

import (
	cm "eInfusion/comm"
	dh "eInfusion/datahub"
	logs "eInfusion/tlogs"
	tsc "eInfusion/trsfscomm"
	"net"
	"os"
	"strconv"
	"time"
)

// Broadcast :对所有连接发送广播
func (ts *TServer) Broadcast(rOrder *cm.Cmd) {
	for _, c := range ts.Connections {
		rOrder.CmdID = dh.NewTCPOrderID(rOrder.CmdID, cm.GetPureIPAddr(c))
		dh.AddToTCPOrderQueue(rOrder)
	}
}

// SendOrderAndMsg :发送给设备单独的某条命令、数据，并返回指定信息给WEB前端
func (ts *TServer) SendOrderAndMsg(rOrder *cm.Cmd, rWebMsg string) error {
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
	// 发送成功记录日志
	// FIXME:抽象化log对象
	logs.LogMain.Info(TCPMsg.SendSuccess, "发送至=> IP："+cm.GetPureIPAddr(ts.Connections[connID]))
	return nil
}

// LoopingTCPOrder :循环发送设备对象内的指令序列
func (ts *TServer) LoopingTCPOrder() {
	// 循环清除超过指定时间周期的【待发列表】
	dm, _ := time.ParseDuration(cm.ConvertIntToStr(ts.ExpireTimeByMinutes) + "m")
	// dm := ts.ExpireTimeByMinutes * time.Minute
	go func() {
		for v := range ts.WaitOrders {
			//判断指令是否在生存期内
			if time.Now().Sub(v.CreateTime) < dm {
				// 由于还在生存期内，继续发回待发队列
				ts.WaitOrders <- v
			} else {
				// 超过生存周期，错误信息回写到前端
				cm.Msg("out of time:", v.SendData.CmdID)
				dh.SendMsgToWeb(cm.NewOrder(v.SendData.CmdID, []byte(TCPMsg.SendFailureForLongTime)))
			}
		}
	}()

	// 循环发送指令，基于datahub指令channel,并存放到相应连接的消息内容中(兼容其它模块发来的数据和指令)
	go func() {
		for od := range dh.GetTCPOrderQueue() {
			if cm.CkErr("", ts.SendOrderAndMsg(od, TCPMsg.SendSuccess)) {
				// 如果发送失败，放入待发送列表
				var wod WaitOrder
				wod.CreateTime = time.Now()
				wod.SendData = od
				ts.WaitOrders <- wod
			}
		}
	}()
}

// setReadTimeout :设定TCP连接接收数据时间(指定连接或不指定)
func (ts *TServer) setReadTimeout(rConnID string, t time.Duration) {
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
func (ts *TServer) madeConn(c *net.TCPConn) {
	connID := cm.GetPureIPAddr(c)
	ts.Lock()
	// 加入到连接列表中
	ts.Connections[connID] = c
	ts.Unlock()
	logs.LogMain.Info("IP:", connID, "上线")
	cm.SepLi(20, "-")
	// v.CreateTime.Add(5*time.Minute)
	// 遍历待发列表，如果符合连接ID即可发送
	go func() {
		for v := range ts.WaitOrders {
			if dh.DecodeToTCPConnID(v.SendData.CmdID) == connID {
				dh.AddToTCPOrderQueue(v.SendData)
			} else {
				// 如果不符合连接ID，放回待发列表
				ts.WaitOrders <- v
			}
		}
	}()
	// ****定时处理(心跳等)
	//	go loopingCall(conn)
}

// lostConn :连接断开
func (ts *TServer) lostConn(c *net.TCPConn) {
	connID := cm.GetPureIPAddr(c)
	ts.Lock()
	delete(ts.Connections, connID)
	ts.Unlock()
	logs.LogMain.Info("IP:", connID, "下线")
	cm.SepLi(30, "*")
}

//echo server Goroutine
func receiveData(c *net.TCPConn) {

	{ // TODO:测试包
		// SendData(c, tsc.CmdGetRcvStatus(cm.ConvertPerTwoOxCharOfStrToBytes("A0000000")))
		// time.Sleep(10 * time.Millisecond)
		// dtID := cm.ConvertStrToBytesByPerTwoChar("B0000000")
		// rvID := cm.ConvertStrToBytesByPerTwoChar("A0000000")

		// // bIP := cm.ConvertStrIPToBytes("192.168.121.12")
		// // bPort := cm.ConvertDecToBytes(7778)
		// orders := tsc.CmdOperateDetect(tsc.G_TsCmd.DelDetect, rvID, 1, dtID)
		// // orders := tsc.CmtsetRcvReconTime(rvID, cm.ConvertDecToBytes(900))
		// // orders := tsc.CmtsetRcvCfg(rvID, bIP, bPort)
		// SendData(c, orders)
	}
	defer c.Close()
	for {
		c.SetReadDeadline(time.Now().Add(5 * time.Minute))
		//	指定接收数据包头的帧长
		recDataHeader := make([]byte, tsc.TrsDefin.HeaderLength)
		_, err := c.Read(recDataHeader)
		if err != nil {
			break
		}
		// 数据包数据内容长度记录变量
		var intPckContentLength int
		// 判断包头是否正确，如果正确，获取长度
		if !DecodeHeader(recDataHeader, &intPckContentLength) {
			logs.LogMain.Info(cm.GetPureIPAddr(c), TCPMsg.HeaderDataError)
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
				DecodeRcvData(r, cm.GetPureIPAddr(c))
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

// RunTCPService :启动TCP服务
func RunTCPService(ts *TServer, port int) {
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
	// 循环处理TCP对象指令
	ts.LoopingTCPOrder()

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

// 以下为channel应用参考实例
// if cm.CkErr("", ts.SendOrderAndMsg(od, TCPMsg.SendSuccess)) {
//发送不成功，则延迟发送
// cTicker := time.NewTicker(12 * time.Second) // 定时
// lastCk := time.After(3 * time.Minute)       // 延时
// defer cTicker.Stop()
// for i := 0; i < 3; i++ {
// 	select {
// 	case <-cTicker.C:
// 		if !cm.CkErr("", ts.SendOrderAndMsg(od, TCPMsg.SendSuccess)) {
// 			// continue
// 			break
// 		}
// 	}
// }
// select {
// case <-lastCk:
// 	if !cm.CkErr("", ts.SendOrderAndMsg(od, TCPMsg.SendSuccess)) {
// 		// continue
// 		break
// 	}
// }
// dh.SendMsgToWeb(cm.NewOrder(od.CmdID, []byte(TCPMsg.SendFailureForLongTime)))

// FIXME:这里有问题,需重新注销函数
// dh.UnregisterReqOrdersUnion()
// }
