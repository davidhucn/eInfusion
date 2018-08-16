package thttp

import (
	cm "eInfusion/comm"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ws接收目前仅限于json
func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	defer conn.Close()
	if cm.CkErr("链接websocket出错！", err) {
		return
	}
	for {
		err := conn.ReadJSON(&cliMsg)
		// _, m, err := conn.ReadMessage()
		if cm.CkErr("websocket接收前端数据出错!", err) {
			// FIXME:制定通讯标准，此处应返回前端页面出错信息
			conn.WriteMessage(1, []byte("can't Exchange the data"))
			break
		}
		// 根据前端应用需求信息发送指令
		for i := 0; i < len(cliMsg); i++ {
			if verifyReqAction(cliMsg[i]) {

			}
		}
		// 回传前端
		conn.WriteJSON(cliMsg)

		// conn.WriteMessage(1, []byte(cliMsg[0].Action))
	}

}
