package pkg

import (
	"fmt"
	json "github.com/json-iterator/go"
	"imsdk/pkg/funcs"
	"imsdk/pkg/request"
)

type jsonString = string

//const imDebugHost = "https://demo-sdk-ser-api.buzzmsg.com/"
//const imReleaseHost = "https://ssi.buzzmsg.com/"

//const imDebugHost = "https://dev-sdk-ser-api.chat.com.tr/"
const imDebugHost = "http://localhost:7500/"
const imReleaseHost = "https://imsi.chat.com.tr/"

type CurlResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ErrCode int         `json:"err_code"`
		ErrMsg  string      `json:"err_msg"`
		Items   interface{} `json:"items"`
	} `json:"data"`
}

func Curl(uri string, data interface{}) (CurlResponse, error) {
	imHost := imDebugHost
	if funcs.GetEnv() == "release" {
		imHost = imReleaseHost
	}
	re, _ := request.Post(imHost+uri, data)
	fmt.Println("imHost+uri-----", imHost+uri)
	var res CurlResponse
	err := json.Unmarshal(re, &res)
	fmt.Printf("Curl res ------- %+v", res)
	return res, err
}
