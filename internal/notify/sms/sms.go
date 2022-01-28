package sms

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/pkg/errors"
)

//接收短信的手机号码。
//
//手机号码格式：
//国内短信：+/+86/0086/86或无任何前缀的11位手机号码，例如1590000****。
//国际/港澳台消息：国际区号+号码，例如852000012****。
//支持对多个手机号码发送短信，手机号码之间以半角逗号（,）分隔。上限为1000个手机号码。批量调用相对于单条调用及时性稍有延迟。
type sms struct {
	regionId     string
	accessKeyId  string
	accessSecret string
	signName     string
	request      *dysmsapi.SendSmsRequest
	phoneNumbers string
}

func New(regionId, accessKeyId, accessSecret string) *sms {
	return &sms{regionId: regionId, accessKeyId: accessKeyId, accessSecret: accessSecret, request: dysmsapi.CreateSendSmsRequest()}
}
func (s *sms) AddReceivers(phone string) {
	s.phoneNumbers = phone
}

func (s sms) Send(ctx context.Context, signName, templateCode, message string) error {
	var err error
	client, err := dysmsapi.NewClientWithAccessKey(s.regionId, s.accessKeyId, s.accessSecret)
	s.request.Scheme = "https"
	s.request.TemplateCode = templateCode
	s.request.SignName = signName
	s.request.PhoneNumbers = s.phoneNumbers
	s.request.TemplateParam = message
	select {
	case <-ctx.Done():
		err = ctx.Err()
	default:
		response, err := client.SendSms(s.request)
		if err != nil {
			err = errors.Wrap(err, "failed to send sms")
		}
		fmt.Printf("response is %v\n", response.Message)

	}
	return err

}
