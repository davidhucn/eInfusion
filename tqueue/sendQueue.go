package tqueue

import (
	"eInfusion/ttcp"
	cm "eInfusion/comm"
	wk "eInfusion/dbwork"
	ep "eInfusion/protocol"
	tcp "eInfusion/ttcp"
	"time"
)

func addQueueMap(rOrderID string, rData []byte) {
	cMkMutex.Lock()
	defer cMkMutex.Unlock()
	sOrders[rOrderID] = rData
	sIDStream <- rOrderID
}

//  删除socket conn 映射
// func delQueueMap(rOrderID string) {
// 	cDelMutex.Lock()
// 	defer cDelMutex.Unlock()
// 	delete(Orders, rOrderID)
// }

// StartSendQueueListener :启动队列处理平台
func StartSendQueueListener() {
	pingTicker := time.NewTicker(6 * time.Second) // 定时
	testAfter := time.After(5 * time.Second)       // 延时
	for sIDStream != nil {
		// TODO:获取设备对应的IP地址
		
		select {
		case oID <- sIDStream:
			strings.Split(oID, "@")[1]
			// ttcp.SendData(,sOrders[oID])
			// cm.Msg("quit")
			// tcp.SendData()
			return
		// case <-pingTicker.C:
			//发送心跳
			// _, err := SendData(conn, []byte("PING"))
			// if err != nil {
			// 	pingTicker.Stop()
			// 	return
			// }
		// case <-testAfter:
			//	doLog("testAfter:")
			//TODO:日志记录
		}
		select {
		case <-tick.C:
		fmt.Printf("%d: case <-tick.C\n", i)
		default:
	}
}

/////////////////////////sample////////////////////////////////
// var tc Clienter
// 	tc.SendStr = make(chan *Request, 1000)
// 	tc.RecvStr = make(chan string)
// 	tc.Connect()

// 	go ProxySendLoop(&tc)
// 	go ProxyRecvLoop(&tc)
//////////////////////////////////////////////////////////////

// AddToSendQueue :根据参数生成为统一MAP对象(sOrders)，等待发送
// func AddToSendQueue(rSorderID string,rTargetID string,rCmdType uint8, rArgs string) {
	
// }

// func AddToSendQueue (rSorderID string, rTargetID string, rCmdType uint8,rArgs string) {
	// if rCmdType == cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.AddDetect) || rCmdType == cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.DelDetect) {
	// 	// 根据DetID获取RcvID
	// 	if wk.GetRcvID(rTargetID) != "" {
	// 		rcvID := cm.ConvertStrToBytesByPerTwoChar(wk.GetRcvID(rTargetID))
	// 		detID := cm.ConvertStrToBytesByPerTwoChar(rTargetID)
	// 		// 获取rcv相关的IP
	// 		ipAdd :=wk.GetRcvIP(rcvID)
	// 		// 重组指令标识:由时间戳+IP地址组成
	// 		orderID := rSorderID + "@" + ipAdd
	// 		addQueueMap(orderID, ep.CmdOperateDetect(rCmdType, rcvID, 1, detID))
	// 	}
	// }
// }
