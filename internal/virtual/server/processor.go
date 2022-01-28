package server

import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	model2 "giot/internal/model"
	"giot/internal/virtual/device"
	"giot/internal/virtual/engine"
	"giot/internal/virtual/store"
	"giot/internal/virtual/wheelTimer"
	"giot/pkg/etcd"
	"giot/pkg/log"
	modbus2 "giot/pkg/modbus"
	"giot/utils/consts"
	"giot/utils/runtime"
	"github.com/RussellLuo/timingwheel"
	"github.com/panjf2000/gnet/pkg/pool/goroutine"
	"strconv"
	"strings"
	"time"
)

type Processor struct {
	modbus     modbus2.Client
	Stg        etcd.Interface
	Tw         *timingwheel.TimingWheel
	Dt         store.DeviceTimerIn
	al         store.AlarmStoreIn
	gu         store.GuidStoreIn
	sl         store.SlaveStoreIn
	workerPool *goroutine.Pool
}

func NewProcessor() *Processor {
	processor := &Processor{modbus: modbus2.NewClient(&modbus2.RtuHandler{}), Stg: etcd.GenEtcdStorage(), Tw: wheelTimer.NewTimer(), Dt: store.NewTimerStore(), al: store.NewAlarmStore(), gu: store.NewGuidStore(), sl: store.NewSlaveStore(), workerPool: goroutine.Default()}
	go processor.watchPoolEtcd()
	return processor
}

type ProcessorIn interface {
	Swift(data <-chan *model2.RemoteData, reg chan *model2.RegisterData)
	handle(data *model2.RemoteData)
	ListenCommand(msg <-chan *model2.ListenMsg)
	watchPoolEtcd()
	activeStore(action, guid, val string) error
	register(data *model2.RegisterData) error
	deleteTask(action, remoteAddr string)
}

func (p *Processor) Swift(rdata <-chan *model2.RemoteData, reg chan *model2.RegisterData) {

	for {
		select {
		case re := <-reg:
			err := p.register(re)
			if err != nil {
				log.Errorf("transfer guid:%s register error  remoteAddr:%s", string(re.D), re.C.RemoteAddr().String())
				return
			}
		case data := <-rdata:
			p.handle(data)
		case <-time.After(200 * time.Millisecond):
			//等待缓冲
		}
	}
}

