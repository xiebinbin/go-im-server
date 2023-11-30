package push

import (
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
)

type NotificationRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
type AndroidConfigRequest struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	ClickAction string `json:"click_action"`
	Extra       string `json:"extra"`
}

func initService() *firebase.App {
	return initializeAppWithServiceAccount()
}

func SendMulticast(tokens []string, notification NotificationRequest, androidConfig AndroidConfigRequest, data map[string]string) {
	data = map[string]string{
		"title":        androidConfig.Title,
		"body":         androidConfig.Body,
		"click_action": androidConfig.ClickAction,
		"extra":        androidConfig.Extra,
	}
	message := &messaging.MulticastMessage{
		Tokens:       tokens,
		Notification: NotificationMessage(notification.Title, notification.Body),
		Android:      AndroidMessage(androidConfig.Title, androidConfig.Body, androidConfig.ClickAction),
		Data:         data,
	}
	sendMulticast(initService(), message)
}

func SendToKenTest() {
	sendToToken(initService())
}

func SendReleaseMulti(tokens []string) {
	message := &messaging.MulticastMessage{
		Tokens:       tokens,
		Notification: NotificationMessage("notification.Title Release", "notification.Body Release"),
		Android:      AndroidMessage("androidConfig.Title Release", "androidConfig.Body Release", "message"),
	}
	sendReleaseMulti(message)
}

func SendTestMulti(tokens []string) {
	message := &messaging.MulticastMessage{
		Tokens:       tokens,
		Notification: NotificationMessage("notification.Title  test", "notification.Body test"),
		Android:      AndroidMessage("androidConfig.Title test", "androidConfig.Body test", "message"),
	}
	sendTestMulti(message)
}

func SendMulticastNotification(tokens []string, notification NotificationRequest, data map[string]string) {
	message := &messaging.MulticastMessage{
		Tokens:       tokens,
		Notification: NotificationMessage(notification.Title, notification.Body),
		Data:         data,
	}
	sendMulticast(initService(), message)
}

func SendMulticastAndroidMsg(tokens []string, androidConfig AndroidConfigRequest, data map[string]string) {
	message := &messaging.MulticastMessage{
		Tokens:  tokens,
		Android: AndroidMessage(androidConfig.Title, androidConfig.Body, androidConfig.ClickAction),
		Data:    data,
	}
	sendMulticast(initService(), message)
}

func SendAll() {
	sendAll(initService())
}
