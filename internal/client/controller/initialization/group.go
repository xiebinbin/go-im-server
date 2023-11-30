package initialization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	chat2 "imsdk/internal/common/dao/chat"
	"imsdk/internal/common/dao/chat/members"
	groupDetail "imsdk/internal/common/dao/group/detail"
	members2 "imsdk/internal/common/dao/group/members"
	"imsdk/pkg/funcs"
	"imsdk/pkg/response"
	"math"
)

func CountGroupMemberAmount(ctx *gin.Context) {
	groups := groupDetail.New().GetAll()
	for _, group := range groups {
		memAmount := members.New().GetChatMembersCount(group.ID)
		uData := groupDetail.Detail{
			Total: int(memAmount),
		}
		err := groupDetail.New().UpByID(group.ID, uData)
		if err != nil {
			return
		}
	}
}

func DealChatInfo(ctx *gin.Context) {
	count := groupDetail.New().GetCount()
	//counters.New().
	row := 500
	page := int(math.Ceil(float64(count) / float64(row)))
	fmt.Println("count:", count, page)
	for i := 0; i < page; i++ {
		data := groupDetail.New().GetListByLimit(int64(row), int64(i*row))
		if len(data) > 0 {
			for _, v := range data {
				ctx.Set("ak", "Chat")
				auids, _ := members2.New().GetGroupMembers(v.ID)
				chatAuids, _ := members.New().GetChatMembers(v.ID)
				if len(chatAuids) > len(auids) {
					diffIds := funcs.DifferenceString(auids, chatAuids)
					for _, id := range diffIds {
						err := members.New().Delete(id, v.ID)
						if err != nil {
							fmt.Println("Delete fail :", id, v.ID)
							continue
						}
					}
					fmt.Println("auids:", chatAuids, auids, v.ID, diffIds)
				}
			}
		}
		fmt.Println("deal page success ********** :", i)
	}
}

func DealChatOnlyTwo(ctx *gin.Context) {
	where := bson.M{"achat_id": bson.M{"$regex": "^g_"}, "only_two": 1}
	data, _ := chat2.New().GetListByCond(where)
	res := make(map[string]interface{}, 0)
	res["count"] = len(data)
	for _, datum := range data {
		err := chat2.New().UpMapByID(datum.ID, map[string]interface{}{
			"only_two": 0,
		})
		res[datum.ID] = err
		fmt.Println("err:", err)
	}
	response.RespData(ctx, res)
	return
}
