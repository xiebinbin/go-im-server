package health

import "github.com/gin-gonic/gin"

func AwsHealthCheck(ctx *gin.Context) {
	ctx.String(200, "success")
	return
}
