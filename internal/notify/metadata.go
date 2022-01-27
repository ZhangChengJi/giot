package notify

type Metadata struct {
	RegionId     string `json:"regionId"`
	AccessKeyId  string `json:"accessKeyId"`
	AccessSecret string `json:"accessSecret"`
	SignName     string `json:"signName"`
	TemplateCode string `json:"templateCode"`
	PhoneNumbers string `json:"phoneNumbers"`
}

type Template struct {
	DeviceName string
	SlaveId    int
	Value      float64
}
