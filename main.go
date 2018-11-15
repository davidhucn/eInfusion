package main

import (
	eh "eInfusion/http"
	et "eInfusion/tcp"
	logs "eInfusion/tlogs"
	"runtime"
)

func init() {
	// 初始化日志
	logs.LogDisable()
	logs.LogConfigLoad()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	wc := eh.NewWebClients()
	go eh.StartHTTPServer(wc, 7779)
	tcpSer := et.NewTCPServer(300, 45)
	et.RunTCPService(tcpSer, 7778)

	// refereSample()
}

// refereSample := func(){
// 	c := make(chan int, 10)
// 	go func() {
// 		for v := range c {
// 			if v == 8 {
// 				print("bingo!")
// 			} else {
// 				println("waiting!")
// 			}
// 		}
// 	}()

// 	println("start...")
// 	for i := 0; i < 9; i++ {
// 		// a = append(a, i)
// 		c <- i
// 		println("len:", len(c))
// 		time.Sleep(2 * time.Second)
// 	}
// }
