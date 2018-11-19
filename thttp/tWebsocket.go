package thttp

// import (
// 	cm "eInfusion/comm"
// 	logs "eInfusion/tlogs"
// 	"net/http"
// 	"sync"

// 	ws "github.com/gorilla/websocket"
// )

// var wsupgrader = ws.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// // WriteBack :回写到前端
// func (c *WSConnet) WriteBack() {
// 	for d := range c.sdData {
// 		// err := c.ws.WriteMessage(websocket.TextMessage, message)
// 		err := c.conn.WriteMessage(ws.TextMessage, d)
// 		if err != nil {
// 			break
// 		}
// 		// cm.SepLi(30, "")
// 		// for one := range WsClis {
// 		// 	cm.Msg(one)
// 		// }
// 		// cm.SepLi(30, "")
// 	}
// 	// c.conn.Close()

// }

// func (c *WSConnet) reader(rSn string) {
// 	for {
// 		err := c.conn.ReadJSON(&clisData)
// 		// _, m, err := wsconnn.ReadMessage()
// 		if err != nil {
// 			c.sdData <- []byte("ws数据连接失败!")
// 			var ms sync.Mutex
// 			ms.Lock()
// 			delete(WsClis, rSn)
// 			ms.Unlock()
// 			break
// 		}
// 		// 根据前端应用需求信息发送指令
// 		// 加入发送消息队列
// 		for i := 0; i < len(clisData); i++ {
// 			// AddToSendQueue(rSn, clisData[i].ID, cm.ConvertBasStrToUint(10, clisData[i].CmdType), clisData[i].Args)
// 			// FIXME:制定通讯标准，此处应返回前端页面完成信息
// 			c.sdData <- []byte("前端发送完成，等待后端服务处理")
// 		}
// 		// wsconnn.WriteMessage(1, []byte(cliMsg[0].Action))
// 	}
// }

// // ws接收目前仅限于json
// func wshandler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := wsupgrader.Upgrade(w, r, nil)
// 	defer conn.Close()
// 	if cm.CkErr("websocket连接出错！", err) {
// 		logs.LogMain.Error("http-websocket 连接出错！IP：【", conn.RemoteAddr().String(), "】")
// 	}
// 	// 获取随机字符串生成标识id
// 	ssn := cm.GetRandString(10)
// 	c := &WSConnet{sdData: make(chan []byte, 1024), conn: conn}
// 	// 登记注册到全局wsConnect对象
// 	var ms sync.Mutex
// 	ms.Lock()
// 	WsClis[ssn] = c
// 	ms.Unlock()
// 	go c.WriteBack()
// 	c.reader(ssn)
// }
