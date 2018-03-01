package comm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"time"
)

// 获取当前时间
func GetCurrentTime() string {
	strTime := time.Now().Format("2006-01-02 15:04:05") //后面的参数是固定的 否则将无法正常输出
	return strTime
}

// 打印到屏幕
func ShowScreen(v ...interface{}) {
	fmt.Println(v...)
}

// 错误判断处理
func CheckError(err error) {
	if err != nil {

	}
}

//根据参数base转换成指定进制，返回数值
func BaseConvert(ref_intBase int, ref_varContent interface{}) string {
	//	reflect.TypeOf(ref_varContent)
	var strBaseValue string
	switch ref_intBase {
	case 16:
		strBaseValue = "x"
	case 10:
		strBaseValue = "d"
	case 2:
		strBaseValue = "b"
	case 1:
		strBaseValue = "v"
	default:
		strBaseValue = "x"
	}
	strRetValue := fmt.Sprintf("%"+strBaseValue, ref_varContent)
	return strRetValue
}

//	根据指定进制要求，把字符串转换成数字int型
func BaseStrToInt(ref_base int, ref_content string) int {
	intRetValue, err := strconv.ParseInt(ref_content, ref_base, 64)
	if err != nil {
		//		scrPrint(" 字符串转换为数字出错: ", err)
		intRetValue = 0
	}
	return int(intRetValue)
}

// 获取当然路径
func GetCurrentDirectory() string {
	var strPath string
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		strPath = "\\"
	} else {
		strPath = "/"
	}
	dir, _ := os.Getwd()
	strPath = dir + strPath
	return strPath
}

// 转化整形成字符型
func ConvertIntToStr(intContent int) string {
	return strconv.Itoa(intContent)
}

//整形转换成字节
func ConvertIntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

// 写入指定文件，如果没有该文件自动生成
func WriteToFileWithBuffer(f_strPath string, f_strContent string, f_boolWriteAppend bool) bool {
	//	this function for complex content to file
	var intFileOpenMode int
	if f_boolWriteAppend {
		intFileOpenMode = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	} else {
		intFileOpenMode = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	fileHandle, err := os.OpenFile(f_strPath, intFileOpenMode, 0666)
	if err != nil {
		return false
	} else {
		// write to the file
		fileHandle.WriteString("\r\n" + f_strContent)
	}
	defer fileHandle.Close()
	return true
}
