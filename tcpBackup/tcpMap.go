package tcpBackup

//import (
//	"bytes"
//	"encoding/binary"
//	"fmt"
//	"log"
//	"net"
//	"os"
//	"strconv"
//	"strings"
//	"sync"
//	"time"
//)

//func main() {
//	tcpStart(8889)
//}

///*
//   定义相关锁
//*/
//var (
//	connMkMutex  sync.Mutex
//	connDelMutex sync.Mutex
//)

///*
//   定义logger
//*/
//var logger *log.Logger

///*
//   初始化log
//*/
//func initLog(logfile *os.File) {
//	// logger = log.New(logfile, "log:", log.Ldate|log.Ltime)
//	logger = log.New(logfile, "prefix ", 0)
//}

///*
//   处理log
//*/
//func doLog(args ...interface{}) {
//	str := time.Now().Format("2006-01-02 15:04:05")
//	var logData string
//	var temp string

//	for _, arg := range args {
//		switch val := arg.(type) {
//		case int:
//			temp = strconv.Itoa(val)
//		case string:
//			temp = val
//		}
//		if len(temp) > 64 { // 限制只打印前64个字符
//			logData = temp[:64]
//		} else {
//			logData = temp
//		}
//		str = str + " " + logData
//	}
//	logger.Println(str)
//}

///*
//   定义socket conn 映射
//*/
//var clisConnMap map[string]*net.TCPConn

///*
//   初始化socket conn 映射
//*/
//func initClisConnMap() {
//	clisConnMap = make(map[string]*net.TCPConn)
//}

///*
//   建立socket conn 映射
//*/
//func mkClisConn(key string, conn *net.TCPConn) {
//	connMkMutex.Lock()
//	defer connMkMutex.Unlock()
//	clisConnMap[key] = conn
//}

///*
//   删除socket conn 映射
//*/
//func delClisConn(key string) {
//	connDelMutex.Lock()
//	defer connDelMutex.Unlock()
//	delete(clisConnMap, key)
//}

///*
//   定义解码器
//*/
//type Unpacker struct {
//	// 头(xy)2bytes + 标识1byte + 包长度2bytes + data
//	// 当然了,头不可能是xy,这里举例子,而且一般还需要转义
//	_buf []byte
//}

//func (unpacker *Unpacker) feed(data []byte) {
//	unpacker._buf = append(unpacker._buf, data...)
//}

//func (unpacker *Unpacker) unpack() (flag byte, msg []byte) {
//	str := string(unpacker._buf)
//	for {
//		if len(str) < 5 {
//			break
//		} else {
//			_, head, data := Partition(str, "xy")
//			if len(head) == 0 { // 没有头
//				if str[len(str)-1] == byte(120) { // 120 => 'x'
//					unpacker._buf = []byte{byte(120)}
//				} else {
//					unpacker._buf = []byte{}
//				}
//				break
//			}

//			buf := bytes.NewReader([]byte(data))
//			msg = make([]byte, buf.Len())
//			var dataLen uint16
//			binary.Read(buf, binary.LittleEndian, &flag)
//			binary.Read(buf, binary.LittleEndian, &dataLen)

//			fmt.Println("DEC:", flag, dataLen)
//			if buf.Len() < int(dataLen) {
//				break
//			}
//			binary.Read(buf, binary.LittleEndian, &msg)
//			unpacker._buf = unpacker._buf[2+1+2+dataLen:]
//		}
//	}
//	return
//}

///*
//   启动服务
//*/
//func tcpStart(port int) {
//	initLog(os.Stderr)
//	initClisConnMap()
//	doLog("tcpStart:")
//	host := ":" + strconv.Itoa(port)
//	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
//	checkError(err)

//	listener, err := net.ListenTCP("tcp", tcpAddr)
//	checkError(err)

//	for {
//		conn, err := listener.AcceptTCP()
//		if err != nil {
//			continue
//		}

//		go handleClient(conn)
//	}
//}

///*
//   socket conn
//*/
//func handleClient(conn *net.TCPConn) {
//	// ****这里是初始化连接处理
//	addr := conn.RemoteAddr().String()
//	doLog("handleClient:", addr)
//	connectionMade(conn)
//	request := make([]byte, 128)
//	defer conn.Close()
//	buf := make([]byte, 0)
//	for {
//		// ****这里是读循环处理
//		readLoopHandled(conn)
//		read_len, err := conn.Read(request)
//		if err != nil {
//			// 这里没使用checkError因为不退出,只是break出去
//			doLog("ERR:", "read err", err.Error())
//			break
//		}

