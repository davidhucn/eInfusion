package httpOperate

import (
	"eInfusion/comm"
	"net/http"

	"github.com/gin-gonic/gin"
	//	gin "gopkg.in/gin-gonic/gin.v1"
	//	"encoding/json"
)

func StartHttpServer() {
	comm.SepLi(60)
	comm.Msg("start http...,Port:", c_Http_Port)
	comm.SepLi(60)

	gin.SetMode(gin.DebugMode) //全局设置环境，此为开发环境，线上环境为gin.ReleaseMode
	router := gin.Default()    //获得路由实例

	//注册接口
	//	router.GET("/simple/server/get", GetHandler)
	router.POST("/simple/server/post", PostHandler)
	//	router.PUT("/simple/server/put", PutHandler)
	//	router.DELETE("/simple/server/delete", DeleteHandler)

	router.Run(":" + comm.ConvertIntToStr(c_Http_Port))
}

func PostHandler(c *gin.Context) {
	//c.GetHeader()
	type JsonHolder struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	holder := JsonHolder{Id: 1, Name: "my name"}
	//若返回json数据，可以直接使用gin封装好的JSON方法
	c.JSON(http.StatusOK, holder)
	return
}
