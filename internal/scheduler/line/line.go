package line

import (
	"bytes"
	"giot/utils/consts"
	"github.com/go-redis/redis"
)

const (
	on  = 1
	off = 0
)

type Line struct {
	Re *redis.Client
}
type LineCache interface {
	ClearAll()
	SetDeviceOnline(deviceId string)
	SetDeviceOffline(deviceId string)
	SetSlaveOnline(deviceId string, slaveId string)
	//BatchSlaveOnline(deviceId string,)
	SetSlaveOffline(deviceId string, slave string)
	BatchSlaveOffline(deviceId string)
}

func getWildcardCacheName(prefix string) string {
	var str bytes.Buffer
	str.WriteString(prefix)
	str.WriteString(consts.SYMBOL)
	str.WriteString(consts.WILDCARD)
	return str.String()
}

func getDeviceCacheName(prefix string) string {
	var str bytes.Buffer
	str.WriteString(consts.LINE_DEVICE)
	str.WriteString(consts.SYMBOL)
	str.WriteString(consts.WILDCARD)
	return str.String()
}

func (l *Line) ClearAll() {
	l.Re.Del(getWildcardCacheName(consts.LINE_DEVICE)) //删除设备在线状态
	l.Re.Del(getWildcardCacheName(consts.LINE_DEVICE)) //删除探头在线状态
}

// 设置设备在线
func (l *Line) SetDeviceOnline(deviceId string) {
	l.Re.HSet(consts.LINE_DEVICE, deviceId, on)
}

// 设置设备离线
func (l *Line) SetDeviceOffline(deviceId string) {
	l.Re.HSet(consts.LINE_DEVICE, deviceId, off)
}

// 设置探头在线
func (l *Line) SetSlaveOnline(deviceId string, slaveId string) {
	l.Re.HSet(consts.LINE_DEVICE+consts.SYMBOL+deviceId, slaveId, on)
}

//批量探头在线
//func (l *Line) BatchSlaveOnline(deviceId string) {
//		l.Re.HSet(consts.LINE_SLAVE+consts.SYMBOL+deviceId, slaveId, 1)
//
//}

// 设置探头离线
func (l *Line) SetSlaveOffline(deviceId string, slaveId string) {
	l.Re.HDel(consts.LINE_SLAVE+consts.SYMBOL+deviceId, slaveId)
}

//批量探头离线
func (l *Line) BatchSlaveOffline(deviceId string) {
	l.Re.Del(consts.LINE_SLAVE + consts.SYMBOL + deviceId)

}
