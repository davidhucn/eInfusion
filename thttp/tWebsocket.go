package thttp

import (
	cm "eInfusion/comm"
	eq "eInfusion/tqueue"
	"net/http"
	"sync"
	"time"

	ws "github.com/gorilla/websocket"
	// // 仅仅应用里面的websocket对象S
	// _ "golang.org/x/net/websocket"
)

var wsupgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// SendBackWS :回写到前端
func SendBackWS(rStrCnt string) {
	wsConn.Write([]byte(rStrCnt))

}

// ws接收目前仅限于json
func wshandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := wsupgrader.Upgrade(w, r, nil)
	defer wsConn.Close()
	// cm.Msg("the type of wsConn:", cm.GetVarType(wsConn))
	if cm.CkErr("链接websocket出错！", err) {
		return
	}
	for {
		err := wsConn.ReadJSON(&clisData)
		// _, m, err := wsConn.ReadMessage()
		if cm.CkErr("websocket接收前端数据出错!", err) {
			// FIXME:制定通讯标准，此处应返回前端页面出错信息
			wsConn.WriteMessage(ws.TextMessage, []byte("can't Exchange the data"))
			break
		}
		// 根据前端应用需求信息发送指令
		// 获取时间戳，来生成orders id
		ssn := cm.GetRandString(10)
		//定时检查返回数据池里面有否相应的数据
		go func() {
			var ms sync.Mutex
			cT := time.NewTicker(500 * time.Millisecond) // 定时
			select {
			case <-cT.C:
				cm.Msg("i am alive!!")
				if _, ok := eq.RcMsgs[ssn]; ok {
					// 回传前端
					wsConn.WriteMessage(ws.TextMessage, []byte(eq.RcMsgs[ssn]))
					// 册除该条数据
					ms.Lock()
					defer ms.Unlock()
					delete(eq.RcMsgs, ssn)
				}
			}
		}()
		// 加入发送消息队列
		for i := 0; i < len(clisData); i++ {
			eq.AddToSendQueue(ssn, clisData[i].ID, cm.ConvertBasStrToUint(10, clisData[i].CmdType), clisData[i].Args)
			wsConn.WriteMessage(ws.TextMessage, []byte("Doing..."))
		}

		// wsConn.WriteMessage(1, []byte(cliMsg[0].Action))
	}

}
