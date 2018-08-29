package tqueue

import (
	cm "eInfusion/comm"
	wk "eInfusion/dbwork"
	ep "eInfusion/protocol"
	logs "eInfusion/tlogs"
	tcp "eInfusion/ttcp"
	"strings"
	"time"
)

func addSendQueueMap(rOrderID string, rData []byte) {
	cMkMutex.Lock()
	defer cMkMutex.Unlock()
	sdOrders[rOrderID] = rData
	sdIDStream <- rOrderID
}

// delQueueMap :删除发送指令映射
func delSendQueueMap(rOrderID string) {
	cDelMutex.Lock()
	defer cDelMutex.Unlock()
	delete(sdOrders, rOrderID)
}

// StartSendQueueListener :启动队列处理平台
func StartSendQueueListener() {
	cTicker := time.NewTicker(12 * time.Second) // 定时
	// testAfter := time.After(5 * time.Second)   // 延时
	for sdIDStream != nil {
		select {
		case oid := <-sdIDStream:
			// 获取ip地址,即指令索引
			sIPAddr := strings.Split(oid, "@")[1]
			// 确定设备是否在线
			if _, ok := tcp.ClisConnMap[sIPAddr]; ok {
				if cm.CkErr("发送数据错误！", tcp.SendData(tcp.ClisConnMap[sIPAddr], sdOrders[oid])) {
					continue
				} else {
					// 发送成功
					// TODO: 返回消息处理
					delSendQueueMap(sIPAddr)
				}
			} else {
				// 不在线,尝试3次，判断是否在线，延时判断，如果在线即发送
				for i := 0; i < 3; i++ {
					select {
					case <-cTicker.C:
						if _, ok := tcp.ClisConnMap[sIPAddr]; ok {
							if cm.CkErr("发送指令至设备错误，不在线！", tcp.SendData(tcp.ClisConnMap[sIPAddr], sdOrders[oid])) {
								continue
							} else {
								// 发送成功
								delSendQueueMap(sIPAddr)
							}
						}
					}
				}
				// 多次尝试无效，警报，记录日志
				// logs.LogMain.Critical("IP地址为：【", sIPAddr, "】多次无法发送数据！,请核查")
				logs.LogMain.Debug("IP地址为：【", sIPAddr, "】多次无法发送数据！,请核查")
			}
			// default:
			// 	logs.LogMain.Debug(sdIDStream, "-- waiting for send queue!")
		}
	}
	// close(sdIDStream)
	// cTicker.Stop()
}

// StartReceiveQueueListener : 监听设备返回消息队列
func StartReceiveQueueListener() {

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
		// 根据DetID获取RcvID,如果rcvID不为空
		if wk.GetRcvID(rTargetID) != "" {
			rcvID := cm.ConvertStrToBytesByPerTwoChar(wk.GetRcvID(rTargetID))
			detID := cm.ConvertStrToBytesByPerTwoChar(rTargetID)
			// 获取rcv相关的IP
			ipAddr := wk.GetRcvIP(wk.GetRcvID(rTargetID))
			// 重组指令标识:由时间戳+IP地址组成
			orderID := rSorderID + "@" + ipAddr
			addSendQueueMap(orderID, ep.CmdOperateDetect(rCmdType, rcvID, 1, detID))
		}
	}
}
