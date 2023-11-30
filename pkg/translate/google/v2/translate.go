package v2

import (
	"cloud.google.com/go/pubsub"
	translate2 "cloud.google.com/go/translate"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
	"imsdk/pkg/app"
	"imsdk/pkg/translate"
	"os"
)

type (
	TransResponse struct {
		Error TransError `json:"error"`
	}
	TransError struct {
		Code    int           `json:"code"`
		Message string        `json:"message"`
		Errors  []TransErrors `json:"errors"`
	}
	TransErrors struct {
		Domain  string `json:"code"`
		Message string `json:"message"`
		Reason  string `json:"reason"`
	}
)

const TransUrl = "https://translation.googleapis.com/language/translate/v2"

func initializeAppWithServiceAccount() *translate2.Client {
	filename := app.Config().GetPublicConfigDir() + "push/firebase/account_key_dev.json"
	if os.Getenv("RUN_ENV") == "release" {
		filename = app.Config().GetPublicConfigDir() + "push/firebase/account_key_pro.json"
	}
	opt := option.WithCredentialsFile(filename)
	ctx := context.Background()
	client, err := translate2.NewClient(ctx, opt)
	if err != nil {
		_ = fmt.Errorf("Translate: %v", err)
	}
	defer client.Close()
	return client
}

func TranslateText(text, targetLang string) (string, error) {
	ctx := context.Background()
	client := initializeAppWithServiceAccount()
	lang, err := language.Parse(targetLang)
	if err != nil {
		return "", fmt.Errorf("language.Parse: %v", err)
	}
	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	fmt.Println("resp:", resp, err)
	if err != nil {
		return "", fmt.Errorf("Translate: %v", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("Translate returned empty response to text: %s", text)
	}
	return resp[0].Text, nil
}

func Trans(text, targetLang string) {
	//data.Key = GetConf().Key
	//data.Key = "AIzaSyA0WVq1Z9_clGzMiygaGzEULSPDhxKylJg"
	header := map[string]string{
		"SocketData-Type": "Application/json",
	}
	data := map[string]interface{}{
		"q":      text,
		"target": targetLang,
	}
	resByte, _ := CurlPost(TransUrl, data, header)
	//resStr := string(resByte)
	var res TransResponse
	err := json.Unmarshal(resByte, &res)
	fmt.Println(err)
}

func GetConf() translate.GoogleV2Conf {
	conf := translate.GetConf().GoogleV2Conf
	return conf
}
func ApiKey() error {
	// Download service account key per https://cloud.google.com/docs/authentication/production.
	// Set environment variable GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
	// This environment variable will be automatically picked up by the client.
	client, err := pubsub.NewClient(context.Background(), "chat_dev")
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()
	// Use the authenticated client.
	_ = client
	return nil
}
