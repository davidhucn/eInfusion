package dbWorks

import (
	cm "eInfusion/comm"
	"time"
	// "unicode/utf8"
)

//检测器对象
// Stat:工作状态,0-关机，1-开机
// Alarm: 是否报警，输液条没有液体
type Detector struct {
	QRCode   string
	ID       string
	RcvID    string
	Capacity int //0,1,2,3
	Stat     int //工作状态：0-关机，1-开机
	Disable  bool
	Alarm    bool //是否报警，输液条没有液体
}

//根据数据生成检测器状态信息
func BinDetectorStat(rdata byte, dt *Detector) {
	// t, _ := strconv.ParseUint(string(sMeta[1]), 10, 64)
	bib := cm.ConvertByteToBinaryOfByte(rdata)
	// TODO:完成数据解析

	// dt.Disable = bib[0]
	// dt.Stat = bib[1]
	// dt.Capacity = bib[2:3]
	// dt.Alarm = bib[4]
	for _, v := range bib {
		cm.Msg(v)
	}
	// cm.Msg(cm.ConvertBasNumberToStr(10, sMeta[0]))
	// if cm.ConvertBasNumberToStr(2, sMeta[0]) == 48 {
	// 	cm.Msg("zero")
	// }

}

//生成索引编号
//TODO:等待下一步细化
func CreateQRID(ref_strID string) string {
	strBranchCode := "1x0"
	strCategoryCode := "CP"
	//批号
	strPHCode := "xx1"
	strTime := cm.ConvertIntToStr(time.Now().Hour()) + cm.ConvertIntToStr(time.Now().Minute()) + cm.ConvertIntToStr(time.Now().Second())

	return strBranchCode + strCategoryCode + strPHCode + strTime + ref_strID
}

// 生成二维码
func StartCreateQRCode() {
	//auto create the qrcode
	for i := 0; i < 10; i++ {
		strName := "B000000" + cm.ConvertIntToStr(i)
		strContent := CreateQRID(strName)
		cm.CreateQRCodePngFile(strContent, 128, strName+".png")
	}
}
