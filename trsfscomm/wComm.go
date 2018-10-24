package trsfscomm

import (
	cm "eInfusion/comm"
	logs "eInfusion/tlogs"

	"github.com/imroc/biu"
)

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

//StartCreateQRCode ：生成二维码
func StartCreateQRCode() {
	//auto create the qrcode
	for i := 0; i < 10; i++ {
		strName := "B000000" + cm.ConvertIntToStr(i)
		strContent := cm.CreateQRID(strName)
		cm.CreateQRCodePngFile(strContent, 128, strName+".png")
	}
}
