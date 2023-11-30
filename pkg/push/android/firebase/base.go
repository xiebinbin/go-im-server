package push

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"google.golang.org/api/option"
	"imsdk/pkg/app"
	"log"
	"os"
	"time"
)

func initializeAppWithServiceAccount() *firebase.App {
	filename := app.Config().GetPublicConfigDir() + "forward/firebase/account_key_dev.json"
	if os.Getenv("RUN_ENV") == "release" {
		filename = app.Config().GetPublicConfigDir() + "forward/firebase/account_key_pro.json"
	}

	opt := option.WithCredentialsFile(filename)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	return app
}

func sendMulticast(app *firebase.App, message *messaging.MulticastMessage) {
	// [START send_multicast]
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	br, err := client.SendMulticast(context.Background(), message)
	fmt.Printf("android messages forward err : %s \n", err)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
	// [END send_multicast]
}

func sendTestMulti(message *messaging.MulticastMessage) {
	filename := app.Config().GetPublicConfigDir() + "forward/firebase/account_key_dev.json"
	opt := option.WithCredentialsFile(filename)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	sendMulticast(app, message)
}

func sendReleaseMulti(message *messaging.MulticastMessage) {
	filename := app.Config().GetPublicConfigDir() + "forward/firebase/account_key_pro.json"
	opt := option.WithCredentialsFile(filename)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	sendMulticast(app, message)
}

func sendToToken(app *firebase.App) {
	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}
	// This registration token comes from the client FCM SDKs.
	//registrationToken := "egundd_RRsWH6J6LIuBvWK:APA91bETWhGw_RWVOvLzdQVZZAPsKy0Sbotvw78zCa25K0KTeT772f6ZkUtV3AbfQ46mZ3HgNnTdEXUdUguYoa7cynu7dczUzc1YKgcixw4CwmgWg-NP99zOFC8NpyqKCayKgOad1HP1"
	registrationToken := "fy-TeSSfSn2FlZXyA9Qcrn:APA91bGRcR-57yZ8-aTuzZ7WhZ7pIHuyOH5HfSWq0dp-zvau5KMF_colE4qLgBfr0PCvuEgGrCSVpaeAm4oZEOIr2KTHhXPNZGrT6xiX_VqlK0njTx5iowZVgW7xP2_HkBVM076kLyFN"
	message := &messaging.Message{
		Notification: NotificationMessage("Chat", "Xiao Mi is a little fairy 0000000"),
		Token:        registrationToken,
		Android:      AndroidMessage("Chat androidMessage", "Xiao Mi is a little fairy 11111111", "moments"),
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("token======")
	fmt.Println("Successfully sent message:", response)
}

func sendAll(app *firebase.App) {
	// This registration token comes from the client FCM SDKs.
	//registrationToken := "d31SWBSRTP-KNw59zDxcKo:APA91bF8ORZOCbvvJFl0n02bQikU6-wdWo9pEEdd1KVHpU_vhvW09eAXvUd9-Cgk9zoDQiqcgSsNknfDN4QVuT5kIzPixewwQbhlQaHBNm2P9DyfpE8-rwmzM91k1o0PNpq43jkr1KTa"
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	// [START send_all] Create a list containing up to 100 messages.
	messages := []*messaging.Message{
		{
			Notification: NotificationMessage("Chat androidMessage", "Xiao Mi is a little fairy 11111111"),
			//Token:        registrationToken,
			Android: AndroidMessage("Chat androidMessage", "Xiao Mi is a little fairy 11111111", "moments"),
			//APNS:         apnsMessage(),
			Topic: "industry-tech",
		},
	}
	br, err := client.SendAll(ctx, messages)
	if err != nil {
		log.Fatalln(err)
	}

	// See the BatchResponse reference documentation
	// for the contents of response.
	fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
	// [END send_all]
}

func RegistrationTokenMessage(registrationToken string) *messaging.Message {
	// [START android_message_golang]
	message := &messaging.Message{
		Token: registrationToken,
	}
	// [END android_message_golang]
	return message
}

func NotificationMessage(title, body string) *messaging.Notification {
	// [START android_message_golang]
	notification := &messaging.Notification{
		Title: title,
		Body:  body,
	}
	// [END android_message_golang]
	return notification
}

func AndroidMessage(title, body, clickAction string) *messaging.AndroidConfig {
	// [START android_message_golang]
	oneHour := time.Duration(1) * time.Hour
	Android := &messaging.AndroidConfig{
		TTL:      &oneHour,
		Priority: "normal",
		Notification: &messaging.AndroidNotification{
			Title:       title,
			Body:        body,
			Priority:    messaging.PriorityHigh,
			ClickAction: clickAction,
			//NotificationCount: count(),
			Icon: "stock_ticker_update",
		},
	} // [END android_message_golang]
	return Android
}

func sendMulticastAndHandleErrors(app *firebase.App) {
	// [START send_multicast_error]
	// Create a list containing up to 100 registration tokens.
	// This registration tokens come from the client FCM SDKs.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	registrationTokens := []string{
		"d31SWBSRTP-KNw59zDxcKo:APA91bF8ORZOCbvvJFl0n02bQikU6-wdWo9pEEdd1KVHpU_vhvW09eAXvUd9-Cgk9zoDQiqcgSsNknfDN4QVuT5kIzPixewwQbhlQaHBNm2P9DyfpE8-rwmzM91k1o0PNpq43jkr1KTa",
		"ed1blp9USqm5zjSja_oPX_:APA91bFZv43kUgydzevsvVxX0_vnPWV21bbyoMfkZd1z5hWwKqySS6YaySnSbjh_2N5cq2SemtFPbJYkTPGI1tZ00pGHL-py0pGpior-WiJv_l4ujia9i8L7uH0swuhs1-djAe7j_-4c",
	}
	message := &messaging.MulticastMessage{
		Tokens: registrationTokens,
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
	}

	br, err := client.SendMulticast(context.Background(), message)
	if err != nil {
		log.Fatalln(err)
	}

	if br.FailureCount > 0 {
		var failedTokens []string
		for idx, resp := range br.Responses {
			if !resp.Success {
				failedTokens = append(failedTokens, registrationTokens[idx])
			}
		}

		fmt.Printf("List of tokens that caused failures: %v\n", failedTokens)
	}
	// [END send_multicast_error]
}

func apnsMessage(title, body string) *messaging.APNSConfig {
	// [START apns_message_golang]
	badge := 42
	APNS := &messaging.APNSConfig{
		Headers: map[string]string{
			"apns-priority": "10",
		},
		Payload: &messaging.APNSPayload{
			Aps: &messaging.Aps{
				Alert: &messaging.ApsAlert{
					Title: title,
					Body:  body,
				},
				Badge: &badge,
			},
		},
	}
	// [END apns_message_golang]
	return APNS
}

func sendToTopic(app *firebase.App) {
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	// [START send_to_topic_golang]
	// The topic name can be optionally prefixed with "/topics/".
	topic := "highScores"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Topic: topic,
	}

	// Send a message to the devices subscribed to the provided topic.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("topic======")
	fmt.Println("Successfully sent message:", response)
	// [END send_to_topic_golang]
}
