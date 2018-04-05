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

//生成索引编号
//等待下一步细化
func CreateQRID(ref_strID string) string {
	strBranchCode := "1x0"
	strCategoryCode := "CP"
	//批号
	strPHCode := "xx1"

	return strBranchCode + strCategoryCode + strPHCode + ref_strID
}
