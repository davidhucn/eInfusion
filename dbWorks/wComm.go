package dbWorks

import (
	"eInfusion/comm"
	. "eInfusion/comm"

	"time"
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
func BinDetectorStat(da byte,

// FIXME:测试，须修改
// dt *Detector
) []byte {
	meta := comm.ConvertBasNumberToStr(2, da)
	return comm.ConvertPerTwoOxCharOfStrToBytes(meta)
}

//生成索引编号
//TODO:等待下一步细化
func CreateQRID(ref_strID string) string {
	strBranchCode := "1x0"
	strCategoryCode := "CP"
	//批号
	strPHCode := "xx1"
	strTime := ConvertIntToStr(time.Now().Hour()) + ConvertIntToStr(time.Now().Minute()) + ConvertIntToStr(time.Now().Second())

	return strBranchCode + strCategoryCode + strPHCode + strTime + ref_strID
}

// 生成二维码
func StartCreateQRCode() {
	//auto create the qrcode
	for i := 0; i < 10; i++ {
		strName := "B000000" + ConvertIntToStr(i)
		strContent := CreateQRID(strName)
		CreateQRCodePngFile(strContent, 128, strName+".png")
	}
}
