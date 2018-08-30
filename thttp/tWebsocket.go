package thttp

import (
	cm "eInfusion/comm"
	eq "eInfusion/tqueue"
	"net/http"
	"sync"
	"time"

	ws "github.com/gorilla/websocket"
)

var wsupgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ws接收目前仅限于json
func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	defer conn.Close()
	// cm.Msg("the type of conn:", cm.GetVarType(conn))
	if cm.CkErr("链接websocket出错！", err) {
		return
	}
	for {
		err := conn.ReadJSON(&clisData)
		// _, m, err := conn.ReadMessage()
		if cm.CkErr("websocket接收前端数据出错!", err) {
			// FIXME:制定通讯标准，此处应返回前端页面出错信息
			conn.WriteMessage(ws.TextMessage, []byte("can't Exchange the data"))
			break
		}
		// 根据前端应用需求信息发送指令
		// 获取时间戳，来生成orders id
		ssn := cm.GetTimeStamp()
		// 加入发送消息队列
		eq.AddToSendQueue(ssn, clisData.ID, cm.ConvertBasStrToUint(10, clisData.CmdType), clisData.Args)

		//定时检查返回数据池里面有否相应的数据
		go func() {
			var ms sync.Mutex
			cT := time.NewTicker(500 * time.Millisecond) // 定时
			select {
			case <-cT.C:
				if _, ok := eq.RcMsgs[ssn]; ok {
					// 回传前端
					conn.WriteJSON([]byte(eq.RcMsgs[ssn]))
					// 册除该条数据
					ms.Lock()
					defer ms.Unlock()
					delete(eq.RcMsgs, ssn)
				}
			}
		}()
		// conn.WriteMessage(1, []byte(cliMsg[0].Action))
	}

}
