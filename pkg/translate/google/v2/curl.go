package v2

import (
	"fmt"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"strings"
)

func CurlPost(url string, data interface{}, headers map[string]string) ([]byte, error) {
	byteData, _ := json.Marshal(data)
	fmt.Println("string(byteData)=======", string(byteData))
	request, _ := http.NewRequest("POST", url, strings.NewReader(string(byteData)))
	//if headers != nil {
	//	for k, v := range headers {
	//		request.Header.Set(k, v)
	//	}
	//}
	//request.Header.Set("SocketData-Type", "Application/json")
	client := &http.Client{}
	fmt.Println("request======", request)
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func CurlGet(url string, headers map[string]string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req) //
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
