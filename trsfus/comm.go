package trsfus

import (
	cm "eInfusion/comm"
	log "eInfusion/tlogs"
	"sync"
	"time"

	"github.com/imroc/biu"
)

// packetHeaderPrefix :帧头前缀数据
var packetHeaderPrefix []byte

// MakePacketHeaderPrefix :生成数据包头前缀
func MakePacketHeaderPrefix(d uint8) []byte {
	packetHeaderPrefix = append(packetHeaderPrefix, d)
	return packetHeaderPrefix
}

// CmdType :指令类型
type CmdType int

const (
	_ CmdType = iota
	// CmdGetReceiverState :获取接收器状态
	CmdGetReceiverState
	// CmdGetDetectorState :获取检测器状态
	CmdGetDetectorState
	// CmdAddDetector :添加检测器
	CmdAddDetector
	// CmdDeleteDetector :删除检测器
	CmdDeleteDetector
	// CmdSetReceiverConfig :设置接收器参数
	CmdSetReceiverConfig
	// CmdSetReceiverReconnectTime :设置接收器心跳连接时间
	CmdSetReceiverReconnectTime
)

// Order :指令对象
type Order struct {
	RcvID string
	DetID string
	Cmd   CmdType
	Args  []string
}

// OrdersQueue :指令对象结合映射，用于接收-发送配对
var OrdersQueue sync.Map // ([string:ID]Order)

// NewOrder :新建指令对象
func NewOrder(rcvID string, detID string, cmd CmdType, args []string) *Order {
	return &Order{
		RcvID: rcvID,
		DetID: detID,
		Cmd:   cmd,
		Args:  args,
	}
}

// Receiver ：检测器对象
type Receiver struct {
	ID            string
	nativeIP      string
	detectAmount  int
	reconnectTime int
	targetIP      string
	targetPort    string
}

// NewReceiver ：新建一个接收器对象
func NewReceiver(id string, detectAmount int, nativeIP string, reconnectTime int,
	targetIP string, targetPort string) *Receiver {
	return &Receiver{
		ID:            id,
		detectAmount:  detectAmount,
		nativeIP:      nativeIP,
		reconnectTime: reconnectTime,
		targetIP:      targetIP,
		targetPort:    targetPort,
	}
}

//Detector :检测器对象
// Stat:工作状态,0-关机，1-开机
// Alarm: 是否报警，输液条没有液体
type Detector struct {
	QRCode   string
	ID       string
	RcvID    string
	Capacity uint8 //0,1,2,3
	PowerOn  uint8 //工作状态：0-关机，1-开机
	Alarm    uint8 //是否报警，0-正常，1－报警，无药水
}

// NewDetector :新建检测器对象
// Capacity uint8 【0,1,2,3】
// PowerOn  uint8 【工作状态：0-关机，1-开机】
// Alarm    uint8 【是否报警，0-正常，1－报警，无药水】
func NewDetector(qr string, id string, rcvID string, capacity uint8, powerOn uint8, alarm uint8) *Detector {
	return &Detector{
		QRCode:   qr,
		ID:       id,
		RcvID:    rcvID,
		Capacity: capacity,
		PowerOn:  powerOn,
		Alarm:    alarm,
	}
}

// ReceiveCmdMap :接收的指令映射
var ReceiveCmdMap = make(map[byte]CmdType, 6)

// SendCmdMap :发送指令映射
var SendCmdMap = make(map[CmdType]byte, 6)

func init() {
	SendCmdMap[CmdGetReceiverState] = 10
	SendCmdMap[CmdGetDetectorState] = 11
	SendCmdMap[CmdAddDetector] = 12
	SendCmdMap[CmdDeleteDetector] = 13
	SendCmdMap[CmdSetReceiverConfig] = 14
	SendCmdMap[CmdSetReceiverReconnectTime] = 15

	ReceiveCmdMap[0x00] = CmdGetReceiverState
	ReceiveCmdMap[0x01] = CmdGetDetectorState
	ReceiveCmdMap[0x02] = CmdAddDetector
	ReceiveCmdMap[0x03] = CmdDeleteDetector
	ReceiveCmdMap[0x04] = CmdSetReceiverConfig
	ReceiveCmdMap[0x05] = CmdSetReceiverReconnectTime

	packetHeaderPrefix = make([]byte, 0)
}

//BinDetectorStat :根据通讯协议，对byte数据生成检测器状态信息（bit）
// 注：目前夹断功能没有开放
func BinDetectorStat(rdata byte, dt *Detector) {
	smd := biu.ByteToBinaryString(rdata)
	// 数据为7位表示检测器状态,如果为0则表示没有打开（如：没电，等)，不进行后续解析
	if string(smd[6]) != "0" {
		dt.PowerOn = cm.ConvertBasStrToUint(10, string(smd[6]))
		dt.Alarm = cm.ConvertBasStrToUint(10, string(smd[3]))
		st := "000000" + string(smd[5]) + string(smd[4])
		biu.ReadBinaryString(st, &dt.Capacity)
	} else {
		dt.PowerOn = 0
		dt.Alarm = 0
		dt.Capacity = 0
		log.DoLog(log.Info, "试图获取检测器：[", dt.ID, "]信息无效，可能没电或者未启动！")
	}
}

//StartCreateQRCode ：生成二维码
func StartCreateQRCode() {
	//auto create the qrcode
	for i := 0; i < 10; i++ {
		strName := "B000000" + cm.ConvertIntToStr(i)
		strContent := CreateQRID(strName)
		cm.CreateQRCodePngFile(strContent, 128, strName+".png")
	}
}

//CreateQRID ：生成索引编号
//TODO:等待下一步细化(硬件供应商提供方案)
func CreateQRID(rID string) string {
	strBranchCode := "1x0"
	strCategoryCode := "CP"
	//批号
	strPHCode := "xx1"
	strTime := cm.ConvertIntToStr(time.Now().Hour()) + cm.ConvertIntToStr(time.Now().Minute()) + cm.ConvertIntToStr(time.Now().Second())
	return strBranchCode + strCategoryCode + strPHCode + strTime + rID
}
