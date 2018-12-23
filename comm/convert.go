package comm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

// ConvertBasStrToInt64 :把指定进制的字符转换成为十进制数值（int型）
// 请注意，只能还原数值的进制，并不能转换进制
func ConvertBasStrToInt64(rBase int, rStrCnt string) int64 {
	intRetValue, err := strconv.ParseInt(rStrCnt, rBase, 64)
	if err != nil {
		intRetValue = 0
	}
	return intRetValue
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

// ConvertStrIPAddrToBytes :把指定十进制IP地址转换成为bytes
func ConvertStrIPAddrToBytes(rIP string) []byte {
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
	ms := ConvertBasNumberToStr(16, ConvertBasStrToInt64(10, rStrCnt))
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

// ConvertTimeToStr : 转换Time为String（指定格式：小时：分钟：秒）
func ConvertTimeToStr(rTime time.Time) string {
	return rTime.Format("2006-01-02 15:04:05")
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
