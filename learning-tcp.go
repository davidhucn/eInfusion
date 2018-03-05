package connection

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

//支持的最大消息长度
const maxLength int = 1<<32 - 1 // 4294967295

var (
	rHeadBytes = [4]byte{0, 0, 0, 0}
	wHeadBytes = [4]byte{0, 0, 0, 0}
	errMsgRead = errors.New("Message read length error")
	errHeadLen = errors.New("Message head length error")
	errMsgLen  = errors.New("Message length is no longer in normal range")
)
var connPool sync.Pool

//从对象池中获取一个对象,不存在则申明
func Newconnection(conn net.Conn) Conn {
	c := connPool.Get()
	if cnt, ok := c.(*connection); ok {
		cnt.rwc = conn
		return cnt
	}
	return &connection{rlen: 0, rwc: conn}
}

type Conn interface {
	Read() (r io.Reader, size int, err error)
	Write(p []byte) (n int, err error)
	Writer(size int, r io.Reader) (n int64, err error)
	RemoteAddr() net.Addr
	LocalAddr() net.Addr
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	Close() (err error)
}

//定义一个结构体,用来封装Conn.
type connection struct {
	rlen  int        //消息长度
	rwc   net.Conn   //原始的网络链接
	rlock sync.Mutex //Conn读锁
	wlock sync.Mutex //Conn写锁
}

//此方法用来读取头部消息.
func (self *connection) rhead() error {
	n, err := self.rwc.Read(rHeadBytes[:])
	if n != 4 || err != nil {
		if err != nil {
			return err
		}
		return errHeadLen
	}
	self.rlen = int(binary.BigEndian.Uint32(rHeadBytes[:]))
	return nil
}

//此方法用来发送头消息
func (self *connection) whead(l int) error {
	if l <= 0 || l > maxLength {
		return errMsgLen
	}
	binary.BigEndian.PutUint32(wHeadBytes[:], uint32(l))
	_, err := self.rwc.Write(wHeadBytes[:])
	return err
}

//头部消息解析之后.返回一个io.Reader借口.用来读取远端发送过来的数据
//封装成limitRead对象.来实现ioReader接口
func (self *connection) Read() (r io.Reader, size int, err error) {
	self.rlock.Lock()
	if err = self.rhead(); err != nil {
		self.rlock.Unlock()
		return
	}
	size = self.rlen
	r = limitRead{r: io.LimitReader(self.rwc, int64(size)), unlock: self.rlock.Unlock}
	return
}

//发送消息前先调用whead函数,来发送头部信息,然后发送body
func (self *connection) Write(p []byte) (n int, err error) {
	self.wlock.Lock()
	err = self.whead(len(p))
	if err != nil {
		self.wlock.Unlock()
		return
	}
	n, err = self.rwc.Write(p)
	self.wlock.Unlock()
	return
}

//发送一个流.必须指定流的长度
func (self *connection) Writer(size int, r io.Reader) (n int64, err error) {
	self.wlock.Lock()
	err = self.whead(int(size))
	if err != nil {
		self.wlock.Unlock()
		return
	}
	n, err = io.CopyN(self.rwc, r, int64(size))
	self.wlock.Unlock()
	return
}

func (self *connection) RemoteAddr() net.Addr {
	return self.rwc.RemoteAddr()
}

func (self *connection) LocalAddr() net.Addr {
	return self.rwc.LocalAddr()
}

func (self *connection) SetDeadline(t time.Time) error {
	return self.rwc.SetDeadline(t)
}

func (self *connection) SetReadDeadline(t time.Time) error {
	return self.rwc.SetReadDeadline(t)
}

func (self *connection) SetWriteDeadline(t time.Time) error {
	return self.rwc.SetWriteDeadline(t)
}

func (self *connection) Close() (err error) {
	err = self.rwc.Close()
	self.rlen = 0
	connPool.Put(self)
	return
}

type limitRead struct {
	r      io.Reader
	unlock func()
}

func (self limitRead) Read(p []byte) (n int, err error) {
	n, err = self.r.Read(p)
	if err != nil {
		self.unlock()
	}
	return n, err
}
