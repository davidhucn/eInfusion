package tqueue

// TargetID :目标设备序号
var TargetID chan string

// CmdType :操作指令类型
var CmdType chan uint8

// Args :命令 、 参数数据
var Args chan string

// SendCmdCnt :指令内容(bytes)
var SendCmdCnt chan []byte