//		if read_len == 0 { // 在gprs时数据不能通过这个判断是否断开连接,要通过心跳包
//			doLog("ERR:", "connection already closed by client")
//			break
//		} else {
//			// request[:read_len]处理
//			buf = append(buf, request[:read_len]...)
//			doLog("<=", addr, string(request[:read_len]))
//			dataReceived(conn, &buf)
//			request = make([]byte, 128) // clear last read content
//		}
//	}
//	// ****这里是连接断开处理
//	connectionLost(conn)
//}

///*
//   连接初始处理(ed)
//*/
//func connectionMade(conn *net.TCPConn) {
//	//初始化连接这个函数被调用

//	// ****建立conn映射
//	addr := conn.RemoteAddr().String()
//	ip := strings.Split(addr, ":")[0]
//	mkClisConn(ip, conn)

//	doLog("connectionMade:", addr)

//	// ****定时处理(心跳等)
//	go loopingCall(conn)
//}

///*
//   读循环处理(ed)
//*/
//func readLoopHandled(conn *net.TCPConn) {
//	//当进入循环读数据这个函数被调用, 主要用于设置超时(好刷新设置超时)

//	// *****设置超时 (要写在for循环里)
//	setReadTimeout(conn, 10*time.Minute)
//}

///*
//   客户端连接发送来的消息处理(ed)
//*/
//func dataReceived(conn *net.TCPConn, pBuf *[]byte) {
//	//一般情况可以用pBuf参数,但是如果有分包粘包的情况就必须使用clisBufMap的buf
//	//clisBufMap的buf不断增大,不管是否使用都应该处理
//	//addr := conn.RemoteAddr().String()
//	doLog("*pBuf:", string(*pBuf))
//	//sendData(clisConnMap["192.168.6.234"], []byte("xxx"))
//	sendData(conn, []byte("echo"))
//}

///*
//   连接断开(ed)
//*/
//func connectionLost(conn *net.TCPConn) {
//	//连接断开这个函数被调用
//	addr := conn.RemoteAddr().String()
//	ip := strings.Split(addr, ":")[0]

//	delClisConn(ip) // 删除关闭的连接对应的clisMap项
//	doLog("connectionLost:", addr)
//}

///*
//   发送数据
//*/
//func sendData(conn *net.TCPConn, data []byte) (n int, err error) {
//	addr := conn.RemoteAddr().String()
//	n, err = conn.Write(data)
//	if err == nil {
//		doLog("=>", addr, string(data))
//	}
//	return
//}

///*
//   广播数据
//*/
//func broadcast(tclisMap map[string]*net.TCPConn, data []byte) {
//	for _, conn := range tclisMap {
//		sendData(conn, data)
//	}
//}

///*
//   定时处理&延时处理
//*/
//func loopingCall(conn *net.TCPConn) {
//	pingTicker := time.NewTicker(30 * time.Second) // 定时
//	testAfter := time.After(5 * time.Second)       // 延时

//	for {
//		select {
//		case <-pingTicker.C:
//			//发送心跳
//			_, err := sendData(conn, []byte("PING"))
//			if err != nil {
//				pingTicker.Stop()
//				return
//			}
//		case <-testAfter:
//			doLog("testAfter:")
//		}
//	}
//}

///*
//   设置读数据超时
//*/
//func setReadTimeout(conn *net.TCPConn, t time.Duration) {
//	conn.SetReadDeadline(time.Now().Add(t))
//}

///*
//   错误处理
//*/
//func checkError(err error) {
//	if err != nil {
//		doLog("ERR:", err.Error())
//		os.Exit(1)
//	}
//}

//func Partition(s string, sep string) (head string, retSep string, tail string) {
//	// Partition(s, sep) -> (head, sep, tail)
//	index := strings.Index(s, sep)
//	if index == -1 {
//		head = s
//		retSep = ""
//		tail = ""
//	} else {
//		head = s[:index]
//		retSep = sep
//		tail = s[len(head)+len(sep):]
//	}
//	return
//}
