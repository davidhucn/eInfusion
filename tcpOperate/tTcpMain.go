package tcpOperate

import (
	"eInfusion/comm"
	"eInfusion/logs"
	ep "eInfusion/protocol"
	"net"
	"os"
	"time"
)

//echo server Goroutine
func EchoFunc(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			//println("Error reading:", err.Error())
			return
		}
		//send reply
		_, err = conn.Write(buf)
		if err != nil {
			//println("Error send reply:", err.Error())
			return
		}
	}
}

//initial listener and run
func StartTcpServer() {
	listener, err := net.Listen("tcp", ":"+c_TcpServer_Port)
	defer listener.Close()
	if err != nil {
		fmt.Println("error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("running ...\n")

	var cur_conn_num int = 0
	conn_chan := make(chan net.Conn)
	ch_conn_change := make(chan int)

	go func() {
		for conn_change := range ch_conn_change {
			cur_conn_num += conn_change
		}
	}()

	go func() {
		for _ = range time.Tick(1e8) {
			fmt.Printf("cur conn num: %f\n", cur_conn_num)
		}
	}()

	for i := 0; i < c_MaxConnectionAmount; i++ {
		go func() {
			for conn := range conn_chan {
				ch_conn_change <- 1
				println("change:", ch_conn_change)
				EchoFunc(conn)
				ch_conn_change <- -1
				println("change:", ch_conn_change)
			}
		}()
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Error accept:", err.Error())
			return
		}
		conn_chan <- conn
	}
}
