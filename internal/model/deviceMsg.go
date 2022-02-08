package model

import (
	"giot/utils"
	"giot/utils/consts"
	"hash/fnv"
	"strings"
	"time"
)

type DeviceMsg struct {
	Ts         time.Time `json:"timestamp"`
	Type       string    `json:"type"`
	Status     bool      `json:"status"`
	DeviceId   string    `json:"deviceId"`
	Name       string    `json:"name" `
	SlaveId    int       `json:"slaveId"`
	ProductId  string    `json:"productId"`
	ModelId    string    `json:"modelId"`
	AlarmId    string    `json:"alarmId"`
	AlarmLevel int       `json:"alarmLevel"`
	Data       uint64    `json:"data"`
	NotifyType string    `json:"notifyType"` //通知类型
	TemplateId string    `json:"templateId"` //通知模版ID
	Actions    []*Action
}

// Change max partitions as you need.
const MAX_PARTITIONS = 10

// TaosEncoder implementations

// If this is setted, sql will use db.table for tablename
func (r DeviceMsg) TaosDatabase() string {
	return "dory_device"
}

// Auto create table using stable and tags
func (r DeviceMsg) TaosSTable() string {
	switch r.Type {
	case consts.DATA:
		return "device_data"
	case consts.ALARM:
		return "device_alarm"
	default:
		return ""

	}

}

// tags must be setted with TaosSTable
func (r DeviceMsg) TaosTags() []interface{} {
	var tags []interface{}
	if r.Type == consts.DATA {
		tags = append(tags, r.ProductId, r.DeviceId, r.SlaveId, r.ModelId)
	} else {
		tags = append(tags, r.ProductId, r.DeviceId, r.SlaveId, r.ModelId, r.AlarmId)
	}
	return tags
}

// Dynamic device id as table name
func (r DeviceMsg) TaosTable() string {
	switch r.Type {
	case consts.DATA:
		return strings.Join([]string{"device_data", r.DeviceId}, "")
	case consts.ALARM:
		return strings.Join([]string{"device_alarm", r.DeviceId}, "")
	default:
		return ""

	}
}

// Use specific column names as you need
func (r DeviceMsg) TaosCols() []string {
	var tags []string
	return tags
}

// Values
func (r DeviceMsg) TaosValues() []interface{} {
	var values []interface{}
	values = append(values, r.Ts)
	if r.Type == consts.DATA {
		values = append(values, r.Data, r.Status)
	} else {
		values = append(values, r.Data, r.AlarmLevel)
	}

	return values
}

// Codec interface

// Encoding method
func (r DeviceMsg) CodecMethod() utils.CodecMethodEnum {
	return utils.MessagePack
}

// How to set partition for an message
func (r DeviceMsg) Partition() int32 {
	h := fnv.New32a()
	h.Write([]byte(r.DeviceId))
	return int32(h.Sum32() % MAX_PARTITIONS)
}
