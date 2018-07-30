package comm

import (
	"os"

	qrcode "github.com/skip2/go-qrcode"
)

//生成二维码图形文件
//默认的程序目录下生成QRPng目录
func CreateQRCodePngFile(ref_strContent string, ref_size int, ref_filename string) error {
	var strPath string
	strPath = GetCurrentDirectory() + "QRPng"
	if !IsExists(strPath) {
		os.Mkdir(strPath, os.ModePerm)
	}
	return qrcode.WriteFile(ref_strContent, qrcode.Medium, ref_size, strPath+GetPathSeparator()+ref_filename)
}

//生成二维码图形
func CreateQRCodeBytes(ref_strContent string, ref_size int) ([]byte, error) {
	var png []byte
	png, err := qrcode.Encode(ref_strContent, qrcode.Medium, ref_size)
	return png, err
}
