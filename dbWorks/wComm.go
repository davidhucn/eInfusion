package dbWorks

//检测器对象
type Detector struct {
	Qcode   string
	ID      string
	RcvID   string
	Stat    string //十进制表示
	Disable bool
}

//生成二维码字符串
func GetQcodeStr(ref_ID) string {
	strBranchCode := "1x0"
	strUnitCode := "CP"
	strPHCode := "xx1"

	return ""
}
