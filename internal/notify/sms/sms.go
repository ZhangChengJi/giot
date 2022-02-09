package sms

import (
	"context"
	"fmt"
	"giot/pkg/log"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/pkg/errors"
)

//接收短信的手机号码。
//
//手机号码格式：
//国内短信：+/+86/0086/86或无任何前缀的11位手机号码，例如1590000****。
//国际/港澳台消息：国际区号+号码，例如852000012****。
//支持对多个手机号码发送短信，手机号码之间以半角逗号（,）分隔。上限为1000个手机号码。批量调用相对于单条调用及时性稍有延迟。
//https://next.api.aliyun.com/api/Dysmsapi/2017-05-25/SendSms?spm=api-workbench.SDK%20Document.0.0.edf51e0fzudjW8&lang=GO&params={%22PhoneNumbers%22:%2218866890352%22,%22SignName%22:%22%E5%A4%9A%E7%91%9E%E4%BA%91%22,%22TemplateCode%22:%22SMS_232169525%22,%22TemplateParam%22:%22{\%22msisdn\%22:\%2218923022\%22,\%22name\%22:\%22%E5%93%88%E5%95%8A\%22,\%22date\%22:\%222022\%22}%22}&tab=DOC
type sms struct {
	client       *dysmsapi20170525.Client
	phoneNumbers string
}
type Template struct {
	Devname   string `json:"devname"`
	Devid     string `json:"devid"`
	Alarmtype string `json:"alarmtype"`
}

func New(accessKeyId, accessKeySecret string) *sms {
	config := &openapi.Config{
		Endpoint: tea.String("dysmsapi.aliyuncs.com"),
		// 您的AccessKey ID
		AccessKeyId: tea.String(accessKeyId),
		// 您的AccessKey Secret
		AccessKeySecret: tea.String(accessKeySecret),
	}
	_result, _err := dysmsapi20170525.NewClient(config)
	if _err != nil {
		log.Errorf("init aliyun sms failed:%s", _err)
	}
	return &sms{client: _result}
}
func (s *sms) AddReceivers(phone string) {
	s.phoneNumbers = phone
}

func (s *sms) Send(ctx context.Context, signName, templateCode, param string) (err error) {
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(s.phoneNumbers),
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(templateCode),
		TemplateParam: tea.String(param),
	}
	select {
	case <-ctx.Done():
		err = ctx.Err()
	default:
		resp, err := s.client.SendSms(sendSmsRequest)
		if err != nil {
			err = errors.Wrap(err, "failed to send sms")
		}
		log.Errorf("send cms failed:%v", resp.Body)
		fmt.Printf("response is %v\n", resp.Body)

	}
	return err

}
