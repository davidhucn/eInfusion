package trsfus

import (
	cm "eInfusion/comm"
	"sync"
	"time"

	"github.com/imroc/biu"
)

// CmdType :指令类型
type CmdType int

const (
	_ CmdType = iota
	// GetReceiverState :获取接收器状态
	GetReceiverState
	// GetDetectorState :获取检测器状态
	GetDetectorState
	// AddDetector :添加检测器
	AddDetector
	// DeleteDetector :删除检测器
	DeleteDetector
	// SetReceiverConfig :设置接收器参数
	SetReceiverConfig
	// SetReceiverReconnectTime :设置接收器心跳连接时间
	SetReceiverReconnectTime
)

// Order :指令对象
type Order struct {
	RcvID []byte
	DetID []byte
	Cmd   []byte
	Args  []byte
}

// OrdersQueue :指令对象结合映射，用于接收-发送配对
var OrdersQueue sync.Map // ([string:ID]Order)

// NewOrder :新建指令对象
func NewOrder(rcvID []byte, detID []byte, cmd []byte, args []byte) *Order {
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
	reconnectTime time.Duration
	targetIP      string
	targetPort    string
}

// NewReceiver ：新建一个接收器对象
func NewReceiver(id string, detectAmount int, nativeIP string, reconnectTime time.Duration,
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
	SendCmdMap[GetReceiverState] = 0x00
	SendCmdMap[GetDetectorState] = 0x01
	SendCmdMap[AddDetector] = 0x02
	SendCmdMap[DeleteDetector] = 0x03
	SendCmdMap[SetReceiverConfig] = 0x04
	SendCmdMap[SetReceiverReconnectTime] = 0x05

	ReceiveCmdMap[0x10] = GetReceiverState
	ReceiveCmdMap[0x11] = GetDetectorState
	ReceiveCmdMap[0x12] = AddDetector
	ReceiveCmdMap[0x13] = DeleteDetector
	ReceiveCmdMap[0x14] = SetReceiverConfig
	ReceiveCmdMap[0x15] = SetReceiverReconnectTime
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
		logs.LogMain.Info("试图获取检测器：[", dt.ID, "]信息无效，可能没电或者未启动！")
		// tlogs.DoLog()
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
