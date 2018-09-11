package refer

// import (
// 	"encoding/json"
// 	"log"
// 	"net"
// 	"strconv"
// 	"time"
// )

// const (
// 	proxy_timeout = 5
// 	proxy_server  = "127.0.0.1:1988"
// 	msg_length    = 1024
// )

// type Request struct {
// 	reqId      int
// 	reqContent string
// 	rspChan    chan<- string // writeonly chan
// }

// //store request map
// var requestMap map[int]*Request

// type Clienter struct {
// 	client  net.Conn
// 	isAlive bool
// 	SendStr chan *Request
// 	RecvStr chan string
// }

// func (c *Clienter) Connect() bool {
// 	if c.isAlive {
// 		return true
// 	} else {
// 		var err error
// 		c.client, err = net.Dial("tcp", proxy_server)
// 		if err != nil {
// 			return false
// 		}
// 		c.isAlive = true
// 		log.Println("connect to " + proxy_server)
// 	}
// 	return true
// }

// //send msg to upstream server
// func ProxySendLoop(c *Clienter) {

// 	//store reqId and reqContent
// 	senddata := make(map[string]string)
// 	for {
// 		if !c.isAlive {
// 			time.Sleep(1 * time.Second)
// 			c.Connect()
// 		}
// 		if c.isAlive {
// 			req := <-c.SendStr

// 			//construct request json string
// 			senddata["reqId"] = strconv.Itoa(req.reqId)
// 			senddata["reqContent"] = req.reqContent
// 			sendjson, err := json.Marshal(senddata)
// 			if err != nil {
// 				continue
// 			}

// 			_, err = c.client.Write([]byte(sendjson))
// 			if err != nil {
// 				c.RecvStr <- string("proxy server close...")
// 				c.client.Close()
// 				c.isAlive = false
// 				log.Println("disconnect from " + proxy_server)
// 				continue
// 			}
// 			//log.Println("Write to proxy server: " + string(sendjson))
// 		}
// 	}
// }

// //recv msg from upstream server
// func ProxyRecvLoop(c *Clienter) {
// 	buf := make([]byte, msg_length)
// 	recvdata := make(map[string]string, 2)
// 	for {
// 		if !c.isAlive {
// 			time.Sleep(1 * time.Second)
// 			c.Connect()
// 		}
// 		if c.isAlive {
// 			n, err := c.client.Read(buf)
// 			if err != nil {
// 				c.client.Close()
// 				c.isAlive = false
// 				log.Println("disconnect from " + proxy_server)
// 				continue
// 			}
// 			//log.Println("Read from proxy server: " + string(buf[0:n]))

// 			if err := json.Unmarshal(buf[0:n], &recvdata); err == nil {
// 				reqidstr := recvdata["reqId"]
// 				if reqid, err := strconv.Atoi(reqidstr); err == nil {
// 					req, ok := requestMap[reqid]
// 					if !ok {
// 						continue
// 					}
// 					req.rspChan <- recvdata["resContent"]
// 				}
// 				continue
// 			}
// 		}
// 	}
// }

// //one handle per request
// func handle(conn *net.TCPConn, id int, tc *Clienter) {

// 	data := make([]byte, msg_length)
// 	handleProxy := make(chan string)
// 	request := &Request{reqId: id, rspChan: handleProxy}

// 	requestMap[id] = request
// 	for {
// 		n, err := conn.Read(data)
// 		if err != nil {
// 			log.Println("disconnect from " + conn.RemoteAddr().String())
// 			conn.Close()
// 			delete(requestMap, id)
// 			return
// 		}
// 		request.reqContent = string(data[0:n])
// 		//send to proxy
// 		select {

// 		case tc.SendStr <- request:
// 		case <-time.After(proxy_timeout * time.Second):
// 			//proxyChan <- &Request{cancel: true, reqId: id}
// 			_, err = conn.Write([]byte("proxy server send timeout."))
// 			if err != nil {
// 				conn.Close()
// 				delete(requestMap, id)
// 				return
// 			}
// 			continue
// 		}

// 		//read from proxy
// 		select {
// 		case rspContent := <-handleProxy:
// 			_, err := conn.Write([]byte(rspContent))
// 			if err != nil {
// 				conn.Close()
// 				delete(requestMap, id)
// 				return
// 			}
// 		case <-time.After(proxy_timeout * time.Second):
// 			_, err = conn.Write([]byte("proxy server recv timeout."))
// 			if err != nil {
// 				conn.Close()
// 				delete(requestMap, id)
// 				return
// 			}
// 			continue
// 		}
// 	}
// }

// func TcpLotusMain(ip string, port int) {
// 	//start tcp server
// 	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), port, ""})
// 	if err != nil {
// 		log.Fatalln("listen port error")
// 		return
// 	}
// 	log.Println("start tcp server " + ip + " " + strconv.Itoa(port))
// 	defer listen.Close()

// 	//start proxy connect and loop
// 	var tc Clienter
// 	tc.SendStr = make(chan *Request, 1000)
// 	tc.RecvStr = make(chan string)
// 	tc.Connect()

// 	go ProxySendLoop(&tc)
// 	go ProxyRecvLoop(&tc)
// 	//listen new request
// 	requestMap = make(map[int]*Request)
// 	var id int = 0
// 	for {

// 		conn, err := listen.AcceptTCP()
// 		if err != nil {
// 			log.Println("receive connection failed")
// 			continue
// 		}
// 		id++
// 		log.Println("connected from " + conn.RemoteAddr().String())
// 		go handle(conn, id, &tc)

// 	}
// }
