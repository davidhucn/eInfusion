package ntcp

import "fmt"

// PacketHeader :TCP包头
type PacketHeader struct {
	length             int    // 头数据长度(涵盖包总长度的帧头长度)
	prefixData         []byte // 桢头数据前缀内容
	packetLengthCursor int    //总长度下标
}

// NewTCPHeader :创建TCP数据包头对象
func NewTCPHeader(headerLength int, headerPrefixContent []byte, packetLengthCursor int) *PacketHeader {
	if packetLengthCursor >= headerLength {
		fmt.Println("错误，帧头数据内定义包总长度的下标非法！")
		return nil
	}
	return &PacketHeader{
		length:             headerLength,
		prefixData:         headerPrefixContent,
		packetLengthCursor: packetLengthCursor,
	}
}
