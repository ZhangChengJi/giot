package device

import (
	"giot/pkg/log"
	"giot/utils"
	"giot/utils/consts"
	"hash/fnv"
	"strings"
	"time"
)

type DeviceMsg struct {
	Ts        time.Time `json:"timestamp"`
	DataType  string    `json:"dataType"`
	Level     int       `json:"level"`
	DeviceId  string    `json:"deviceId"`
	Status    string    `json:"status"`
	Name      string    `json:"name" `
	SlaveId   int       `json:"slaveId"`
	SlaveName string    `json:"slaveName"`
	Address   string    `json:"address"`
	Data      uint16    `json:"data"`
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
	switch r.DataType {
	case consts.DATA:
		return "device_data"
	case consts.ALARM:
		return "device_alarm"
	default:
		log.Sugar.Errorf("无法匹配到表")
		return ""

	}

}

// tags must be setted with TaosSTable
func (r DeviceMsg) TaosTags() []interface{} {
	var tags []interface{}
	if r.DataType == consts.DATA {
		tags = append(tags, r.DeviceId, r.SlaveId)
	} else {
		tags = append(tags, r.DeviceId, r.SlaveId)
	}
	return tags
}

// Dynamic device id as table name
func (r DeviceMsg) TaosTable() string {
	switch r.DataType {
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
	if r.DataType == consts.DATA {
		values = append(values, r.Data, r.Level)
	} else {
		values = append(values, r.Data, r.Level)
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
