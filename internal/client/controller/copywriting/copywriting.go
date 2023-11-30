package copywriting

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/copywriting"
	"imsdk/pkg/response"
)

func List(ctx *gin.Context) {
	data := copywriting.List()
	response.ResData(ctx, data)
	return
}
