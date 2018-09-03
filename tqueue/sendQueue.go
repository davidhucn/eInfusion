package tqueue

import (
	cm "eInfusion/comm"
	wk "eInfusion/dbwork"
	ep "eInfusion/protocol"
	logs "eInfusion/tlogs"
	tcp "eInfusion/ttcp"
	"strings"
	"sync"
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
// 主动发送方式
func StartSendQueueListener() {
	// testAfter := time.After(5 * time.Second)   // 延时
	for sdIDStream != nil {
		select {
		case oid := <-sdIDStream:
			var exm sync.Mutex
			// 获取ip地址,即指令索引
			sIPAddr := strings.Split(oid, "@")[1]
			ssn := strings.Split(oid, "@")[0]
			// 确定设备是否在线
			if _, ok := tcp.ClisConnMap[sIPAddr]; ok {
				if cm.CkErr("发送数据错误！", tcp.SendData(tcp.ClisConnMap[sIPAddr], sdOrders[oid])) {
					continue
				} else {
					// 发送成功
					// 把当前结果回写前端
					exm.Lock()
					defer exm.Unlock()
					RcMsgs[ssn] = "请求执行成功，等待反馈结果..."
					delSendQueueMap(sIPAddr)
					return
				}
			} else {
				cTicker := time.NewTicker(12 * time.Second) // 定时
				lastCk := time.After(1 * time.Minute)       // 延时
				// period :=
				defer cTicker.Stop()
				// 不在线,尝试3次，判断是否在线，延时判断，如果在线即发送
				for i := 0; i < 3; i++ {
					select {
					case <-cTicker.C:
						if _, ok := tcp.ClisConnMap[sIPAddr]; ok {
							if cm.CkErr("发送指令错误，不在线！", tcp.SendData(tcp.ClisConnMap[sIPAddr], sdOrders[oid])) {
								continue
							} else {
								// 发送成功
								// 把当前结果回写前端
								exm.Lock()
								defer exm.Unlock()
								RcMsgs[ssn] = "请求执行成功，等待反馈结果..."
								delSendQueueMap(sIPAddr)
								return
							}
						}
					}
				}
				// 最后3分钟后尝试1次
				select {
				case <-lastCk:
					if _, ok := tcp.ClisConnMap[sIPAddr]; ok {
						if cm.CkErr("发送指令错误，不在线！", tcp.SendData(tcp.ClisConnMap[sIPAddr], sdOrders[oid])) {
							continue
						} else {
							// 发送成功
							// 把当前结果回写前端
							exm.Lock()
							defer exm.Unlock()
							RcMsgs[ssn] = "请求执行成功，等待反馈结果..."
							delSendQueueMap(sIPAddr)
							return
						}
					}
				}
				// 多次尝试无效、失败，警报，记录日志
				// 把当前结果回写前端
				exm.Lock()
				defer exm.Unlock()
				RcMsgs[ssn] = "请求执行失败，设备未连线..."
				logs.LogMain.Info("IP地址为：【", sIPAddr, "】多次无法发送数据！,请核查")
			}
		}
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
func AddToSendQueue(rSorderID string, rDetID string, rCmdType uint8, rArgs string) {
	addDet := cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.AddDetect)
	delDet := cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.DelDetect)
	if rCmdType == addDet || rCmdType == delDet {
		// 根据DetID获取RcvID,如果rcvID不为空
		if wk.GetRcvID(rDetID) != "" {
			rcvID := cm.ConvertStrToBytesByPerTwoChar(wk.GetRcvID(rDetID))
			detID := cm.ConvertStrToBytesByPerTwoChar(rDetID)
			// 获取rcv相关的IP
			ipAddr := wk.GetRcvIP(wk.GetRcvID(rDetID))
			// 重组指令标识:由时间戳+IP地址组成
			orderID := rSorderID + "@" + ipAddr
			addSendQueueMap(orderID, ep.CmdOperateDetect(rCmdType, rcvID, 1, detID))
		} else {
			logs.LogMain.Error("错误：没有目标设备编码或者无法获取相关设备编码！")
		}
	}
}
