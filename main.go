package main

import (
	"eInfusion/ndb"
	"eInfusion/ntcp"
	"runtime"
)

func init() {
	// 初始化日志
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// go trsfscomm.StartCreateQRCode()
	// wc := eh.NewWebClients()
	// go eh.StartHTTPServer(wc, 7779)
	// tcpSer := et.NewTCPServer(300, 45)
	// et.RunTCPService(tcpSer, 7778)
	ndb.InitDB()
	ntcp.StartTCPService()
	// var s sync.Map
	// s.Store("g", 97)
	// s.Store("l", 100)
	// s.Store("e", 200)

	// s.Range(func(k, v interface{}) bool {
	// 	// fmt.Println("iterate:", k, v)
	// 	if k == "l" {
	// 		fmt.Print(v)
	// 	}
	// 	return true
	// })

}

// refereSample := func() {
// 	c := make(chan int, 10)
// 	go func() {
// 		for v := range c {;
// 			if v == 8 {
// 				print("bingo!")
// 			} else {
// 				println("waiting!")
// 			}
// 		}
// }()

// 	println("start...")
// 	for i := 0; i < 9; i++ {
// 		// a = append(a, i)
// 		c <- i
// 		println("len:", len(c))
// 		time.Sleep(2 * time.Second)
// 	}
// }
