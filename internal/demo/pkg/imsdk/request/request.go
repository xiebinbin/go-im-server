package request

import (
	"context"
	"encoding/json"
	"fmt"
	"imsdk/internal/demo/pkg/imsdk"

	//"imsdk/internal/common/pkg/config"
	"imsdk/internal/demo/pkg/imsdk/utils"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	HTTP  = "HTTP"
	HTTPS = "HTTPS"

	DefaultHttpPort = "5000"

	GET     = "GET"
	PUT     = "PUT"
	POST    = "POST"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"

	Json = "application/json"
	Xml  = "application/xml"
	Raw  = "application/octet-stream"
	Form = "application/x-www-form-urlencoded"

	Header = "Header"
	Query  = "Query"
	Body   = "Body"
	Path   = "Path"

	HeaderSeparator = "\n"
	//ImsdkServerHost = "http://172.17.0.1:7500/"
	//ImsdkServerDebugHost   = " https://dev-sdk-ser-api.chat.com.tr/"
	ImsdkServerReleaseHost = "https://imsi.chat.com.tr/"
	ImsdkServerDebugHost   = "http://localhost:7500/"
)

type Request struct {
	Scheme         string
	Host           string
	Method         string
	Domain         string
	Port           string
	ReadTimeout    time.Duration
	ConnectTimeout time.Duration
	userAgent      map[string]string
	version        string
	actionName     string
	AcceptFormat   string
	QueryParams    map[string]string
	Headers        map[string]string
	FormParams     map[string]string
	Content        []byte
}

func New() *Request {
	return &Request{}
}

func GetImsdkServerHost(ctx context.Context, model imsdk.ModelType) string {
	//imsdkServerHost, _ := config.GetIMSDKServer()
	//return imsdkServerHost
	if model == imsdk.ModelRelease {
		return ImsdkServerReleaseHost
	}
	return ImsdkServerDebugHost
}

func (request *Request) Post() ([]byte, error) {
	params := request.GetContent()
	//byteData, _ := json.Marshal(params)
	//fmt.Println(" strings.NewReader(string(byteData)):", strings.NewReader(string(byteData)))
	req, err := http.NewRequest(request.Method, request.Host, strings.NewReader(string(params)))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", Json)
	req.Header.Set("req-id", utils.GenerateNonce())
	if request.Headers != nil {
		for k, v := range request.Headers {
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

func (request *Request) Get() ([]byte, error) {
	params := request.GetQueryParams()
	byteData, _ := json.Marshal(params)
	fmt.Println(" strings.NewReader(string(byteData)):", strings.NewReader(string(byteData)))
	req, err := http.NewRequest(request.Method, request.Host, strings.NewReader(string(byteData)))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", Json)
	if request.Headers != nil {
		for k, v := range request.Headers {
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

func (request *Request) SetContentType(contentType string) *Request {
	if contentType == "" {
		contentType = Json
	}
	request.addHeaderParam("Content-Type", contentType)
	return request
}

func (request *Request) GetContentType() (contentType string, contains bool) {
	contentType, contains = request.Headers["Content-Type"]
	return
}

func (request *Request) SetHost(host string) *Request {
	request.Host = host
	return request
}

func (request *Request) GetHost() string {
	return request.Host
}

func (request *Request) SetMethod(method string) *Request {
	request.Method = method
	return request
}

func (request *Request) GetMethod() string {
	return request.Method
}

func (request *Request) SetContent(content []byte) *Request {
	request.Content = content
	return request
}

func (request *Request) GetContent() []byte {
	return request.Content
}

func (request *Request) SetVersion(version string) *Request {
	request.version = version
	return request
}

func (request *Request) GetVersion() string {
	return request.version
}

func (request *Request) GetActionName() string {
	return request.actionName
}

func (request *Request) GetQueryParams() map[string]string {
	return request.QueryParams
}

func (request *Request) GetFormParams() map[string]string {
	return request.FormParams
}

func (request *Request) GetUserAgent() map[string]string {
	return request.userAgent
}

func (request *Request) AppendUserAgent(key, value string) {
	newkey := true
	if request.userAgent == nil {
		request.userAgent = make(map[string]string)
	}
	if strings.ToLower(key) != "core" && strings.ToLower(key) != "go" {
		for tag, _ := range request.userAgent {
			if tag == key {
				request.userAgent[tag] = value
				newkey = false
			}
		}
		if newkey {
			request.userAgent[key] = value
		}
	}
}

func (request *Request) addHeaderParam(key, value string) {
	request.Headers[key] = value
}

func (request *Request) addQueryParam(key, value string) {
	request.QueryParams[key] = value
}

func (request *Request) addFormParam(key, value string) {
	request.FormParams[key] = value
}

func (request *Request) GetDomain() string {
	return request.Domain
}

func (request *Request) SetDomain(host string) *Request {
	request.Domain = host
	return request
}
