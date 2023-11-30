package aliyun

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
)

type StsClient struct {
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Region          string `json:"region"`
	Role            *Role  `json:"role"`
}

type Role struct {
	Arn         string `json:"arn"`
	SessionName string `json:"session_name"`
}

type StsResponse struct {
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	SessionToken    string `json:"session_token"`
	Expire          string `json:"expire"`
}

func NewSTSClient(region, accessKeyId, accessKeySecret string) *StsClient {
	return &StsClient{
		Region:          region,
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
}

func (s *StsClient) SetRole(arn, sessionName string) *StsClient {
	s.Role = &Role{
		Arn:         arn,
		SessionName: sessionName,
	}
	return s
}

func (s *StsClient) GetSts() (StsResponse, error) {
	client, err := sts.NewClientWithAccessKey(s.Region, s.AccessKeyId, s.AccessKeySecret)
	if err != nil {
		fmt.Println(err.Error(), "------step1")
		return StsResponse{}, nil
	}
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"

	request.RoleArn = s.Role.Arn
	request.RoleSessionName = s.Role.SessionName
	response, er := client.AssumeRole(request)
	if er != nil {
		fmt.Println(er.Error(), "------step2")
		return StsResponse{}, nil
	}
	fmt.Printf("response is %#v\n", response.BaseResponse.GetHttpStatus())
	return StsResponse{
		AccessKeyId:     response.Credentials.AccessKeyId,
		AccessKeySecret: response.Credentials.AccessKeySecret,
		SessionToken:    response.Credentials.SecurityToken,
		Expire:          response.Credentials.Expiration,
	}, nil
}
