package ntcp

import (
	"bytes"
	cm "eInfusion/comm"
	"errors"
	"net"
	"sync"
	"time"
)

// Client holds info about connection
type Client struct {
	conn   *net.TCPConn
	server *TServer
}

// TServer :tcpserver
type TServer struct {
	listenPort               string // Address to open connection: localhost:9999
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewDataReceived        func(c *Client, p []byte)
	clients                  sync.Map //[ID:string]*Client
	timeOutDuration          time.Duration
	packetHeader             *PacketHeader
	sendQueue                chan *cm.Cmd // 发送队列
	waitQueue                sync.Map     //等待队列，下线客户端
}

// NewTCPServer :Creates new tcp tcpserver instance
func NewTCPServer(listenPort string, timeOutDuration time.Duration, packetHeader *PacketHeader) *TServer {
	s := &TServer{
		listenPort:      listenPort,
		timeOutDuration: timeOutDuration,
		packetHeader:    packetHeader,
		sendQueue:       make(chan *cm.Cmd, 1024),
	}
	s.WhenNewClientConnected(func(c *Client) {})
	s.WhenNewDataReceived(func(c *Client, p []byte) {})
	s.WhenClientConnectionClosed(func(c *Client, err error) {})
	return s
}

// Read client data from channel
func (c *Client) listen() {
	connID := cm.GetPureIPAddr(c.conn) // TCP连接ID，ip地址
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

// SendData :发送bytes到客户端
func (c *Client) SendData(b []byte) error {
	// 判断是否客户端在线，在线则发送
	// TODO:
	_, err := c.conn.Write(b)
	// od := cm.NewCmd(c.conn
	// 	), b)
	// c.server.clients.Load()
	return err
}

// Boradcast :广播到所有客户端
func (s *TServer) Boradcast(b []byte) error {
	s.clients.Range(func(cID, conn interface{}) bool {
		// TODO:完成广播

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
	tcpAddr, err := net.ResolveTCPAddr("tcp4", s.listenPort)
	if cm.CkErr(TCPMsg.SourceError, err) {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if cm.CkErr(TCPMsg.SourceError, err) {
		panic(err)
	}
	cm.SepLi(60, "")
	cm.Msg(TCPMsg.StartServiceMsg, ",监听地址：")
	cm.SepLi(60, "")
	defer listener.Close()

	// 循环发送列表内指令
	go func() {
		for od := range s.sendQueue {
			c, ok := s.clients.Load(od.ID)
			if ok {
				time.Sleep(15 * time.Millisecond)
				_, err := c.(*net.TCPConn).Write(od.Data)
				if !cm.CkErr(TCPMsg.SendError, err) {
					// 发送成功
					s.clients.Delete(od.ID)
				} else {
					// TODO:发送不成功,存储至待发送列表内

				}
			}
		}
	}()
	for {
		c, _ := listener.AcceptTCP()
		client := &Client{
			conn:   c,
			server: s,
		}
		client.conn.SetReadDeadline(time.Now().Add(s.timeOutDuration))
		go client.listen()
	}
}
