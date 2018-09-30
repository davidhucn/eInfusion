package http

import (
	cm "eInfusion/comm"
	"net/http"
	"sync"

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
			c.WsConnection.WriteMessage(ws.TextMessage, v.Cmd)
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

func (w *WebClients) reader(rSn string) {
	for {
		err := w.Connections[rSn].ReadJSON(&clisData)
		err := c.conn.ReadJSON(&clisData)
		// _, m, err := wsconnn.ReadMessage()
		if err != nil {
			c.sdData <- []byte("ws数据连接失败!")
			var ms sync.Mutex
			ms.Lock()
			delete(WsClis, rSn)
			ms.Unlock()
			break
		}
		// 根据前端应用需求信息发送指令
		// 加入发送消息队列
		for i := 0; i < len(clisData); i++ {
			AddToSendQueue(rSn, clisData[i].ID, cm.ConvertBasStrToUint(10, clisData[i].CmdType), clisData[i].Args)
			// FIXME:制定通讯标准，此处应返回前端页面完成信息
			c.sdData <- []byte("前端发送完成，等待后端服务处理")
		}
		// wsconnn.WriteMessage(1, []byte(cliMsg[0].Action))
	}
}

// ws接收目前仅限于json
func wshandler(wc *WebClients, w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	defer conn.Close()
	if cm.CkErr(WebMsg.WSConnectError+"IP Addr:"+cm.GetPureIPAddr(conn.RemoteAddr().String()), err) {
		return
	}
	// 获取随机字符串生成标识id
	// conID := cm.GetRandString(10)
	// 登记注册到全局wsConnect对象
	wc.Lock()
	wc.Connections[cm.GetRandString(10)] = conn
	wc.Unlock()
	go wc.WriteBack()
	wc.reader(cm.GetRandString(8))
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
	// websocket处理方法
	r.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})
	r.Run(":" + cm.ConvertIntToStr(iPort))
	// r.Run(":12312")
}
