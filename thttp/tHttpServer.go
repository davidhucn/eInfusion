package thttp

import (
	cm "eInfusion/comm"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	// "github.com/tidwall/gjson"
)

func StartHTTPServer(iPort int) {
	cm.SepLi(60, "")
	cm.Msg("start http...,Port:", iPort)
	cm.SepLi(60, "")

	gin.SetMode(gin.DebugMode) //全局设置环境，此为开发环境，线上环境为gin.ReleaseMode
	r := gin.Default()         //获得路由实例
	// r.LoadHTMLGlob("view/ping.tmpl")
	r.LoadHTMLFiles("view/alarm.html")
	//	router.GET("/simple/server/get", GetHandler)
	r.GET("/", func(c *gin.Context) {
		// st := "ws://localhost:12312/ws"
		// resp := map[string]string{"url": st}
		// c.JSON(http.StatusOK, resp)
		// c.HTML(http.StatusOK, "ping.tmpl", resp)
		// c.HTML(http.StatusOK, "ping.tmpl", gin.H{
		// 	"url": "ws://localhost:12312/ws",
		// })
		c.HTML(http.StatusOK, "alarm.html", nil)
	})
	// websocket处理方法
	r.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})
	// r.Run(":" + cm.ConvertIntToStr(iPort))
	r.Run(":12312")
}

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

	// ticker := time.NewTicker(5 * time.Second)
	// for range ticker.C {
	// 	data := []byte("Running")
	// 	conn.WriteMessage(1, data)
	// }
	// msg := []byte("running")
	// ticker := time.NewTicker(5 * time.Second)
	// for range ticker.C {
	// 	conn.WriteMessage(1, msg)
	// }
	// conn.WriteMessage(len(msg), msg)

	for {
		err := conn.ReadJSON(&cliMsg)
		// _, m, err := conn.ReadMessage()
		if err != nil {
			cm.Msg("error:", err)
			conn.WriteMessage(1, []byte("can't Exchange the data"))
			break
		}
		cm.Msg("meta:", cliMsg)
		// err = json.Unmarshal(m, &cliMsg)

		// cm.Msg(string(m))
		// if err != nil {
		// 	cm.Msg("error:", err)
		// }

		conn.WriteJSON(cliMsg)
		// conn.WriteMessage(1, []byte(cliMsg[0].Action))

		// fmt.Printf("%+v", reqMsg)
	}

}

// func testJ() {
// 	// var jsonBlob = []byte(` [
// 	//     { "Name" : "Platypus" , "Order" : "Monotremata" } ,
// 	//     { "Name" : "Quoll" ,     "Order" : "Dasyuromorphia" }
// 	// ] `)
// 	var jsonBlob = []byte(`[{"Name":"Platypus","Order":"Monotremata"}]`)
// 	type Animal struct {
// 		Name  string
// 		Order string
// 	}
// 	var animals []Animal
// 	err := json.Unmarshal(jsonBlob, &animals)
// 	if err != nil {
// 		fmt.Println("error:", err)
// 	}
// 	cm.SepLi(50, "")
// 	cm.Msg("local json:", cm.ConvertByteOfAscToStr(jsonBlob))
// 	cm.SepLi(50, "")
// 	cm.Msg("[]byte(local json):", jsonBlob)
// 	cm.Msg("animals:", animals)
// 	// fmt.Printf("%+v", animals)
// }
