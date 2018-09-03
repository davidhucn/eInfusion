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

// WriteBack :回写到前端
func (c *WSConnet) WriteBack() {
	for d := range c.sdData {
		// err := c.ws.WriteMessage(websocket.TextMessage, message)
		err := c.conn.WriteMessage(ws.TextMessage, d)
		if err != nil {
			break
		}
	}
	c.conn.Close()
}

func (c *WSConnet) reader() {
	for {
		err := c.conn.ReadJSON(&clisData)
		// _, m, err := wsconnn.ReadMessage()
		if cm.CkErr("websocket接收前端数据出错!", err) {
			// FIXME:制定通讯标准，此处应返回前端页面出错信息
			// c.conn.WriteMessage(ws.TextMessage, []byte("can't Exchange the data"))
			c.sdData <- []byte("can't Exchange the data")
			break
		}
		// 根据前端应用需求信息发送指令
		// 获取随机字符串生成orders id
		ssn := cm.GetRandString(10)
		//定时检查返回数据池里面有否相应的数据
		// go func() {
		// 	var ms sync.Mutex
		// 	cT := time.NewTicker(500 * time.Millisecond) // 定时
		// 	select {
		// 	case <-cT.C:
		// 		cm.Msg("i am alive!!")
		// 		if _, ok := eq.RcMsgs[ssn]; ok {
		// 			// 回传前端
		// 			c.conn.WriteMessage(ws.TextMessage, []byte(eq.RcMsgs[ssn]))
		// 			// 册除该条数据
		// 			ms.Lock()
		// 			defer ms.Unlock()
		// 			delete(eq.RcMsgs, ssn)
		// 		}
		// 	}
		// }()
		// 加入发送消息队列
		for i := 0; i < len(clisData); i++ {
			eq.AddToSendQueue(ssn, clisData[i].ID, cm.ConvertBasStrToUint(10, clisData[i].CmdType), clisData[i].Args)
			c.conn.WriteMessage(ws.TextMessage, []byte("Doing..."))
		}

		// wsconnn.WriteMessage(1, []byte(cliMsg[0].Action))
	}
}

// ws接收目前仅限于json
func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	defer conn.Close()
	if cm.CkErr("链接websocket出错！", err) {
		return
	}
	c := &WSConnet{sdData: make(chan []byte, 1024), conn: conn}
	// 登记注册到全局wsConnect对象

	go c.WriteBack()
	c.reader()

}
