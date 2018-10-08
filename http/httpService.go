package http

import (
	cm "eInfusion/comm"
	"net/http"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

var wsupgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WriteBack :回写到前端
func (w *WebClients) WriteBack() {
	var c *WsObject
	for _, c = range w.Connections {
		for v := range c.Orders {
			c.WSConnection.WriteMessage(ws.TextMessage, v.Cmd)
		}
	}
	// if cm.CkErr(WebMsg.WSSendDataError, c[].WriteMessage(ws.TextMessage, d.Cmd)) {
	// 	break
	// }
	// err := c.conn.WriteMessage(ws.TextMessage, d)
	// if err != nil {
	// 	break
	// }
	// cm.SepLi(30, "")
	// for one := range WsClis {
	// 	cm.Msg(one)
	// }
	// cm.SepLi(30, "")

	// c.conn.Close()
}

// reader :接受web数据
func (w *WebClients) reader(rConID string) {
	for {
		if cm.CkErr(WebMsg.WSReceiveDataError, w.Connections[rConID].WSConnection.ReadJSON(&clisData)) {
			od := NewOrder([]byte(WebMsg.WSReceiveDataError), cm.GetRandString(8))
			w.Connections[rConID].Orders <- od
			delete(w.Connections, rConID)
			break
		}
		// 根据前端应用需求信息发送指令
		// 加入发送消息队列
		// for i := 0; i < len(clisData); i++ {
		// 	AddToSendQueue(rSn, clisData[i].ID, cm.ConvertBasStrToUint(10, clisData[i].CmdType), clisData[i].Args)
		// 	// FIXME:制定通讯标准，此处应返回前端页面完成信息
		// 	c.sdData <- []byte("前端发送完成，等待后端服务处理")
		// }
		// wsconnn.WriteMessage(1, []byte(cliMsg[0].Action))
	}
}

// ws接收目前仅限于json
func wshandler(wc *WebClients, w http.ResponseWriter, r *http.Request) {
	con, err := wsupgrader.Upgrade(w, r, nil)
	defer con.Close()
	if cm.CkErr(WebMsg.WSConnectError+" IP Addr:"+cm.GetPureIPAddr(con.RemoteAddr().String()), err) {
		return
	}
	// TODO:记录用户操作的内容
	{
		sis, _ := CStore.Get(r, "session-name")
		sis.Values["foo"] = "bar"
		sis.Save(r, w)
	}

	// 获取随机字符串生成标识id
	// conID := cm.GetRandString(10)
	// 登记注册到全局wsConnect对象
	conID := cm.GetRandString(10)
	wc.Lock()
	wc.Connections[conID].WSConnection = con
	wc.Unlock()
	go wc.WriteBack()
	wc.reader(conID)
}

// StartHTTPServer :开始运行httpServer
func StartHTTPServer(iPort int) {
	cm.SepLi(60, "")
	cm.Msg("start http...,Port:", iPort)
	cm.SepLi(60, "")
	gin.SetMode(gin.DebugMode) //全局设置环境，此为开发环境，线上环境为gin.ReleaseMode
	r := gin.Default()         //获得路由实例
	r.LoadHTMLGlob("view/alarm.tmpl")
	// r.LoadHTMLFiles("view/alarm.html")
	r.GET("/", func(c *gin.Context) {
		st := "ws://localhost:" + cm.ConvertIntToStr(iPort) + "/ws"
		para := map[string]string{"url": st}
		// c.JSON(http.StatusOK, resp)
		c.HTML(http.StatusOK, "alarm.tmpl", para)
		// c.HTML(http.StatusOK, "ping.tmpl", gin.H{
		// 	"url": "ws://localhost:12312/ws",
		// })
		// c.HTML(http.StatusOK, "alarm.html", nil)
	})
	wcs := NewWebClients()
	// websocket处理方法
	r.GET("/ws", func(c *gin.Context) {
		wshandler(wcs, c.Writer, c.Request)
	})
	r.Run(":" + cm.ConvertIntToStr(iPort))
	// r.Run(":12312")
}
