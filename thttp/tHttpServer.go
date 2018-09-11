package thttp

import (
	cm "eInfusion/comm"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StartHTTPServer :开始运行httpServer
func StartHTTPServer(iPort int) {
	// go StartSendQueueListener()

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

// ticker := time.NewTicker(5 * time.Second)
// for range ticker.C {
// 	data := []byte("Running")
// 	conn.WriteMessage(1, data)
