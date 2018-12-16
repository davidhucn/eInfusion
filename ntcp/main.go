package ntcp

import (
	cm "eInfusion/comm"
	"log"
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
	onNewDataReceive         func(c *Client, p []byte)
	clients                  sync.Map //[string]*Client
	timeOutDuration          time.Duration
	packetHeader             *THeader
}

// Read client data from channel
func (c *Client) listen() {
	connID := cm.GetRandString(8)
	c.server.onNewClientCallback(c)
	// 添加到服务器客户端映射中client id 是string
	c.server.clients.Store(connID, c)
	defer c.conn.Close()
	for {

		if err != nil {
			c.conn.Close()
			c.server.clients.Delete(connID)
			c.server.onClientConnectionClosed(c, err)
			return
		}
		c.server.onNewDataReceive(c, p)
	}
}

// VerifyPacketHeader :判断包头是否符合要求
func (c *Client) verifyPacketHeader() bool {
	length := c.server.packetHeader.length
	headerBuffer := make([]byte, length)
	c.conn.Read(headerBuffer)
	for b:=range 
}

// Send text message to client
func (c *Client) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// Called right after tcpserver starts listening new client
func (s *TServer) onNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *TServer) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *TServer) onNewDataReceive(callback func(c *Client, p []byte)) {
	s.onNewDataReceive = callback
}

// Start network tcpserver
func (s *TServer) Listen() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP tcpserver.")
	}
	defer listener.Close()

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

// NewTCPServer :Creates new tcp tcpserver instance
func NewTCPServer(address string, timeOutDuration time.Duration, packetHeader *THeader) *TServer {
	cm.Msg(TCPMsg.StartServiceMsg, ",监听地址：", address)
	s := &TServer{
		address:         address,
		timeOutDuration: timeOutDuration,
		packetHeader:    packetHeader,
	}

	s.onNewClient(func(c *Client) {})
	s.onNewDataReceive(func(c *Client, p []byte) {})
	s.OnClientConnectionClosed(func(c *Client, err error) {})

	return s
}
