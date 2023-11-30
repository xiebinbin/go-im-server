package singapore

import (
	"context"
	"encoding/json"
	"fmt"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"imsdk/pkg/redis"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type CallBackOfflineParams struct {
	AChatId    string      `json:"im_achat_id"`
	AMID       string      `json:"im_amid"`
	Title      string      `json:"im_title,omitempty"`
	Body       string      `json:"im_body,omitempty"`
	SenderId   string      `json:"im_sender_id,omitempty"`
	ReceiveIds string      `json:"im_receive_ids,omitempty"`
	Type       string      `json:"im_type"`
	Extra      interface{} `json:"im_sdk_content"`
	DeviceIds  []string    `json:"im_device_ids"`
}

type TokenResp struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Errors struct {
		Token string `json:"token"`
	} `json:"errors"`
}

type CallbackResp struct {
	Data struct {
		Message string `json:"message"`
	} `json:"data"`
	Errors struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func GetToken() string {
	logCtx := log.WithFields(context.Background(), map[string]string{"action": "PostFormToThird GetToken"})
	cacheTag := "imsdk:singapore:token:" + os.Getenv("RUN_ENV")
	val, _ := redis.Client.Get(cacheTag).Result()
	if val != "" {
		return val
	}
	host := "https://dev.api.bossjob.com/api-auth/token/generate"
	secretKey := "4XDRmsekvvVqRBnxUgNDWgUFCReUNbms"
	if os.Getenv("RUN_ENV") == "release" {
		host = "https://api.bossjob.com/api-auth/token/generate"
		secretKey = "CLTeAvgsknKPnMDetwmTXGyX3p3WhDMg"
	}
	data := map[string]interface{}{
		"secret_key": secretKey,
		"client_id":  "td-im",
	}
	headers := map[string]string{}
	res, err := PostForm(host, data, headers)
	var resData TokenResp
	err = json.Unmarshal(res, &resData)
	log.Logger().Info(logCtx, "PostForm err:", resData, err)
	if resData.Data.Token != "" {
		fmt.Println("resData.Data.Token:", resData.Data.Token)
		redis.Client.Set(cacheTag, resData.Data.Token, time.Second*3600)
		return resData.Data.Token
	}
	fmt.Println("res:", resData, err)
	return ""
}

func PostForm(host string, data map[string]interface{}, headers map[string]string) ([]byte, error) {
	logCtx := log.WithFields(context.Background(), map[string]string{"action": "PostFormToThird"})
	formStrings := funcs.HttpBuildQuery(data)
	request, _ := http.NewRequest("POST", host, strings.NewReader(formStrings))
	log.Logger().Info(logCtx, "PostForm data:", host, data, headers)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset:utf-8;")
	if len(headers) > 0 {
		for k, i := range headers {
			request.Header.Set(k, i)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Logger().Error(logCtx, "PostForm err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	res, er := ioutil.ReadAll(resp.Body)
	log.Logger().Info(logCtx, "PostForm res:", string(res), er)
	return res, er
}
