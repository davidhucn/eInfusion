package http

import (
	cm "eInfusion/comm"
	dh "eInfusion/datahub"
	"net/http"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

var wsupgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// loopingAddToQueueFromDataHub :循环读取datahub包内WebMsgQueue对象发送到前端web
func (w *WebClients) loopingAddToQueueFromDataHub() {
	for dh.WebMsgQueue != nil {
		select {
		case od := <-dh.WebMsgQueue:
			w.Orders <- od
		}
	}
}

// WriteBack :回写到WEB前端
func (w *WebClients) loopingWriteBack() {
	// 遍历cmd对象
	for cd := range w.Orders {
		// 截取cmd中websocket连接ID
		wsConID := DecodeToWSConnID(cd.CmdID)
		if c, ok := w.Connections[wsConID]; ok {
			cm.CkErr(WebMsg.WSSendDataError, c.WriteMessage(ws.TextMessage, cd.Cmd))
		}
	}
}

// reader :接受web数据
func (w *WebClients) reader(rWSConnID string) {
	for {
		if cm.CkErr(WebMsg.WSReceiveDataError, w.Connections[rWSConnID].ReadJSON(&clisData)) {
			odID := NewWSOrderID(rWSConnID)
			od := cm.NewOrder(odID, []byte(WebMsg.WSReceiveDataError))
			w.Orders <- od
			w.Lock()
			delete(w.Connections, rWSConnID)
			w.Unlock()
			break
		}
		// 根据前端应用需求信息发送指令
		// 加入发送消息队列
		for i := 0; i < len(clisData); i++ {
			odID := rWSConnID + "#" + cm.GetRandString(6)
			dh.SendOrderToDeviceByTCP(odID, clisData[i].ID, cm.ConvertBasStrToUint(10, clisData[i].CmdType), clisData[i].Args)
			// FIXME:制定通讯标准，此处应返回前端页面完成信息代码
			od := cm.NewOrder(odID, []byte(WebMsg.WSSendDataSuccess))
			w.Orders <- od
		}
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
		sis.Values["foo"] = "bar" //改成数组
		sis.Save(r, w)
	}
	// 获取随机字符串生成标识id
	conID := cm.GetRandString(10)
	wc.Lock()
	// 登记注册到全局wsConnect对象
	wc.Connections[conID] = con
	wc.Unlock()
	go wc.loopingWriteBack()
	go wc.loopingAddToQueueFromDataHub()
	wc.reader(conID)
}

// StartHTTPServer :开始运行httpServer
func StartHTTPServer(wc *WebClients, iPort int) {
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
	// wcs := NewWebClients()
	// websocket处理方法
	r.GET("/ws", func(c *gin.Context) {
		wshandler(wc, c.Writer, c.Request)
	})
	r.Run(":" + cm.ConvertIntToStr(iPort))
	// r.Run(":12312")
}
