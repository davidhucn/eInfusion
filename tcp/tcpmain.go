package tcp

//@Author:david
//@Date:2018/09/28
//@Purpose: Tcp Socket Server(Epoll模型） ,应用sync map对象

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
func (ds *Devices) Broadcast(rOrder cm.Cmd) {
	for _, c := range ds.Connections {
		_, err := c.Write(rOrder.Cmd)
		if cm.CkErr(TCPMsg.SendError, err) {
			continue
		}
	}
}

// SendOrderByTCP :发送给设备单独的某条命令、数据
func (ds *Devices) SendOrderByTCP(rOrder *cm.Cmd, rWebMsg string) error {
	// 获取tcp连接id
	connID := DecodeToTCPConnID(rOrder.CmdID)
	if _, ok := ds.Connections[connID]; !ok {
		return cm.ConvertStrToErr(TCPMsg.CanNotFindConnection)
	}
	time.Sleep(15 * time.Millisecond)
	_, err := ds.Connections[connID].Write(rOrder.Cmd)
	if cm.CkErr(TCPMsg.SendError, err) {
		return cm.ConvertStrToErr(TCPMsg.SendError)
	}
	if rWebMsg != "" {
		od := cm.NewOrder(rOrder.CmdID, []byte(rWebMsg))
		dh.SendMsgToWeb(od)
	}
	// 	// TODO:抽象化log对象
	logs.LogMain.Info("=> IP："+ds.Connections[connID].RemoteAddr().String(), TCPMsg.SendSuccess)
	return nil
}

// LoopingSendOrders :循环发送设备对象内的指令序列
func (ds *Devices) LoopingSendOrders() {
	for od := range ds.Orders {
		if cm.CkErr(TCPMsg.SendError, ds.SendOrderByTCP(od, TCPMsg.SendSuccess)) {
			//发送不成功，则延迟发送
			cTicker := time.NewTicker(12 * time.Second) // 定时
			lastCk := time.After(1 * time.Minute)       // 延时
			defer cTicker.Stop()
			for i := 0; i < 3; i++ {
				select {
				case <-cTicker.C:
					if !cm.CkErr("", ds.SendOrderByTCP(od, TCPMsg.SendSuccess)) {
						continue
					}
				}
			}
			select {
			case <-lastCk:
				if !cm.CkErr("", ds.SendOrderByTCP(od, TCPMsg.SendSuccess)) {
					continue
				}
			}
			dh.SendMsgToWeb(cm.NewOrder(od.CmdID, []byte(TCPMsg.SendFailureForLongTime)))
		}
	}
}

// RetrieveTCPOrdersFromDataHub : 循环检测datahub包内DeviceTCPOrders对象，加入到发送队列
func (ds *Devices) RetrieveTCPOrdersFromDataHub() {
	for dh.DeviceTCPOrderQueue != nil {
		select {
		case od := <-dh.DeviceTCPOrderQueue:
			ds.Orders <- od
		}
	}
}

// setReadTimeout:设置读数据超时
func setReadTimeout(conn *net.TCPConn, t time.Duration) {
	conn.SetReadDeadline(time.Now().Add(t))
}

// setReadTimeout :设定TCP连接接收数据时间
func (ds *Devices) setReadTimeout(rConnID string, t time.Duration) {
	if rConnID == "" {
		for _, c := range ds.Connections {
			c.SetReadDeadline(time.Now().Add(t))
		}
		return
	}
	if c, ok := ds.Connections[rConnID]; ok {
		c.SetReadDeadline(time.Now().Add(t))
	}
}

//madeConn :连接初始处理(ed)
func (ds *Devices) madeConn(c *net.TCPConn) {
	//初始化连接这个函数被调用
	connID := cm.GetPureIPAddr(c.RemoteAddr().String())
	ds.Lock()
	ds.Connections[connID] = c
	ds.Unlock()
	logs.LogMain.Info("IP:", connID, "上线")
	cm.SepLi(20, "-")
	// ****定时处理(心跳等)
	//	go loopingCall(conn)
}

// lostConn :连接断开
func (ds *Devices) lostConn(c *net.TCPConn) {
	connID := cm.GetPureIPAddr(c.RemoteAddr().String())
	ds.Lock()
	delete(ds.Connections, connID)
	ds.Unlock()
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
		// // orders := ep.CmdSetRcvReconTime(rvID, cm.ConvertDecToBytes(900))
		// // orders := ep.CmdSetRcvCfg(rvID, bIP, bPort)
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

// StartTCPJob :TCP启动服务入口
func (ds *Devices) StartTCPJob(port int) {
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
		go ds.LoopingSendOrders()
		go ds.RetrieveTCPOrdersFromDataHub()
	}
	connStream := make(chan *net.TCPConn)
	//打开N个Goroutine等待连接，Epoll模式
	for i := 0; i < ds.MaxTCPConnect; i++ {
		go func() {
			for cs := range connStream {
				ds.madeConn(cs)
				//	接收数据
				receiveData(cs)
				ds.lostConn(cs)
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
