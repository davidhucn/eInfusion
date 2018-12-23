package comm

import (
	lg "eInfusion/tlogs"
	"fmt"
	"math/rand"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Cmd :指令对象
type Cmd struct {
	Data       []byte
	ID         string
	CreateTime time.Time
}

// NewCmd :生成新的命令对象
func NewCmd(rID string, rCnt []byte) *Cmd {
	return &Cmd{
		Data:       rCnt,
		ID:         rID,
		CreateTime: time.Now(),
	}
}

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
func GetPureIPAddr(rConn *net.TCPConn) string {
	ip := rConn.RemoteAddr().String()
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
func CkErr(rMsgTitle string, errType lg.LogType, rErr error) bool {
	if rErr != nil {
		if rMsgTitle != "" {
			lg.DoLog(errType, rMsgTitle, rErr)
		}
		return true
	}
	return false
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
	if CkErr("获取文件句柄失败！", lg.Error, err) {
		fileHandle.WriteString("\r\n" + rStrCnt)
	} else {
		return false
	}
	return true
}
