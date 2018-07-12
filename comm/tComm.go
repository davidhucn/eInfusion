package comm

import (
	"bytes"
	"eInfusion/logs"
	"encoding/binary"

	// "encoding/hex"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//获取时间戳
func GetTimeStamp() string {
	return string(time.Now().Unix())
}

func GetRealIPAddr(ip string) string {
	return strings.Split(ip, ":")[0]
}

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
func SepLi(ref_num int, ref_Char string) {
	s := "-"
	if ref_Char != "" {
		s = ref_Char
	}
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
func CkErr(ref_MsgTitle string, ref_err error) bool {
	if ref_err != nil {
		logs.LogMain.Error(ref_MsgTitle, ref_err)
		return true
	}
	return false
}

//把数字内容（int/uint）根据参数base转换成指定进制，返回
func ConvertBasNumberToStr(ref_intBase int, ref_varContent interface{}) string {
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
//转换过程中可能会丢失0，因此需要补0
func ConvertOxBytesToStr(ref_content []byte) string {
	var strRet string
	for i := 0; i < len(ref_content); i++ {
		strCon := ConvertBasNumberToStr(16, ref_content[i])
		if len(strCon) == 1 {
			strCon = "0" + strCon
		}
		strRet += strCon
	}
	return strRet
}

//	把指定进制的字符转换成为十进制数值（int型）
func ConvertBasStrToInt(ref_intBase int, ref_content string) int {
	intRetValue, err := strconv.ParseInt(ref_content, ref_intBase, 64)
	if err != nil {
		intRetValue = 0
	}
	return int(intRetValue)
}

//	根据指定进制要求，把字符串转换成数字Uint8型
func ConvertBasStrToUint(ref_intBase int, ref_content string) uint8 {
	intRetValue, err := strconv.ParseUint(ref_content, ref_intBase, 64)
	if err != nil {
		intRetValue = 0
	}
	return uint8(intRetValue)
}

//根据开始、结束下标返回相应的字符串内容返回bytes
func ConvertBasStrToBytes(ref_strContent string, ref_intBegin int, ref_intEnd int, ref_intBase int) []byte {
	var bT []byte
	n := len(ref_strContent)
	if ref_intEnd < n && ref_intBegin > 0 {
		for i := ref_intBegin; i <= ref_intEnd; i++ {
			strT := string(ref_strContent[i])
			bT = append(bT, ConvertBasStrToUint(ref_intBase, strT))
		}
	}
	return bT
}

//根据开始、结束下标返回相应的字符串内容返回string
func GetPartOfStrToStr(ref_strContent string, ref_intBegin int, ref_intEnd int) string {
	var strR string
	n := len(ref_strContent)
	if ref_intEnd < n && ref_intBegin >= 0 {
		for i := ref_intBegin; i <= ref_intEnd; i++ {
			strR += string(ref_strContent[i])
		}
	}
	return strR
}

//把字符串内容按每两字符对应一个byte组成新的bytes，返回[]byte
func ConvertPerTwoOxCharOfStrToBytes(ref_s string) []byte {
	var bT []byte
	var n int = len(ref_s)
	for i := 0; i < n-1; i++ {
		j := i + 1
		strP := string(ref_s[i]) + string(ref_s[j])
		bT = append(bT, ConvertBasStrToUint(16, string(strP)))
	}
	return bT
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
