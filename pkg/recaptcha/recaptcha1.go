package recaptcha

import (
	recapitulates "cloud.google.com/go/recaptchaenterprise/v2/apiv1"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/cloud/recaptchaenterprise/v1"
	"imsdk/pkg/app"
	"imsdk/pkg/errno"
)

const (
	ProjectID = "chat-dev"
	SiteKey   = "6Ld36ewgAAAAACvK8qeJNM0IeRdibUXOAamir_RS"
	Action    = ""
)

func CaptchaOne(token string) error {
	return createAssessment(ProjectID, SiteKey, token, Action)
}

/**
* Create an assessment to analyze the risk of an UI action.
*
* @param projectID: GCloud Project ID
* @param recaptchaSiteKey: Site key obtained by registering a domain/app to use recaptcha services.
* @param token: The token obtained from the client on passing the recaptchaSiteKey.
* @param recaptchaAction: Action name corresponding to the token.
 */
func createAssessment(projectID, recaptchaSiteKey, token, recaptchaAction string) error {

	// Create the recaptcha client.
	// TODO: To avoid memory issues, move this client generation outside
	// of this example, and cache it (recommended) or call client.close()
	// before exiting this method.
	filename := app.Config().GetPublicConfigDir() + "forward/firebase/recaptcha.json"
	opt := option.WithCredentialsFile(filename)
	ctx := context.Background()
	client, err := recapitulates.NewClient(ctx, opt)
	if err != nil {
		fmt.Printf("Error creating reCAPTCHA client\n")
		return errno.Add("Error creating reCAPTCHA client", errno.DefErr)
	}
	defer client.Close()

	// Set the properties of the event to be tracked.
	event := &recaptchaenterprise.Event{
		Token:   token,
		SiteKey: recaptchaSiteKey,
	}
	fmt.Println("222222222222")
	assessment := &recaptchaenterprise.Assessment{
		Event: event,
	}

	// Build the assessment request.
	request := &recaptchaenterprise.CreateAssessmentRequest{
		Assessment: assessment,
		Parent:     fmt.Sprintf("projects/%s", projectID),
	}
	fmt.Println("33333333333")
	response, err := client.CreateAssessment(
		ctx,
		request)

	if err != nil {
		fmt.Printf("%v", err.Error())
		return errno.Add("recaptcha response err : "+err.Error(), errno.DefErr)
	}
	fmt.Println("response==========", response.TokenProperties.Valid)
	fmt.Println("response TokenProperties==========", response.TokenProperties.Action)

	// Check if the token is valid.
	if response.TokenProperties.Valid == false {
		fmt.Printf("The CreateAssessment() call failed because the token"+
			" was invalid for the following reasons: %v",
			response.TokenProperties.InvalidReason)
		return errno.Add("response.TokenProperties.InvalidReason", errno.DefErr)
	}

	// Check if the expected action was executed.
	if response.TokenProperties.Action == recaptchaAction {
		// Get the risk score and the reason(s).
		// For more information on interpreting the assessment,
		// see: https://cloud.google.com/recaptcha-enterprise/docs/interpret-assessment
		fmt.Printf("The reCAPTCHA score for this token is:  %v",
			response.RiskAnalysis.Score)
		for _, reason := range response.RiskAnalysis.Reasons {
			fmt.Printf(reason.String() + "\n")
		}
	}
	//fmt.Printf("The action attribute in your reCAPTCHA tag does " +
	//	"not match the action you are expecting to score")
	return nil
}
