package httpOperate

import (
	cm "eInfusion/comm"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func StartHTTPServer(iPort int) {
	cm.SepLi(60, "")
	cm.Msg("start http...,Port:", iPort)
	cm.SepLi(60, "")

	gin.SetMode(gin.DebugMode) //全局设置环境，此为开发环境，线上环境为gin.ReleaseMode
	r := gin.Default()         //获得路由实例
	r.LoadHTMLGlob("view/ping.tmpl")
	// r.LoadHTMLFiles("view/index.html")
	//	router.GET("/simple/server/get", GetHandler)
	r.GET("/", func(c *gin.Context) {
		// st := cm.ConvertIntToStr(iPort)
		// resp := map[string]string{"port": st}
		// c.JSON(http.StatusOK, resp)
		// c.HTML(http.StatusOK, "index.html", resp)
		c.HTML(http.StatusOK, "ping.tmpl", gin.H{
			"url": "ws://localhost:12312/ws",
		})
		// c.HTML(http.StatusOK, "index.html", nil)
	})
	// websocket处理方法
	r.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})
	// r.Run(":" + cm.ConvertIntToStr(iPort))
	r.Run("localhost:12312")
}

// func rootHandle(c *gin.Context) {
// 	//	c.GetPostForm()
// 	//c.GetHeader()
// 	//	type JsonHolder struct {
// 	//		Id   int    `json:"id"`
// 	//		Name string `json:"name"`
// 	//	}
// 	//	holder := JsonHolder{Id: 1, Name: "my name"}
// 	//若返回json数据，可以直接使用gin封装好的JSON方法
// 	// buf := make([]byte, 1024)
// 	// n, _ := c.Request.Body.Read(buf)
// 	// comm.SepLi(60, "")
// 	// comm.Msg(string(buf[0:n]))
// 	// comm.SepLi(60, "")
// 	// resp := map[string]string{"hello": "world"}
// 	// c.JSON(http.StatusOK, resp)
// 	// return
// 	c.JSON(http.StatusOK,gin.H{
// 		"port":
// 	}
// }
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		cm.Msg("Failed to set websocket upgrade: %v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}
