package pushkit

import (
	"fmt"
	"imsdk/pkg/app"
	"imsdk/pkg/push/ios/apns2/certificate"
	"imsdk/pkg/push/ios/apns2/payload"
	"log"
	"time"
)

type PushRequest struct {
	DeviceToken   string `json:"device_token"`
	Badge         int    `json:"badge,omitempty"`
	AlertTitle    string `json:"alert_title,omitempty"`
	AlertSubTitle string `json:"alert_subtitle,omitempty"`
	AlertBody     string `json:"alert_body,omitempty"`
	Sound         string `json:"sound,omitempty"`
	AlertAction   string `json:"alert_action,omitempty"`
	Category      string `json:"category,omitempty"`
}

func Push(certificatePath, env string, request PushRequest) {
	cert, pemErr := certificate.FromPemFile(certificatePath, "123456")
	if pemErr != nil {
		log.Fatalf("Error retrieving certificate `%v`: %v", certificatePath, pemErr)
	}
	//str, _ := json.Marshal(cert)
	client := NewClient(cert)
	if env != app.ReleaseModel {
		client.Development()
	} else {
		client.Production()
	}
	//client.Development()
	payloadContent := payload.NewPayload().Badge(request.Badge).Category(request.Category).
		AlertTitle(request.AlertTitle).AlertSubtitle("").AlertBody(request.AlertBody).Sound("default").
		AlertAction(request.AlertAction)
	var notification = &Notification{
		DeviceToken: request.DeviceToken,
		Topic:       "com.chat.im.voip",
		Priority:    10,
		PushType:    PushTypeVOIP,
		CollapseID:  "1",
		Expiration:  time.Time{},
		Payload:     payloadContent,
		//ApnsID:      "1",
	}
	res, err := client.Push(notification)
	if err != nil {
		//log.Fatal("Error: ", err)
		fmt.Printf("APNS Error: %v", err)
	} else {
		fmt.Printf("%v: '%v'\n", res.StatusCode, res.Reason)
	}
}
