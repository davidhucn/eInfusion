package comm

import (
	"bytes"
	logs "eInfusion/tlogs"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// GetRandString : 生成随机字符串
func GetRandString(length int) string {
	rand.Seed(time.Now().UnixNano())
	rs := make([]string, length)
	for start := 0; start < length; start++ {
		t := rand.Intn(3)
		if t == 0 {
			rs = append(rs, strconv.Itoa(rand.Intn(10)))
		} else if t == 1 {
			rs = append(rs, string(rand.Intn(26)+65))
		} else {
			rs = append(rs, string(rand.Intn(26)+97))
		}
	}
	return strings.Join(rs, "")
}

// GetPureIPAddr : 获取IP中纯的地址，去除字符串中的端口数据
func GetPureIPAddr(ip string) string {
	return strings.Split(ip, ":")[0]
}

// GetCurrentTime : 获取当前时间
func GetCurrentTime() string {
	strTime := time.Now().Format("2006-01-02 15:04:05") //后面的参数是固定的 否则将无法正常输出
	return strTime
}

// GetCurrentDate :获取当前日期
func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// SepLi :生成分隔行
func SepLi(rN int, rChar string) {
	s := "-"
	if rChar != "" {
		s = rChar
	}
	Msg(strings.Repeat(s, rN))
}

// TrimSpc :去除左右空格
func TrimSpc(rStr string) string {
	return strings.TrimSpace(rStr)
}

// Msg :打印到屏幕
func Msg(v ...interface{}) {
	fmt.Println(v...)
}

// GetVarType :获取变量类型
func GetVarType(rVal interface{}) string {
	return fmt.Sprint(reflect.TypeOf(rVal))
}

// CkErr :处理错误(如果有错误，返回true,无错则返回false),同时记录日志
func CkErr(rMsgTitle string, rErr error) bool {
	if rErr != nil {
		logs.LogMain.Error(rMsgTitle, rErr)
		return true
	}
	return false
}

// ConvertStrToErr ：转换string内容为错误对象
func ConvertStrToErr(s string) error {
	return fmt.Errorf("%s", s)
}

// ConvertBytesOfAscToStr :转换ascii为字符串，可用于http前端数据传至后端
func ConvertBytesOfAscToStr(rMetBt []byte) string {
	var s string
	for i := 0; i < len(rMetBt); i++ {
		s += fmt.Sprintf("%c", rMetBt[i])
	}
	return s
}

// ConvertBasNumberToStr :把数值类型数据（仅支持int/uint）转换成指定进制数值，返回字符串
func ConvertBasNumberToStr(rBase int, rVal interface{}) string {
	//	reflect.TypeOf(ref_varContent)
	var strBaseValue string
	switch rBase {
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
	strRetValue := fmt.Sprintf("%"+strBaseValue, rVal)
	return strRetValue
}

// ConvertOxBytesToStr :转换Bytes为十六进制字符串 (此方法没有用HEX包)
//转换过程中可能会丢失0，因此需要补0
func ConvertOxBytesToStr(rCnt []byte) string {
	var strRet string
	for i := 0; i < len(rCnt); i++ {
		strCon := ConvertBasNumberToStr(16, rCnt[i])
		if len(strCon) == 1 {
			strCon = "0" + strCon
		}
		strRet += strCon
	}
	return strRet
}

// ConvertBasStrToInt :把指定进制的字符转换成为十进制数值（int型）
// 请注意，只能还原数值的进制，并不能转换进制
func ConvertBasStrToInt(rBase int, rStrCnt string) int {
	intRetValue, err := strconv.ParseInt(rStrCnt, rBase, 64)
	if err != nil {
		intRetValue = 0
	}
	return int(intRetValue)
}

// ConvertBasStrToUint :根据指定进制要求，把字符串转换成数字Uint8型
//  请注意，只能还原数值的进制，并不能转换进制
func ConvertBasStrToUint(rBase int, rStrCnt string) uint8 {
	intRetValue, err := strconv.ParseUint(rStrCnt, rBase, 64)
	if err != nil {
		intRetValue = 0
	}
	return uint8(intRetValue)
}

// ConvertHexUnitToDecUnit :16进制数字字符转换成10进制unit
func ConvertHexUnitToDecUnit(rData uint8) uint8 {
	return ConvertBasStrToUint(10, ConvertBasNumberToStr(16, rData))
}

// ConvertStrIPAddToBytes :把指定十进制IP地址转换成为bytes
func ConvertStrIPAddToBytes(rIP string) []byte {
	st := strings.SplitN(rIP, ".", 4)
	var bs []byte
	for i := 0; i < len(st); i++ {
		mv, _ := strconv.ParseUint(st[i], 10, 64)
		t := ConvertBasStrToUint(10, ConvertBasNumberToStr(10, mv))
		bs = append(bs, byte(t))
	}
	return bs
}

// ConvertEvenDecToBytes :偶数十进制数值转换为bytes
// 请注意：只支持偶数位
func ConvertEvenDecToBytes(rStrCnt string) []byte {
	var bs []byte
	ms := ConvertBasNumberToStr(16, ConvertBasStrToInt(10, rStrCnt))
	if len(ms) > 2 {
		t := ConvertStrToBytesByPerTwoChar(ms)
		for i := 0; i < len(t); i++ {
			bs = append(bs, t[i])
		}
	}
	return bs
}

// ConvertDecToBytes :十进制数值转换为bytes
func ConvertDecToBytes(rValue int64) []byte {
	var rbs []byte
	mbs := ConvertIntToBytes(rValue)
	for i := 0; i < len(mbs); i++ {
		if mbs[i] != 0 {
			rbs = append(rbs, mbs[i])
		}
	}
	return rbs
}

// ConvertBasPartOfStrToBytes :根据开始、结束下标返回相应的字符串内容返回bytes
func ConvertBasPartOfStrToBytes(rStrCnt string, rBegin int, rEnd int, rBase int) []byte {
	var bT []byte
	n := len(rStrCnt)
	if rEnd <= n && rBegin >= 0 {
		for i := rBegin; i <= rEnd; i++ {
			strT := string(rStrCnt[i])
			bT = append(bT, ConvertBasStrToUint(rBase, strT))
		}
	}
	return bT
}

//GetPartOfStr :根据开始、结束下标返回相应的字符串内容返回string
func GetPartOfStr(rStrCnt string, rBegin int, rEnd int) string {
	var strR string
	n := len(rStrCnt)
	if rEnd < n && rBegin >= 0 {
		for i := rBegin; i <= rEnd; i++ {
			strR += string(rStrCnt[i])
		}
	}
	return strR
}

// ConvertByteToBinaryOfBytes :转换byte内的数据为二进制的byte切片
func ConvertByteToBinaryOfBytes(rByte byte) []byte {
	var bT []byte
	s := ConvertBasNumberToStr(2, rByte)
	for i := 0; i < len(s); i++ {
		t, _ := strconv.ParseUint(string(s[i]), 10, 64)
		tt := uint8(t)
		bT = append(bT, tt)
	}
	return bT
}

// ConvertBytesToInt :字节转换成整形
func ConvertBytesToInt(n []byte) int {
	bytesbuffer := bytes.NewBuffer(n)
	var x int32
	binary.Read(bytesbuffer, binary.BigEndian, &x)

	return int(x)
}

// ConvertStrToBytesByPerTwoChar :把字符串内容按每两字符对应一个byte组成新的bytes，返回[]byte
// 注意：目前只支持偶数位字符转换
func ConvertStrToBytesByPerTwoChar(rStrCnt string) []byte {
	var bT []byte
	n := len(rStrCnt)
	i := 0
	for i < n-1 {
		j := i + 1
		strP := string(rStrCnt[i]) + string(rStrCnt[j])
		bT = append(bT, ConvertBasStrToUint(16, string(strP)))
		i = j + 1
	}
	return bT
}

// IsExists :判断是否存在,true表示存在
func IsExists(rPath string) bool {
	_, err := os.Stat(rPath)
	if err != nil {
		return false
	}
	return true
}

// GetCurrentDirectory :获取当然路径
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

// GetPathSeparator :根据操作系统环境，返回分隔符
func GetPathSeparator() string {
	var strPath string
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		strPath = "\\"
	} else {
		strPath = "/"
	}
	return strPath
}

// ConvertIntToStr :转化整形成字符型
func ConvertIntToStr(intContent int) string {
	return strconv.Itoa(intContent)
}

// ConvertIntToBytes :整形转换成bytes
// 注意：64位数字，可能会产生几位0的数值byte
func ConvertIntToBytes(n int64) []byte {
	tmp := int64(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	// binary.Write(bytesBuffer, binary.LittleEndian, tmp)
	return bytesBuffer.Bytes()
}

// WrToFilWithBuffer :写入指定文件，如果没有该文件自动生成
func WrToFilWithBuffer(rFilePath string, rStrCnt string, rIsAppend bool) bool {
	//	this function for complex content to file
	var intFileOpenMode int
	if rIsAppend {
		intFileOpenMode = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	} else {
		intFileOpenMode = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	fileHandle, err := os.OpenFile(rFilePath, intFileOpenMode, 0666)
	defer fileHandle.Close()
	if CkErr("获取文件句柄失败！", err) {
		fileHandle.WriteString("\r\n" + rStrCnt)
	} else {
		return false
	}
	return true
}
