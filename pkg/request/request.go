package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func Get(url string, headers map[string]string) ([]byte, error) {
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

func Post(host string, data interface{}) ([]byte, error) {
	byteData, _ := json.Marshal(data)
	request, err := http.NewRequest("POST", host, strings.NewReader(string(byteData)))
	if err != nil {
		return []byte{}, err
	}
	request.Header.Set("Content-Type", "Application/json")
	client := &http.Client{}
	resp, err := client.Do(request) //
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
