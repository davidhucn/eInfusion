package comm

import (
	"fmt"
	"os"
	"strconv"
)

// 打印到屏幕
func ScrPrint(v ...interface{}) {
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
