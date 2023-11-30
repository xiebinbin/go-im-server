package test

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"imsdk/internal/common/dao/conversation/usergroupsetting"
	"imsdk/internal/common/dao/message/usermessage"
	"imsdk/internal/common/model/user/token"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/proto"
	"imsdk/pkg/redis"
	"imsdk/pkg/response"
	"strings"
)

const (
	Address = "127.0.0.1:50052"
)

func GetUrlParams(ctx *gin.Context) {
	p := ctx.Param("action")
	fmt.Println("p----", p)
	pa := strings.Trim(p, "/")
	response.ResData(ctx, pa)
	return
}

func Redis(ctx *gin.Context) {
	var params struct {
		Tag string `json:"tag" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	redis.Client.Del(params.Tag)
	//filename := app.Config().GetPublicConfigDir() + "test.json"
	//filePtr, err := os.Open(filename)
	//if err != nil {
	//	fmt.Printf("Open file failed [Err:%s]\n", err.Error())
	//	return
	//}
	//defer filePtr.Close()
	//settings := make(map[string]interface{}, 0)
	//decoder := json.NewDecoder(filePtr)
	//err = decoder.Decode(&settings)
	//fmt.Println("size------", len(settings))
	//response.ResData(ctx, settings)
	//return
	isMq := true
	userDevices := []string{"1", "2", "3", "4"}
	connMap := map[string][]string{
		"host1": {"1", "2"},
		"host2": {"3", "1"},
	}
	var onlineIds []string
	if isMq {
		for _, conIds := range connMap {
			onlineIds = append(onlineIds, conIds...)
		}
	}
	offlineIds := funcs.DifferenceString(userDevices, onlineIds)
	fmt.Println("offlineIds=====", offlineIds)
	//key := "db0"
	//res := redis.Client.Set(key, "db0", time.Second*180)
	//val := redis.Client.Get(key).Val()
	//fmt.Println("err-----", res.Err())
	//fmt.Println("val-----", val)
}

func SequenceTest(ctx *gin.Context) {
	ret := usermessage.New().TestMsgWriteBack(ctx)
	response.RespData(ctx, ret)
}

func GetSequence(ctx *gin.Context) {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()
	c := __.NewSequenceClient(conn)
	req := &__.SequenceRequest{Uid: "2x2qlr88wdcz", Type: __.SequenceRequest_TYPE_MESSAGE}
	res, err := c.GetSequence(context.Background(), req)
	if err != nil {
		grpclog.Fatalln(err)
	}
	response.ResData(ctx, res.Sequence)
	return
}

func PushAgreeApplyTest(ctx *gin.Context) {
	muteSettingInfo := usergroupsetting.New().GetMuteSettingsForChatId("g_b1e5449fcf72a5ed", []string{"5bea05fac8f3e778"})
	fmt.Println(muteSettingInfo)
	return
}

func Login(ctx *gin.Context) {
	tokenInfo, _ := token.GetToken(ctx, token.GetTokenParams{
		Uid:      "a16249cfebb6",
		AK:       "Chat",
		DeviceId: ctx.Value(base.HeaderFieldDeviceId).(string),
	})
	res := GetAuthResponse{
		Id:           "a16249cfebb6",
		Token:        tokenInfo.Token,
		RefreshToken: tokenInfo.RefreshToken,
	}
	fmt.Println(res)
}
