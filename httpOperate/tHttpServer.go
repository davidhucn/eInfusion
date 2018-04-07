package httpOperate

import (
	"eInfusion/comm"
	"net/http"

	"github.com/gin-gonic/gin"
	//	"gopkg.in/gin-gonic/gin.v1"
)

func StartHttpServer() {
	comm.Msg("start http..")
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "It works")
	})

	router.Run(":7090")
}
