package tqueue

// StartTrsfusionQueue :启动队列处理平台
func StartQueueWait() {
	// pingTicker := time.NewTicker(10 * time.Second) // 定时
	// testAfter := time.After(5 * time.Second)       // 延时

	// for HTTPReqStream != nil {
	// 	select {
	// 	case <-HTTPReqStream:
	// 		cm.Msg("quit")
	// 		return
	// 	case <-pingTicker.C:
	// 		//发送心跳
	// 		// _, err := SendData(conn, []byte("PING"))
	// 		// if err != nil {
	// 		// 	pingTicker.Stop()
	// 		// 	return
	// 		// }
	// 	case <-testAfter:
	// 		//	doLog("testAfter:")
	// 		//TODO:日志记录
	// 	}
	// 	// select {
	// 	// case <-tick.C:
	// 	// fmt.Printf("%d: case <-tick.C\n", i)
	// 	// default:
	// }
}

func AddReqQueue(TargetID chan string, CmdType chan uint8, Args chan string) {

}
