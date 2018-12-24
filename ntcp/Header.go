package ntcp

import "fmt"

// PacketHeader :TCP包头
type PacketHeader struct {
	length             uint8  // 头数据长度(涵盖包总长度的帧头长度)
	packetPrefix       []byte // 桢头数据前缀数据,校验数据 (1~N)
	packetLengthCursor uint8  //总长度下标（0~N）
}

// NewTCPHeader :创建TCP数据包头对象
// {headerLength} :头数据总长度 (范围：1~N) ;
// {headerPrefixContent} :头数据包前缀 ;
// {packetLengthCursor} :头数据包存储数据包总长度的下标（范围：0~N）;
func NewTCPHeader(headerLength uint8, headerPrefixData []byte, packetLengthCursor uint8) *PacketHeader {
	if packetLengthCursor >= headerLength {
		fmt.Println("错误，帧头数据内定义包总长度的下标非法！")
		return nil
	}
	return &PacketHeader{
		length:             headerLength,
		packetPrefix:       headerPrefixData,
		packetLengthCursor: packetLengthCursor,
	}
}
