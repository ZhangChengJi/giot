package notify

type Metadata struct {
	AccessKeyId  string `json:"accessKeyId"`
	AccessSecret string `json:"accessSecret"`
	Sms          *SmsMetadata
	Voice        *VoiceMetadata
}
type SmsMetadata struct {
	SignName    string `json:"signName"`
	Code        string `json:"code"`
	PhoneNumber string `json:"phoneNumber"`
	Param       string `json:"param"`
}
type VoiceMetadata struct {
}

type Template struct {
	DeviceName string
	SlaveId    int
	Value      uint16
}
