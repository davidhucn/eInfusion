package ntcp

// THeader :TCP包头
type THeader struct {
	length             int
	content            []byte
	contentStartCursor int
	packetLengthCursor int
}

// NewTCPHeader :创建TCP数据包头对象
func NewTCPHeader(headerLength int, headerContentStartCursor int, headerContent []byte, packetLengthCursor int) *THeader {
	return &THeader{
		length:             headerLength,
		content:            headerContent,
		contentStartCursor: headerContentStartCursor,
		packetLengthCursor: packetLengthCursor,
	}
}
