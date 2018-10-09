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

// SendOrderIndicate :发送给设备单独的某条命令和数据
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
func (ds *Devices) Broadcast(rOrder *cm.Cmd) {
	for _, c := range ds.Connections {
		_, err := c.Write(rOrder.Cmd)
		if cm.CkErr(TCPMsg.SendError, err) {
			break
		}
	}
}

// LoopingSendTCPOrders : 循环检测datahub包内DeviceTCPOrders对象，发送指令到相应设备去
func (ds *Devices) LoopingSendTCPOrders() {
	for dh.DeviceTCPOrder != nil {
		select {
		case od := <-dh.DeviceTCPOrder:
			conID := strings.Split(od.CmdID, "@")[1]
			if _, ok := ds.Connections[conID]; !ok {
				cm.CkErr("发送指令至设备错误！", cm.ConvertStrToErr(TCPMsg.CanNotFindConnection))
				continue
			}
			time.Sleep(15 * time.Millisecond)
			_, err := ds.Connections[conID].Write(od.Cmd)
			if cm.CkErr(TCPMsg.SendError, err) {
				continue
			} else {
				// 发送成功

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

//连接初始处理(ed)
func (d *Devices) madeConn(c *net.TCPConn) {
	//初始化连接这个函数被调用
	connID := cm.GetPureIPAddr(c.RemoteAddr().String())
	d.Lock()
	d.Connections[connID] = c
	d.Unlock()
	logs.LogMain.Info("IP:", connID, "上线")
	cm.SepLi(20, "-")
	// ****定时处理(心跳等)
	//	go loopingCall(conn)
}

// 连接断开
func (d *Devices) lostConn(c *net.TCPConn) {
	connID := cm.GetPureIPAddr(c.RemoteAddr().String())
	d.Lock()
	delete(d.Connections, connID)
	d.Unlock()
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
func StartTCPService(d *Devices, port int) {
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
	connStream := make(chan *net.TCPConn)
	//打开N个Goroutine等待连接，Epoll模式
	for i := 0; i < MaxTCPConnectLimit; i++ {
		go func() {
			for cs := range connStream {
				d.madeConn(cs)
				//	接收数据
				receiveData(cs)
				d.lostConn(cs)
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
