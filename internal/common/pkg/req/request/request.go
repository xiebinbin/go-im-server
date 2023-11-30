package request

import (
	"context"
	"errors"
	json "github.com/json-iterator/go"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"io/ioutil"
	"net/http"
	"strings"
)

func PostJson(ctx context.Context, host string, data interface{}) ([]byte, error) {
	byteData, _ := json.Marshal(data)
	logCtx := log.WithFields(ctx, map[string]string{"action": "PostJson"})
	log.Logger().Info(logCtx, "params: ", string(byteData))
	request, _ := http.NewRequest("POST", host, strings.NewReader(string(byteData)))
	request.Header.Set("Content-Type", "Application/json")
	client := &http.Client{}
	resp, err := client.Do(request) //
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("fail; http status: " + resp.Status)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func FormRequest(host string, data map[string]interface{}) (*http.Response, error) {
	formStrings := funcs.HttpBuildQuery(data)
	request, _ := http.NewRequest("POST", host, strings.NewReader(formStrings))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	return client.Do(request)
}

func InnerReq(ctx context.Context, host string, data interface{}) ([]byte, error) {
	byteData, _ := json.Marshal(data)
	if !strings.Contains(host, "http://") && !strings.Contains(host, "https://") {
		host = "http://" + host
	}
	req, er := http.NewRequest("POST", host, strings.NewReader(string(byteData)))
	logCtx := log.WithFields(ctx, map[string]string{"action": "InnerReq"})
	log.Logger().Info(logCtx, "InnerReq: err: ", er, " host: ", host, " params: ", string(byteData))
	req.Header.Set("Content-Type", "Application/json")
	client := &http.Client{}
	resp, err := client.Do(req) //
	if err != nil {
		log.Logger().Error(logCtx, "************* request.InnerReq ")
		log.Logger().Error(logCtx, "InnerReq client.Do err:", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("fail; http status: " + resp.Status)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
