package dbWorks

import . "eInfusion/comm"

import "time"

//检测器对象
type Detector struct {
	QRCode  string
	ID      string
	RcvID   string
	Stat    string //十进制表示
	Disable bool
}

//生成索引编号
//等待下一步细化
func CreateQRID(ref_strID string) string {
	strBranchCode := "1x0"
	strCategoryCode := "CP"
	//批号
	strPHCode := "xx1"
	strTime := ConvertIntToStr(time.Now().Hour()) + ConvertIntToStr(time.Now().Minute()) + ConvertIntToStr(time.Now().Second())

	return strBranchCode + strCategoryCode + strPHCode + strTime + ref_strID
}

func StartCreateBRCode() {
	//auto create the qrcode
	for i := 0; i < 10; i++ {
		strName := "B000000" + ConvertIntToStr(i)
		strContent := CreateQRID(strName)
		CreateQRCodePngFile(strContent, 128, strName+".png")
	}
}
