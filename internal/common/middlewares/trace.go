package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Cors(c *gin.Context) {
	method := c.Request.Method
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "SocketData-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, SocketData-Type, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")
	//if method == "OPTIONS" {
	//	c.AbortWithStatus(http.StatusNoContent)
	//}
	fmt.Println("Cors method:", method)
	c.Next()
}
