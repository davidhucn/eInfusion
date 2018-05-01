package httpOperate

import (
	"eInfusion/comm"
	"net/http"

	"github.com/gin-gonic/gin"
	//	"gopkg.in/gin-gonic/gin.v1"
)

func init() {
	//	gin.SetMode(gin.ReleaseMode)
}
func StartHttpServer() {
	comm.SepLi(60)
	comm.Msg("start http...,Port:", c_Http_Port)
	comm.SepLi(60)
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "It works")
	})

	router.Run(":" + comm.ConvertIntToStr(c_Http_Port))
}
