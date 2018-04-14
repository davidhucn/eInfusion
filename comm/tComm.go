package comm

import (
	"bytes"
	"eInfusion/logs"
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// 获取当前时间
func GetCurrentTime() string {
	strTime := time.Now().Format("2006-01-02 15:04:05") //后面的参数是固定的 否则将无法正常输出
	return strTime
}

//获取当前日期
func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

//生成分隔行
func SepLi(ref_num int) string {
	s := "-"
	Msg(strings.Repeat(s, ref_num))
}

//去除左右空格
func TrimSpc(ref_str string) string {
	return strings.TrimSpace(ref_str)
}

// 打印到屏幕
func Msg(v ...interface{}) {
	fmt.Println(v...)
}

//获取变量类型
func GetVarType(ref_var interface{}) string {
	return fmt.Sprint(reflect.TypeOf(ref_var))
}

//处理错误
//如果有错误，返回true,无错则返回false
func CkErr(ref_Msg string, ref_err error) bool {
	if ref_err != nil {
		logs.LogMain.Error(ref_Msg, ref_err)
		return true
	}
	return false
}

//string转byte
func ConvertStrToBytes(s string) byte {
	return *(*byte)(unsafe.Pointer(&s))
}

//byte转string
// convert b to string without copy
func ConvertBytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//根据参数base转换成指定进制，返回
func ConvertBasToStr(ref_intBase int, ref_varContent interface{}) string {
	//	reflect.TypeOf(ref_varContent)
	var strBaseValue string
	switch ref_intBase {
	case 16:
		strBaseValue = "X"
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

//转换16进制Bytes为string
func ConvertOxBytesToStr(ref_content []byte) string {
	var strRet string
	for i := 0; i < len(ref_content); i++ {
		strCon := ConvertBasToStr(16, ref_content[i])
		if len(strCon) == 1 {
			strCon = "0" + strCon
		}
		strRet += strCon
	}
	return strRet
}

//	根据指定进制要求，把字符串转换成数字int型
func ConvertBasStrToInt(ref_base int, ref_content string) int {
	intRetValue, err := strconv.ParseInt(ref_content, ref_base, 64)
	if err != nil {
		intRetValue = 0
	}
	return int(intRetValue)
}

//	根据指定进制要求，把字符串转换成数字Uint8型
func ConvertBasStrToUint(ref_base int, ref_content string) uint8 {
	intRetValue, err := strconv.ParseUint(ref_content, ref_base, 64)
	if err != nil {
		intRetValue = 0
	}
	return uint8(intRetValue)
}

//判断是否存在,true表示存在
func IsExists(ref_Path string) bool {
	_, err := os.Stat(ref_Path)
	if err != nil {
		return false
	}
	return true
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

//根据环境，返回分隔符
func GetPathSeparator() string {
	var strPath string
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		strPath = "\\"
	} else {
		strPath = "/"
	}
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
func WrToFilWithBuffer(f_strPath string, f_strContent string, f_boolWriteAppend bool) bool {
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