func (p *Processor) handle(data *model2.RemoteData) {
	pdu, err := p.modbus.ReadCode(data.Frame) //解码
	if err != nil {
		log.Errorf("data Decode failed:%s", data.Frame)
		return
	}
	if slave, err := p.sl.GetSlave(context.TODO(), data.RemoteAddr, pdu.SaveId); err != nil { //获取属性ID
		log.Errorf("salve:%s not found", pdu.SaveId)
		return
	} else {
		da, _ := strconv.ParseFloat(string(data.Frame), 2)

		if al, err := p.al.Get(context.TODO(), data.RemoteAddr); err != nil {
			device.DataChan <- &model2.DeviceMsg{Type: consts.DATA, DeviceId: slave.DeviceId, ProductId: slave.ProductId, Name: slave.DeviceName, Status: true, Data: da, ModelId: slave.AttributeId, SlaveId: int(slave.SlaveId)}
			log.Warnf("remoteAddr:%s not alarm rule", data.RemoteAddr)
		} else {
			al.AlarmRule(da, slave)
		}

	}
}
func (p *Processor) ListenCommand(msg <-chan *model2.ListenMsg) {
	for {
		select {
		case m := <-msg:
			if m.ListenType == 1 { //1代表tcp 任务
				p.deleteTask(consts.ActionAll, m.RemoteAddr)

			} else {

			}
		case <-time.After(300 * time.Millisecond):
			//等待缓冲
		}
	}
}
func (p *Processor) deleteTask(action, remoteAddr string) {
	switch action {
	case consts.ActionAll:
		timer, err := p.Dt.Get(context.TODO(), remoteAddr)
		if err != nil {
			return
		}
		if len(timer) > 0 {
			p.gu.Delete(context.TODO(), timer[0].Guid) //guid删除
		}
		for _, t := range timer {
			t.T.Stop()
			t.Conn.Close() //TODO 强制关闭连接是否有必要?
		}
		p.Dt.Delete(context.TODO(), remoteAddr) //定时删除
		p.sl.Delete(context.TODO(), remoteAddr) //从机删除
		p.al.Delete(context.TODO(), remoteAddr) //告警删除
	case consts.ActionCode:
		p.Dt.Delete(context.TODO(), remoteAddr) //定时删除
	case consts.ActionSlave:
		p.sl.Delete(context.TODO(), remoteAddr) //从机删除
	case consts.ActionAlarm:
		p.al.Delete(context.TODO(), remoteAddr) //告警删除

	}

}
func (p *Processor) watchPoolEtcd() {
	c, cancel := context.WithCancel(context.TODO())
	ch := p.Stg.Watch(c, "transfer/")
	p.workerPool.Submit(func() {
		defer runtime.HandlePanic()
		defer cancel()
		for event := range ch {
			if event.Canceled {
				log.Warnf("watch failed: %s", event.Error)
			}
			for i := range event.Events {
				switch event.Events[i].Type {
				case etcd.EventTypePut:
					fmt.Println(event.Events[i].Key)
					fmt.Println(event.Events[i].Value)
					//key := event.Events[i].Key[len("transfer/"+guid):]
					ret := strings.Split(event.Events[i].Key, "/")
					p.activeStore(ret[2], ret[1], event.Events[i].Value)

					//key := event.Events[i].Key[len(s.opt.BasePath)+1:]
					//objPtr, err := s.StringToObjPtr(event.Events[i].Value, key)
					//if err != nil {
					//	log.Warnf("value convert to obj failed: %s", err)
					//	continue
					//}
					//s.cache.Store(key, objPtr)
				case etcd.EventTypeDelete:
					fmt.Println("delete...")
					ret := strings.Split(event.Events[i].Key, "/")
					guid, err := p.gu.Get(context.TODO(), ret[1])
					if err != nil {
						return
					}
					if ret[2] == consts.ActionCode {
						p.deleteTask(consts.ActionAll, guid)
					} else {
						p.deleteTask(ret[2], guid)
					}
				}
			}
		}
	})
}

func (p *Processor) activeStore(action, guid, val string) error {
	remoteAddr, err := p.gu.Get(context.TODO(), guid)
	if err != nil {
		log.Warnf("not found guid:%s Unable to query remoteAddr", guid)
		return err
	}
	switch action {
	case consts.ActionCode:
		timers, err := p.Dt.Get(context.TODO(), remoteAddr)
		if err != nil {
			log.Errorf("Unable to get remoteAddr:%s gnet.conn", remoteAddr)
			return err
		}
		for _, t := range timers {
			t.T.Stop()
		}
		de, err := metaDataCompile(val)
		if err != nil {
			log.Errorf("guid:%v transfer metadata transform error.", guid)
			return err
		}

		var timerList []*wheelTimer.SyncTimer
		//4. 封装定时器
		for _, v := range de.Ft {
			if len(v.FCode) > 0 {
				task := &wheelTimer.SyncTimer{
					Guid:       de.Guid,
					Conn:       timers[0].Conn,
					RemoteAddr: remoteAddr,
					Time:       v.Tm,
					Directives: v.FCode,
				}
				task.T = p.Tw.ScheduleFunc(&wheelTimer.DeviceScheduler{Interval: task.Time}, task.Execute)
				timerList = append(timerList, task)

			}
		}
		p.Dt.Update(context.TODO(), remoteAddr, timerList)
		s, _ := p.Dt.Get(context.TODO(), remoteAddr)
		for _, c := range s {
			fmt.Println(c)
		}

	case consts.ActionSlave:
		var slaves []*model2.Slave
		err = json.Unmarshal([]byte(val), &slaves)
		if err != nil {
			log.Errorf("json unmarshal failed: %s", err)
			return err
		}

		p.sl.Update(context.TODO(), remoteAddr, slaves)

	case consts.ActionAlarm:
		var alarms []*model2.Alarm
		err = json.Unmarshal([]byte(val), &alarms)
		if err != nil {
			log.Errorf("json unmarshal failed: %s", err)
			return err
		}
		if len(alarms) > 0 {
			alarmRule := engine.NewAlarmRule(alarms)
			p.al.Update(context.TODO(), remoteAddr, alarmRule)
		}
	}

	return nil
}

