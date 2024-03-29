package device

import (
	"giot/pkg/log"
	"giot/utils"
	"giot/utils/consts"
	"hash/fnv"
	"strconv"
	"strings"
	"time"
)

type DeviceMsg struct {
	Ts           time.Time `json:"ts"`
	DataType     string    `json:"dataType"`
	Level        int       `json:"level"`
	DeviceId     string    `json:"deviceId"`
	GroupId      int32     `json:"groupId"`
	Status       string    `json:"status"`
	Name         string    `json:"name" `
	SlaveId      int       `json:"slaveId"`
	SlaveName    string    `json:"slaveName"`
	Address      string    `json:"address"`
	Data         string    `json:"data"`
	Unit         string    `json:"unit"`
	PropertyName string    `json:"propertyName"`
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
		tags = append(tags, r.DeviceId, r.SlaveId, r.GroupId)
	} else {
		tags = append(tags, r.DeviceId, r.SlaveId, r.GroupId)
	}
	return tags
}

// Dynamic device id as table name
func (r DeviceMsg) TaosTable() string {
	switch r.DataType {
	case consts.DATA:
		return strings.Join([]string{"device_data", r.DeviceId, strconv.Itoa(r.SlaveId), strconv.Itoa(int(r.GroupId))}, "_")
	case consts.ALARM:
		return strings.Join([]string{"device_alarm", r.DeviceId, strconv.Itoa(r.SlaveId), strconv.Itoa(int(r.GroupId))}, "_")
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
		values = append(values, r.Data, r.Level, r.SlaveName, r.PropertyName, r.Unit)
	} else {
		values = append(values, r.Data, r.Level, r.SlaveName, r.PropertyName, r.Unit)
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
