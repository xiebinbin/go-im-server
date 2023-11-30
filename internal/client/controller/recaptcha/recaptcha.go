package recaptcha

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"imsdk/pkg/errno"
	"imsdk/pkg/recaptcha"
	"imsdk/pkg/response"
	"log"
	"net/http"
)

var recaptchaPublicKey string

const (
	pageTop = `<!DOCTYPE HTML><html><head>
<style>.error{color:#ff0000;} .ack{color:#0000ff;}</style><title>Recaptcha Test</title></head>
<body><div style="width:100%"><div style="width: 50%;margin: 0 auto;">
<h3>Recaptcha Test</h3>
<p>This is a test form for the go-recaptcha package</p>`
	form = `<form action="/" method="POST">
	    <script src="https://www.google.com/recaptcha/api.js"></script>
			<div class="g-recaptcha" data-sitekey="%s"></div>
    	<input type="submit" name="button" value="Ok">
</form>`
	pageBottom = `</div></div></body></html>`
	anError    = `<p class="error">%s</p>`
	anAck      = `<p class="ack">%s</p>`
)

// processRequest accepts the http.Request object, finds the reCaptcha form variables which
// were input and sent by HTTP POST to the server, then calls the recaptcha package's Confirm()
// method, which returns a boolean indicating whether or not the client answered the form correctly.
func processRequest(request *http.Request) bool {
	recaptchaResponse, responseFound := request.Form["g-recaptcha-response"]
	if responseFound {
		result, err := recaptcha.Confirm("127.0.0.1", recaptchaResponse[0])
		if err != nil {
			log.Println("recaptcha server error", err)
		}
		return result
	}
	return false
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm() // Must be called before writing response
	fmt.Fprint(writer, pageTop)
	if err != nil {
		fmt.Fprintf(writer, fmt.Sprintf(anError, err))
	} else {
		_, buttonClicked := request.Form["button"]
		if buttonClicked {
			if processRequest(request) {
				fmt.Fprint(writer, fmt.Sprintf(anAck, "Recaptcha was correct!"))
			} else {
				fmt.Fprintf(writer, fmt.Sprintf(anError, "Recaptcha was incorrect; try again."))
			}
		}
	}
	fmt.Fprint(writer, fmt.Sprintf(form, recaptchaPublicKey))
	fmt.Fprint(writer, pageBottom)
}

func ReCaptcha(ctx *gin.Context) {
	var params struct {
		Token string `json:"token" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	//result, _ := recaptcha.Confirm("110.185.192.123", token)
	recaptcha.CaptchaOne(params.Token)
	//fmt.Println("result-------", result)
	//recaptcha.CaptchaOne()
}
