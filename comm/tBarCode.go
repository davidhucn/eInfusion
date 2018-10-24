package comm

import (
	"os"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

//CreateQRCodePngFile :生成二维码图形文件
//默认的程序目录下生成QRPng目录
func CreateQRCodePngFile(rStrCnt string, rCntSize int, rSfileName string) error {
	var strPath string
	strPath = GetCurrentDirectory() + "QRPng"
	if !IsExists(strPath) {
		os.Mkdir(strPath, os.ModePerm)
	}
	return qrcode.WriteFile(rStrCnt, qrcode.Medium, rCntSize, strPath+GetPathSeparator()+rSfileName)
}

//CreateQRCodeBytes :生成二维码图形
func CreateQRCodeBytes(rStrCnt string, rCntSize int) ([]byte, error) {
	var png []byte
	png, err := qrcode.Encode(rStrCnt, qrcode.Medium, rCntSize)
	return png, err
}

//CreateQRID ：生成索引编号
//TODO:等待下一步细化
func CreateQRID(rID string) string {
	strBranchCode := "1x0"
	strCategoryCode := "CP"
	//批号
	strPHCode := "xx1"
	strTime := ConvertIntToStr(time.Now().Hour()) + ConvertIntToStr(time.Now().Minute()) + ConvertIntToStr(time.Now().Second())
	return strBranchCode + strCategoryCode + strPHCode + strTime + rID
}
