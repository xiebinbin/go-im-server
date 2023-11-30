package user

import (
	"fmt"
	"imsdk/pkg/funcs"
	"imsdk/pkg/redis"
)

type LoginRequest struct {
	Prefix string `json:"prefix"`
	Phone  string `json:"phone" binding:"required"`
}

type GetAuthRequest struct {
	Token string `json:"token" binding:"required"`
}

type AuthParams struct {
	Nonce     string `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
	Ver       int64  `json:"ver"`
}

type LoginResponse struct {
	Token string `json:"token"`
	AUId  string `json:"auid"`
}

type GetAuthResponse struct {
	AK       string `json:"ak"`
	AUId     string `json:"auid"`
	AuthCode string `json:"authcode"` // jsonString
}

//var ak = "68oni7jrg31qcsaijtg76qln"
//var sk = "5hNMDvExyAsbMv8PPDwCzEVwx62JSZ2NjpxmTw9pgtig"

func Login(request LoginRequest) (LoginResponse, error) {
	var data LoginResponse
	userInfo := map[string]string{
		"phone": request.Phone,
		"id":    funcs.Md516(request.Phone),
	}
	token := funcs.SHA1Base64(request.Phone)
	data = LoginResponse{
		AUId:  userInfo["id"],
		Token: token,
	}
	_, err := redis.Client.Set(token, request.Phone, 0).Result()
	if err != nil {
		fmt.Println("login save redis err----", err, err.Error())
	}
	return data, nil
}
