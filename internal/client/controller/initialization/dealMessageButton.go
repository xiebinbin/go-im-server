package initialization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/internal/common/dao/message/detail"
	"imsdk/internal/common/model/message"
	"imsdk/pkg/response"
	"math"
)

func DealMsgType19(ctx *gin.Context) {
	where := bson.M{"type": 19}
	//where := bson.M{}
	dao := detail.New()
	count := dao.GetCount(where)
	row := 1000
	page := int(math.Ceil(float64(count) / float64(row)))
	fmt.Println("count:", count, page)
	for i := 0; i < page; i++ {
		data := dao.GetListByPage(int64(row), int64(i*row), where)
		if len(data) > 0 {
			for _, v := range data {
				err := message.DealMessageButtonTemp(message.SendMessageParamsTemp{
					Content:    v.Content,
					Type:       v.Type,
					Mid:        v.ID,
					CreateTime: v.CreatedAt,
				})
				if err != nil {
					fmt.Println("err----", v.ID, err)
				}
			}
		}
		fmt.Println("deal page success ********** :", i)
	}
	response.ResData(ctx, map[string]interface{}{
		"count": count,
		"page":  page,
	})
	return
}
