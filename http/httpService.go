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

// loopingWebMsg :回写到WEB前端
func (w *WebClients) loopingWebMsg() {
	// 循环获取datahub包内的ws消息数据
	go func() {
		for dh.WebMsgQueue != nil {
			select {
			case od := <-dh.WebMsgQueue:
				w.Orders <- od
			}
		}
	}()
	// 遍历cmd对象，循环发送
	go func() {
		for cd := range w.Orders {
			// 截取cmd中websocket连接ID
			wsConID := DecodeToWSConnID(cd.CmdID)
			if wsConID != "" {
				if c, ok := w.Connections[wsConID]; ok {
					cm.CkErr(WebMsg.WSSendDataError, c.WriteMessage(ws.TextMessage, cd.Cmd))
				}
			}
		}
	}()
}

// reader :接受web数据
func (w *WebClients) receiveWebRequest(rWSConnID string) {
	for {
		// 如果连接错误，则返回错误信息并退出
		if cm.CkErr(WebMsg.WSReceiveDataError, w.Connections[rWSConnID].ReadJSON(&clisData)) {
			odID := NewWSOrderID(rWSConnID)
			od := cm.NewOrder(odID, []byte(WebMsg.WSReceiveDataError))
			w.Orders <- od
			w.UnregisterWSConn(rWSConnID)
			break
		}
		// 根据前端应用需求信息发送指令
		// 加入发送消息队列
		for i := 0; i < len(clisData); i++ {
			odID := NewWSOrderID(rWSConnID)
			dh.SendOrderToDeviceByTCP(odID, clisData[i].TargetID, cm.ConvertBasStrToUint(10, clisData[i].CmdType), clisData[i].Args)
			// FIXME:制定通讯标准，此处应返回前端页面完成信息代码
			od := cm.NewOrder(odID, []byte(WebMsg.WSSendDataSuccess))
			w.Orders <- od
		}
	}
}

// RegisterWSConn :注册ws连接至连接池内
func (w *WebClients) RegisterWSConn(rWSConnID string, rConn *ws.Conn) {
	w.Lock()
	w.Connections[rWSConnID] = rConn
	w.Unlock()
}

// UnregisterWSConn :注销ws连接池内websocket
func (w *WebClients) UnregisterWSConn(rWSConnID string) {
	w.Lock()
	delete(w.Connections, rWSConnID)
	w.Unlock()
}

// ws接收目前仅限于json
func wshandler(wc *WebClients, w http.ResponseWriter, r *http.Request) {
	con, err := wsupgrader.Upgrade(w, r, nil)
	defer con.Close()
	if cm.CkErr(WebMsg.WSConnectError+" IP Addr:"+con.RemoteAddr().String(), err) {
		return
	}
	{
		// 用户操作记录
		// sis, _ := CStore.Get(r, "session-name")
		// sis.Values["qid"] = wsConnID //改成数组
		// sis.Save(r, w)
	}
	// 获取随机字符串生成标识id
	wsConnID := cm.GetRandString(10)
	// 登记注册到全局wsConnect对象
	wc.RegisterWSConn(wsConnID, con)
	wc.loopingWebMsg()
	wc.receiveWebRequest(wsConnID)
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
