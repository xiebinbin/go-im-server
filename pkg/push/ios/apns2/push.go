package apns2

import (
	"fmt"
	"imsdk/pkg/app"
	"imsdk/pkg/push/ios/apns2/certificate"
	"imsdk/pkg/push/ios/apns2/payload"
	"log"
	"time"
)

type PushRequest struct {
	DeviceToken   string    `json:"device_token"`
	Badge         int       `json:"badge,omitempty"`
	AlertTitle    string    `json:"alert_title,omitempty"`
	AlertSubTitle string    `json:"alert_subtitle,omitempty"`
	AlertBody     string    `json:"alert_body,omitempty"`
	Sound         string    `json:"sound,omitempty"`
	AlertAction   string    `json:"alert_action,omitempty"`
	Category      string    `json:"category,omitempty"`
	Topic         string    `json:"topic,omitempty"`
	CollapseID    string    `json:"collapse_id,omitempty"`
	PushType      EPushType `json:"push_type,omitempty"`
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
	sound := "default"
	if request.Sound != "" {
		sound = request.Sound
	}
	//“passive”, “active”, " time-sensitive”, or “critical”
	payloadContent := payload.NewPayload().Badge(request.Badge).Category(request.Category).
		AlertTitle(request.AlertTitle).AlertSubtitle("").AlertBody(request.AlertBody).Sound(sound).
		AlertAction(request.AlertAction).SoundName(sound).InterruptionLevel("time-sensitive")
	pushType := PushTypeAlert
	topic := "com.chat.im"
	switch request.PushType {
	case PushTypeVOIP:
		pushType = request.PushType
		topic = "com.chat.im.voip"
		break
	}
	var notification = &Notification{
		DeviceToken: request.DeviceToken,
		Topic:       topic,
		Priority:    10,
		PushType:    pushType,
		CollapseID:  request.CollapseID,
		Expiration:  time.Time{},
		Payload:     payloadContent,
		//ApnsID:      "1",
	}
	fmt.Println("apns notification : ", notification, "request:", request, "sound:", sound)
	res, err := client.Push(notification)
	fmt.Println("apns notification-res : ", res, " ,err: ", err)
	if err != nil {
		//log.Fatal("Error: ", err)
		fmt.Printf("APNS Error: %v", err)
	} else {
		fmt.Printf("%v: '%v'\n", res.StatusCode, res.Reason)
	}
}
