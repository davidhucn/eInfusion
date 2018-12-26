package ntcp

import (
	"bytes"
	cm "eInfusion/comm"
	"eInfusion/tlogs"
	"errors"
	"net"
	"sync"
	"time"
)

// Client :客户端对象-holds info about connection
type Client struct {
	conn   *net.TCPConn
	server *TServer
}

// NewClient :新建客户端对象
func newClient(c *net.TCPConn, s *TServer) *Client {
	return &Client{
		conn:   c,
		server: s,
	}
}

// Read client data from channel
func (c *Client) listen() {
	connID := cm.GetPureIPAddr(c.conn) // TCP连接ID，ip地址
	c.server.onNewClientCallback(c)
	c.server.clients.Store(connID, c.conn) // 添加到服务器客户端映射中client id 是string
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
		headerPrefixLength := len(c.server.packetHeader.packetPrefix)
		// 判断头数据是否正确
		if !bytes.Equal(headerBuffer[:headerPrefixLength], c.server.packetHeader.packetPrefix) {
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

// SendData :把指令存储到sendQueue，发送队列中，系统将自动发送
func (c *Client) SendData(b []byte) {
	od := cm.NewCmd(cm.GetPureIPAddr(c.conn), b)
	c.server.sendQueue <- od
}

// VerifyLegal :审核客户端为合法
func (c *Client) VerifyLegal() bool {
	// TODO:
	return true
}

// Boradcast :广播到所有客户端
func (s *TServer) Boradcast(b []byte) error {
	s.clients.Range(func(cID, conn interface{}) bool {
		od := cm.NewCmd(cID.(string), b)
		s.sendQueue <- od
		return true
	})
	return nil
}

// TServer :tcpserver
type TServer struct {
	listenPort               string // Address to open connection:
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewDataReceived        func(c *Client, p []byte)
	clients                  sync.Map //[ID:string]*Client
	readExpireTime           time.Duration
	packetHeader             *PacketHeader
	sendQueue                chan *cm.Cmd // 发送队列
	waitQueue                []*cm.Cmd    //等待队列，下线客户端([ID:string]*cmd)
	waitSendExpireTime       time.Duration
}

// NewTCPServer :Creates new tcp tcpserver instance
func NewTCPServer(listenPort string, readExpireTime time.Duration, waitSendExpireTime time.Duration, packetHeader *PacketHeader) *TServer {
	s := &TServer{
		listenPort:         listenPort,
		readExpireTime:     readExpireTime,
		packetHeader:       packetHeader,
		sendQueue:          make(chan *cm.Cmd, 2048),
		waitQueue:          make([]*cm.Cmd, 0),
		waitSendExpireTime: waitSendExpireTime,
	}
	s.WhenNewClientConnected(func(c *Client) {})
	s.WhenNewDataReceived(func(c *Client, p []byte) {})
	s.WhenClientConnectionClosed(func(c *Client, err error) {})
	return s
}

// GetClient :获取指定Client对象
func (s *TServer) GetClient(id string) *Client {
	if id != "" {
		c, ok := s.clients.Load(id)
		if ok {
			return c.(*Client)
		}
	}
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
	if cm.CkErr(TCPMsg.SourceError, tlogs.Error, err) {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if cm.CkErr(TCPMsg.SourceError, tlogs.Error, err) {
		panic(err)
	}
	cm.SepLi(60, "")
	cm.Msg(TCPMsg.StartServiceMsg, ",监听地址：")
	cm.SepLi(60, "")
	defer listener.Close()

	// 循环发送列表内指令
	go func() {
		for od := range s.sendQueue {
			c, ok := s.clients.Load(od.ID) //此处ID为IP地址
			if ok {
				time.Sleep(15 * time.Millisecond)
				_, err := c.(*net.TCPConn).Write(od.Data)
				if !cm.CkErr(TCPMsg.SendError, tlogs.Error, err) {
					// 发送成功
					s.clients.Delete(od.ID)
				} else {
					s.waitQueue = append(s.waitQueue, od)
				}
			}
		}
	}()
	// 校验发送列表sendQueue和待发送队列waitQueue的时间，如果超过指令时间，则清除
	// 循环清除超过指定时间周期的【待发列表】
	go func() {
		for i, v := range s.waitQueue {
			//判断指令是否在生存期内
			if time.Now().Sub(v.CreateTime) >= s.waitSendExpireTime {
				// 超时:册除相应待发队列
				s.waitQueue = append(s.waitQueue[:i], s.waitQueue[i+1])
				// TODO:错误信息回写到前端
			}
		}
	}()

	for {
		c, _ := listener.AcceptTCP()
		client := newClient(c, s)
		client.conn.SetReadDeadline(time.Now().Add(s.readExpireTime))
		go client.listen()
	}
}
