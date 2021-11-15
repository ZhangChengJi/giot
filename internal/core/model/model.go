package model

type Product struct {
	Id                int16  `json:"id" xorm:"id"`
	ClassifiedId      int16  `json:"classifiedId" xorm:"classified_id"`           //分类ID
	OrgId             int16  `json:"orgId" xorm:"org_id"`                         //机构ID
	Name              string `json:"name" xorm:"name"`                            //产品名称
	ClassifiedName    string `json:"classifiedName" xorm:"classified_name"`       //产品分类
	MessageProtocol   string `json:"messageProtocol" xorm:"message_protocol"`     //消息协议 16hex  json
	TransportProtocol string `json:"transportProtocol" xorm:"transport_protocol"` //传输协议 http tcp mqtt
	DeviceType        string `json:"deviceType" xorm:"device_type"`               //设备类型 直连设备 网关设备  网关子设备
	Describe          string `json:"describe" xorm:"describe"`                    //描述

}

type Device struct {
	Name string `json:"name"`
}
