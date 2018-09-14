package thttp

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

// WriteBackWsByID : 通过ip地址找到对应的ws连接，并回写到前端（供外部调用）
func WriteBackWsByID(rIPaddr string, rStrCnt string) bool {
	//TODO:
	// if _, yes := WsClis[rIP]; yes {
	// 	WsClis[sn].sdData <- []byte(strCnt)
	// 	return true
	// }
	// ws不在线时,待处理
	// return false
	return true
}

// 回写到前端ws应用消息,模块内部使用
func wsWriteBack(sn string, strCnt string) bool {
	if _, yes := WsClis[sn]; yes {
		WsClis[sn].sdData <- []byte(strCnt)
		return true
	}
	// ws不在线时,待处理
	return false
}

// StartSendQueueListener :启动队列处理平台【主动发送模式】
func StartSendQueueListener() {
	// testAfter := time.After(5 * time.Second)   // 延时
	for sdIDStream != nil {
		select {
		case oid := <-sdIDStream:
			// var exm sync.Mutex
			// 获取ip地址,即指令索引
			sIPAddr := strings.Split(oid, "@")[1]
			ssn := strings.Split(oid, "@")[0]
			// 确定设备是否在线
			if _, ok := tcp.ClisConnMap[sIPAddr]; ok {
				//直接发送指令
				if cm.CkErr("发送指令至设备错误！", tcp.SendData(tcp.ClisConnMap[sIPAddr], sdOrders[oid])) {
					logs.LogMain.Error("发送指令失败，IP：", sIPAddr) // 发送失败时,日志记录
					continue
				} else {
					// 发送成功
					// 如果ws在线则把回传信息
					wsWriteBack(ssn, "设备操作指令发送成功，请稍候...")
					delSendQueueMap(sIPAddr)
					continue
				}
			} else { //当设备不在线时
				cTicker := time.NewTicker(12 * time.Second) // 定时
				lastCk := time.After(1 * time.Minute)       // 延时
				defer cTicker.Stop()
				// 不在线,尝试3次，判断是否在线，延时判断，如果在线即发送
				for i := 0; i < 3; i++ {
					select {
					case <-cTicker.C:
						if _, ok := tcp.ClisConnMap[sIPAddr]; ok {
							if cm.CkErr("发送指令错误，系统检测为在线状态！", tcp.SendData(tcp.ClisConnMap[sIPAddr], sdOrders[oid])) {
								continue
							} else {
								// 发送成功
								// 如果ws在线则把回传信息
								wsWriteBack(ssn, "设备操作指令发送成功，请稍候...")
								delSendQueueMap(sIPAddr)
								continue
							}
						}
					}
				}
				// 最后3分钟后尝试1次
				select {
				case <-lastCk:
					if _, ok := tcp.ClisConnMap[sIPAddr]; ok {
						if cm.CkErr("发送指令错误，系统检测为在线状态！", tcp.SendData(tcp.ClisConnMap[sIPAddr], sdOrders[oid])) {
							continue
						} else {
							// 发送成功
							// 如果ws在线则把回传信息
							wsWriteBack(ssn, "设备操作指令发送成功，请稍候...")
							delSendQueueMap(sIPAddr)
							continue
						}
					}
				}
				// 多次尝试无效、失败，警报，记录日志
				// 如果ws在线则把回传信息
				wsWriteBack(ssn, "由于设备长时间断线，操作指令发送失败...")
				delSendQueueMap(sIPAddr)
				logs.LogMain.Info("IP地址为：【", sIPAddr, "】的设备多次通讯失败！,请核查")
			}
		}
	}
}

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
