package dbWorks

//import . "eInfusion/comm"

//检测器对象
type Detector struct {
	QRCode  string
	ID      string
	RcvID   string
	Stat    string //十进制表示
	Disable bool
}

//生成二维码字符串
func GetQRCodeStr(ref_strID string) string {
	strBranchCode := "1x0"
	strCategoryCode := "CP"
	strPHCode := "xx1"

	return strBranchCode + strCategoryCode + strPHCode + ref_strID
}
