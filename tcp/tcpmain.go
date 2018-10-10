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
	"strings"
	"time"
)

// SendOrderAndWebMsg :发送指令至设备并回传指定消息
func (ds *Devices) SendOrderAndWebMsg(rOrder *cm.Cmd, rMsg string) bool {
	// 发送指令失败时可能为设备不在线或者设备接收故障，成功则返回前端消息
	if !cm.CkErr(TCPMsg.SendError, ds.SendOrderIndicate(rOrder)) {
		// 发送成功，回写到web
		od := cm.NewOrder(rOrder.CmdID, []byte(rMsg))
		dh.SendMsgBackToWeb(od)
		return true
	}
	return false
}

// SendOrderIndicate :发送给设备单独的某条命令、数据
func (ds *Devices) SendOrderIndicate(rOrder *cm.Cmd) error {
	// 获取tcp连接id
	conID := strings.Split(rOrder.CmdID, "@")[1]
	if _, ok := ds.Connections[conID]; !ok {
		return cm.ConvertStrToErr(TCPMsg.CanNotFindConnection)
	}
	time.Sleep(15 * time.Millisecond)
	_, err := ds.Connections[conID].Write(rOrder.Cmd)
	if cm.CkErr(TCPMsg.SendError, err) {
		return cm.ConvertStrToErr(TCPMsg.SendError)
	}
	// 	// TODO:抽象化log对象
	logs.LogMain.Info("=> IP："+ds.Connections[conID].RemoteAddr().String(), TCPMsg.SendSuccess)
	return nil
}

// Broadcast :对所有连接发送广播
func (ds *Devices) Broadcast(rOrder cm.Cmd) {
	for _, c := range ds.Connections {
		_, err := c.Write(rOrder.Cmd)
		if cm.CkErr(TCPMsg.SendError, err) {
			continue
		}
	}
}

// LoopingSendTCPOrders : 循环检测datahub包内DeviceTCPOrders对象，发送指令到相应设备去
func (ds *Devices) LoopingSendTCPOrders() {
	for dh.DeviceTCPOrderQueue != nil {
		select {
		case od := <-dh.DeviceTCPOrderQueue:
			if !ds.SendOrderAndWebMsg(od, TCPMsg.SendSuccess) {
				cTicker := time.NewTicker(12 * time.Second) // 定时
				lastCk := time.After(1 * time.Minute)       // 延时
				defer cTicker.Stop()
				// 不在线,尝试3次，判断是否在线，延时判断，如果在线即发送
				for i := 0; i < 3; i++ {
					select {
					case <-cTicker.C:
						if ds.SendOrderAndWebMsg(od, TCPMsg.SendSuccess) {
							continue
						}
					}
				}
				select {
				case <-lastCk:
					ds.SendOrderAndWebMsg(od, TCPMsg.SendSuccess)
				}
				dh.SendMsgBackToWeb(cm.NewOrder(od.CmdID, []byte(TCPMsg.SendFailureForLongTime)))
			}
		}
	}

}

//   定时处理&延时处理
//func loopingCall(conn *net.TCPConn) {
//	pingTicker := time.NewTicker(30 * time.Second) // 定时
//	testAfter := time.After(5 * time.Second)       // 延时

//	for {
//		select {
//		case <-pingTicker.C:
//			//发送心跳
//			_, err := SendData(conn, []byte("PING"))
//			if err != nil {
//				pingTicker.Stop()
//				return
//			}
//		case <-testAfter:
//			//	doLog("testAfter:")
//			//TODO:日志记录
//		}
//	}
//}

// setReadTimeout:设置读数据超时
func setReadTimeout(conn *net.TCPConn, t time.Duration) {
	conn.SetReadDeadline(time.Now().Add(t))
}

// func (ds *Devices) setReadTimeout(t time.Duration) {
// 	for cs := range ds.Connections {

// 	}
// }

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

// StartTCPService :TCP启动服务入口
func StartTCPService(ds *Devices, port int) {
	host := ":" + strconv.Itoa(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if cm.CkErr(TCPMsg.SourceError, err) {
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	defer listener.Close()
	if cm.CkErr(TCPMsg.SourceError, err) {
		os.Exit(1)
	}
	cm.SepLi(60, "")
	logs.LogMain.Info(TCPMsg.StartServiceMsg + "（" + cm.GetCurrentDate() + "）")
	cm.Msg("Transfusion System Server Port", host)
	cm.SepLi(60, "")
	// 循环发送指令
	{
		go ds.LoopingSendTCPOrders()
	}
	connStream := make(chan *net.TCPConn)
	//打开N个Goroutine等待连接，Epoll模式
	for i := 0; i < MaxTCPConnectLimit; i++ {
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
