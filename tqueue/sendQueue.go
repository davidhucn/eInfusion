package tqueue

import (
	cm "eInfusion/comm"
	wk "eInfusion/dbwork"
	ep "eInfusion/protocol"
)

// StartSendQueueListener :启动队列处理平台
func StartSendQueueListener() {
	// pingTicker := time.NewTicker(10 * time.Second) // 定时
	// testAfter := time.After(5 * time.Second)       // 延时

	// for HTTPReqStream != nil {
	// 	select {
	// 	case <-HTTPReqStream:
	// 		cm.Msg("quit")
	// 		return
	// 	case <-pingTicker.C:
	// 		//发送心跳
	// 		// _, err := SendData(conn, []byte("PING"))
	// 		// if err != nil {
	// 		// 	pingTicker.Stop()
	// 		// 	return
	// 		// }
	// 	case <-testAfter:
	// 		//	doLog("testAfter:")
	// 		//TODO:日志记录
	// 	}
	// 	// select {
	// 	// case <-tick.C:
	// 	// fmt.Printf("%d: case <-tick.C\n", i)
	// 	// default:
	// }
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
func AddToSendQueue(rSorderID string, rTargetID string, rCmdType uint8, rArgs string) {
	if rCmdType == cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.AddDetect) || rCmdType == cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.DelDetect) {
		// 根据DetID获取RcvID
		rcvID := cm.ConvertStrToBytesByPerTwoChar(wk.GetRcvID(rTargetID))
		detID := cm.ConvertStrToBytesByPerTwoChar(rTargetID)
		// sOrders[rSorderID] <- ep.CmdOperateDetect(rCmdType, rcvID, 1, detID)
		// m2[rSorderID] <- "chan"
	}
}