/**
  注册
*/
func (p *Processor) register(data *model2.RegisterData) error {
	//开始
	//1. 判断是否注册过，如果注册过无需重复注册
	remoteAddr := data.C.RemoteAddr().String()
	if wt, err := p.gu.Get(context.TODO(), string(data.D)); err == nil {
		if remoteAddr == wt {
			re, _ := p.modbus.WriteSingleRegister(1, 1, 1, modbus2.Success)
			data.C.AsyncWrite(re)
			log.Warnf("remoteAddr:%s alike no need to register again", remoteAddr)
			return err
		}
	}

	//2. etcd查询是否有元数据
	guid := string(data.D)
	val, err := p.Stg.Get(context.Background(), "transfer/"+guid+"/code")
	if err != nil {
		re, _ := p.modbus.WriteSingleRegister(1, 1, 1, modbus2.Error)
		data.C.AsyncWrite(re)
		log.Warnf("guid:%v transfer metadata not found.", guid)
		return err
	}
	//3. 认证成功开始配置元数据信息
	if len(val) > 0 {
		de, err := metaDataCompile(val)
		re, _ := p.modbus.WriteSingleRegister(1, 1, 1, modbus2.Success)
		data.C.AsyncWrite(re)
		if err != nil {
			log.Errorf("guid:%v transfer metadata transform error.", guid)
			return err
		}

		var timerList []*wheelTimer.SyncTimer
		//4. 封装定时器
		for _, v := range de.Ft {
			if len(v.FCode) > 0 {
				task := &wheelTimer.SyncTimer{
					Guid:       de.Guid,
					Conn:       data.C,
					RemoteAddr: remoteAddr,
					Time:       v.Tm,
					Directives: v.FCode,
				}
				task.T = p.Tw.ScheduleFunc(&wheelTimer.DeviceScheduler{Interval: task.Time}, task.Execute)
				timerList = append(timerList, task)

			}
		}
		//先存储一下guid和远程地址对应关系
		p.gu.Create(context.TODO(), guid, remoteAddr)
		//5. 获取从机信息
		sa, err := p.Stg.Get(context.Background(), "device/"+guid+"/salve")
		var slaves []*model2.Slave
		err = json.Unmarshal([]byte(sa), &slaves)
		if err != nil {
			log.Errorf("json unmarshal failed: %s", err)
			return err
		}

		//6. 获取告警规则
		tr, err := p.Stg.Get(context.Background(), "device/"+guid+"/alarm")
		var alarms []*model2.Alarm
		err = json.Unmarshal([]byte(tr), &alarms)
		if err != nil {
			log.Errorf("json unmarshal failed: %s", err)
			return err
		}
		if len(timerList) > 0 {
			p.Dt.Create(context.TODO(), remoteAddr, timerList)
		}
		if len(slaves) > 0 {
			p.sl.Create(context.TODO(), remoteAddr, slaves)
		}
		if len(alarms) > 0 {
			alarmRule := engine.NewAlarmRule(alarms)
			p.al.Create(context.TODO(), remoteAddr, alarmRule)
		}

		log.Infof("register on success,guid:%s remoteAddr:%s", data.D, data.C.RemoteAddr())
	}
	return nil
}
func metaDataCompile(val string) (*model2.TimerActive, error) {
	ma := &model2.TimerActive{}
	err := json.Unmarshal([]byte(val), ma)
	if err != nil {
		log.Errorf("json unmarshal failed: %s", err)
		return nil, err
	}
	return ma, nil
}
