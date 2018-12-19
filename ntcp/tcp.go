package ntcp

import (
	"bytes"
	cm "eInfusion/comm"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// Client holds info about connection
type Client struct {
	conn   net.Conn
	server *TServer
}

// TServer :tcpserver
type TServer struct {
	address                  string // Address to open connection: localhost:9999
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewDataReceived        func(c *Client, p []byte)
	clients                  sync.Map //[ID:string]*Client
	timeOutDuration          time.Duration
	packetHeader             *PacketHeader
	sendQueue                sync.Map //[ID:string][]byte
}

// NewTCPServer :Creates new tcp tcpserver instance
func NewTCPServer(address string, timeOutDuration time.Duration, packetHeader *PacketHeader) *TServer {
	fmt.Println(TCPMsg.StartServiceMsg, address)
	s := &TServer{
		address:         address,
		timeOutDuration: timeOutDuration,
		packetHeader:    packetHeader,
	}
	s.WhenNewClientConnected(func(c *Client) {})
	s.WhenNewDataReceived(func(c *Client, p []byte) {})
	s.WhenClientConnectionClosed(func(c *Client, err error) {})
	return s
}

// Read client data from channel
func (c *Client) listen() {
	connID := cm.GetRandString(10) // TCP连接ID，获取随机字符串
	c.server.onNewClientCallback(c)
	c.server.clients.Store(connID, c) // 添加到服务器客户端映射中client id 是string
	defer c.conn.Close()
	// 数据头包总长度
	headerLength := c.server.packetHeader.length
	headerBuffer := make([]byte, headerLength)
	for {
		_, err := c.conn.Read(headerBuffer)
		if err != nil {
			c.conn.Close()
			c.server.clients.Delete(connID)
			c.server.onClientConnectionClosed(c, err)
			return
		}
		// 头数据包前缀校验数据长度
		headerPrefixLength := len(c.server.packetHeader.prefixData)
		// 判断头数据是否正确
		if !bytes.Equal(headerBuffer[:headerPrefixLength], c.server.packetHeader.prefixData) {
			// 接收头数据包内数据不符合规定，则下线
			c.conn.Close()
			c.server.clients.Delete(connID)
			c.server.onClientConnectionClosed(c, errors.New(TCPMsg.HeaderDataError))
			return
		}
		pLenCursor := c.server.packetHeader.packetLengthCursor // 数据包总长度
		packetLength := cm.ConvertHexUnitToDecUnit(headerBuffer[pLenCursor])
		if packetLength >= 128 { // 数据包长度限定
			c.conn.Close()
			c.server.clients.Delete(connID)
			c.server.onClientConnectionClosed(c, errors.New(TCPMsg.ReceiveDataOutOfRange))
			return
		}
		// 开始接收内容
		packetDataLength := packetLength - headerLength
		dataBuffer := make([]byte, packetDataLength)
		c.conn.Read(dataBuffer)
		p := make([]byte, 0)
		for _, h := range headerBuffer {
			p = append(p, h)
		}
		for _, d := range dataBuffer {
			p = append(p, d)
		}
		c.server.onNewDataReceived(c, p)
	}
}

// SendBytes :发送bytes到客户端
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

// Boradcast :广播到所有客户端
// TODO:完成广播
func (s *TServer) Boradcast(b []byte) error {
	s.clients.Range(func(k, v interface{}) bool {
		fmt.Println(k, v)
		return true
	})
	return nil
}

// WhenNewClientConnected :当有新连接时
func (s *TServer) WhenNewClientConnected(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// WhenClientConnectionClosed :当连接下线时
func (s *TServer) WhenClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// WhenNewDataReceived :接收到新数据包
func (s *TServer) WhenNewDataReceived(callback func(c *Client, p []byte)) {
	s.onNewDataReceived = callback
}

// Listen :开始启动tcp服务
func (s *TServer) Listen() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", s.address)
	if cm.CkErr(TCPMsg.SourceError, err) {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if cm.CkErr(TCPMsg.SourceError, err) {
		panic(err)
	}
	cm.SepLi(60, "")
	cm.Msg(TCPMsg.StartServiceMsg, ",监听地址：", tcpAddr.IP.String())
	cm.SepLi(60, "")
	defer listener.Close()
	// 循环发送等待发送列表内指令
	go func() {
		s.sendQueue.Range(func(key, v interface{}) bool {
			tConn, ok := s.clients.Load(key)
			if ok {
				// TODO:TCP发送数据
			}
			return true
		})
	}()

	for {
		conn, _ := listener.Accept()
		client := &Client{
			conn:   conn,
			server: s,
		}
		client.conn.SetReadDeadline(time.Now().Add(s.timeOutDuration))
		go client.listen()
	}
}
