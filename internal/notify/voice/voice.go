package voice

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
)

type sms struct {
	regionId     string
	accessKeyId  string
	accessSecret string
	signName     string
	request      *dyvmsapi.SingleCallByVoiceRequest
	phoneNumbers string
}

//func New(regionId, accessKeyId, accessSecret string) *sms {
//	return &sms{regionId: regionId, accessKeyId: accessKeyId, accessSecret: accessSecret, request: dyvmsapi.CreateSingleCallByVoiceRequest()}
//}
//func (s *sms) AddReceivers(phone string) {
//	s.phoneNumbers = phone
//}
//
//func (s sms) Send(ctx context.Context, calledShowNumber, calledNumber, message string) error {
//	client, err := dyvmsapi.NewClientWithAccessKey(s.regionId, s.accessKeyId, s.accessSecret)
//	s.request.AcceptFormat = "json"
//	s.request.CalledShowNumber = calledShowNumber
//	s.request.CalledNumber = calledNumber
//	s.request.PhoneNumbers = s.phoneNumbers
//	s.request.TemplateParam = message
//}
