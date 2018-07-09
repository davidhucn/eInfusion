package httpOperate

import (
	"eInfusion/comm"
	"net/http"

	"github.com/gin-gonic/gin"
	//	gin "gopkg.in/gin-gonic/gin.v1"
	//	"encoding/json"
)

func StartHttpServer(iPort int) {
	comm.SepLi(60, "")
	comm.Msg("start http...,Port:", iPort)
	comm.SepLi(60, "")

	gin.SetMode(gin.DebugMode) //全局设置环境，此为开发环境，线上环境为gin.ReleaseMode
	router := gin.Default()    //获得路由实例

	//注册接口
	//	router.GET("/simple/server/get", GetHandler)
	router.POST("/r/", PostHandler)
	//	router.PUT("/simple/server/put", PutHandler)
	//	router.DELETE("/simple/server/delete", DeleteHandler)

	router.Run(":" + comm.ConvertIntToStr(iPort))
}

func PostHandler(c *gin.Context) {
	//	c.GetPostForm()
	//c.GetHeader()
	//	type JsonHolder struct {
	//		Id   int    `json:"id"`
	//		Name string `json:"name"`
	//	}
	//	holder := JsonHolder{Id: 1, Name: "my name"}
	//若返回json数据，可以直接使用gin封装好的JSON方法
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	comm.SepLi(60, "")
	comm.Msg(string(buf[0:n]))
	comm.SepLi(60, "")
	resp := map[string]string{"hello": "world"}
	c.JSON(http.StatusOK, resp)
	return
}
