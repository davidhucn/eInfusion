package thttp

import (
	cm "eInfusion/comm"
	eq "eInfusion/tqueue"
	"net/http"

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
	cm.Msg("timeout:", wsupgrader.HandshakeTimeout)
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
		for i := 0; i < len(clisData); i++ {
			// 获取时间戳，来生成orders
			ssn := cm.GetTimeStamp()
			// 加入发送消息队列
			eq.AddToSendQueue(ssn, clisData[i].ID, cm.ConvertBasStrToUint(10, clisData[i].CmdType), clisData[i].Args)
		}

		// TODO: 引入chan，接收返回指令结果，回写前端
		// go func() {
		// 	for cs := range connStream {

		// 	}
		// }()
		// 回传前端
		conn.WriteJSON(clisData)
		// conn.WriteMessage(1, []byte(cliMsg[0].Action))
	}

}
