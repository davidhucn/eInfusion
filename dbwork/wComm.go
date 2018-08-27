package dbwork

import (
	cm "eInfusion/comm"
	"eInfusion/logs"
	"time"

	"github.com/imroc/biu"
)

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

//BinDetectorStat :根据数据生成检测器状态信息
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
	}
}

//CreateQRID ：生成索引编号
//TODO:等待下一步细化
func CreateQRID(rID string) string {
	strBranchCode := "1x0"
	strCategoryCode := "CP"
	//批号
	strPHCode := "xx1"
	strTime := cm.ConvertIntToStr(time.Now().Hour()) + cm.ConvertIntToStr(time.Now().Minute()) + cm.ConvertIntToStr(time.Now().Second())

	return strBranchCode + strCategoryCode + strPHCode + strTime + rID
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
