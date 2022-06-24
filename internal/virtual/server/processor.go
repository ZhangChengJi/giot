package server

import (
	"context"
	"encoding/json"
	"fmt"
	"giot/internal/model"
	"giot/internal/virtual/device"
	"giot/internal/virtual/lineTimer"
	"giot/internal/virtual/store"
	"giot/internal/virtual/wheelTimer"
	"giot/pkg/etcd"
	"giot/pkg/log"
	modbus2 "giot/pkg/modbus"
	"giot/utils/consts"
	"giot/utils/runtime"
	"github.com/RussellLuo/timingwheel"
	"github.com/panjf2000/gnet/pkg/pool/goroutine"
	log4j "log"
	"strings"
	"time"
)

type Processor struct {
	modbus      modbus2.Client
	Stg         etcd.Interface
	Timer       *timingwheel.TimingWheel
	TimerStore  store.DeviceTimerIn
	LineStore   lineTimer.LineStoreIn
	guidStore   store.GuidStoreIn
	slaveStore  store.SlaveStoreIn
	deviceStore store.DeviceStoreIn
	workerPool  *goroutine.Pool
}

func NewProcessor() *Processor {
	processor := &Processor{modbus: modbus2.NewClient(&modbus2.RtuHandler{}), Stg: etcd.GenEtcdStorage(), Timer: wheelTimer.NewTimer(), TimerStore: store.NewTimerStore(), LineStore: lineTimer.NewLineStore(), guidStore: store.NewGuidStore(), slaveStore: store.NewSlaveStore(), deviceStore: store.NewDeviceStore(), workerPool: goroutine.Default()}
	go processor.watchPoolEtcd()
	return processor

}

type ProcessorIn interface {
	Swift(reg chan *model.RegisterData)
	Handle(da chan *model.RemoteData)
	ListenCommand(msg chan *model.ListenMsg)
	watchPoolEtcd()
	activeStore(guid, val string) error
	register(data *model.RegisterData) error
	deleteTask(remoteAddr string)
	checkLine(guid, remoteAddr string, duration time.Duration)
}

func (p *Processor) Swift(reg chan *model.RegisterData) {

	for {
		select {
		case re := <-reg:
			p.register(re)

		case <-time.After(200 * time.Millisecond):
			//等待缓冲
		}
	}
}

func (p *Processor) Handle(da chan *model.RemoteData) {

	for {
		select {
		case data := <-da:
			info, err := p.deviceStore.Get(context.TODO(), data.RemoteAddr)
			if err == nil && info != nil {
				var pdu *modbus2.ProtocolDataUnit
				if info.IsType() {
					pdu, err = p.modbus.ReadIndustryCode(data.Frame) //解码
				} else {
					pdu, err = p.modbus.ReadHomeCode(data.Frame) //解码
				}
				if err == nil {
					if slave, err := p.slaveStore.GetSlave(context.TODO(), data.RemoteAddr, pdu.SaveId); err == nil { //获取属性ID
						slave.DataTime = time.Now()
						if slave.LineStatus == "" || slave.LineStatus == consts.OFFLINE {
							fmt.Printf("时间:%v----->slave:%v上⬆️线\n", time.Now().Format("2006-01-02 15:04:05"), slave.SlaveId)
							slave.LineStatus = consts.ONLINE
							device.OnlineChan <- &device.DeviceMsg{Ts: time.Now(), Status: consts.ONLINE, DeviceId: info.GuId, SlaveId: int(slave.SlaveId)}
						}
						slave.Alarm.AlarmRule(slave.SlaveId, ByteToFloat32(pdu.Data), pdu.FunctionCode, info)
					} else {
						log.Errorf("salve:%s not found", pdu.SaveId)
					}
				}
			}
		case <-time.After(200 * time.Millisecond):
			//等待缓冲
		}
	}

}

func ByteToFloat32(bytes []byte) float32 {
	b2 := bytes[1]
	return float32(b2)
}
func (p *Processor) ListenCommand(msg chan *model.ListenMsg) {
	for {
		select {
		case m := <-msg:
			if m.ListenType == 1 { //1代表tcp 任务
				p.deleteTask(m.RemoteAddr)

			} else {
			}
		case <-time.After(300 * time.Millisecond):
			//等待缓冲
		}
	}
}
func (p *Processor) deleteTask(remoteAddr string) {

	timer, err := p.TimerStore.Get(context.TODO(), remoteAddr)
	if err != nil {
		return
	}
	if timer != nil {
		p.guidStore.Delete(context.TODO(), timer.Guid) //远程地址和guid对应关系删除
	}

	timer.T.Stop()
	timer.Conn.Close() //TODO 强制关闭连接是否有必要?
	line, err := p.LineStore.Get(context.TODO(), remoteAddr)
	if err != nil {
		return
	}
	line.T.Stop()
	p.TimerStore.Delete(context.TODO(), remoteAddr)  //定时删除
	p.LineStore.Delete(context.TODO(), remoteAddr)   //删除slave在线检查
	p.slaveStore.Delete(context.TODO(), remoteAddr)  //从机删除
	p.deviceStore.Delete(context.TODO(), remoteAddr) //设备数据删除
	device.OnlineChan <- &device.DeviceMsg{Ts: time.Now(), Status: consts.OFFLINE, DeviceId: timer.Guid}

}
func (p *Processor) watchPoolEtcd() {
	c, cancel := context.WithCancel(context.TODO())
	ch := p.Stg.Watch(c, "device/")
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
					log.Infof("etcd device data key:%v ,update...", event.Events[i].Key)
					fmt.Println(event.Events[i].Value)
					//key := event.Events[i].Key[len("transfer/"+guid):]
					//giot/device/296424434E48313836FFD805/code
					ret := strings.Split(event.Events[i].Key, "/")
					p.activeStore(ret[2], event.Events[i].Value)

					//key := event.Events[i].Key[len(s.opt.BasePath)+1:]
					//objPtr, err := s.StringToObjPtr(event.Events[i].Value, key)
					//if err != nil {
					//	log.Warnf("value convert to obj failed: %s", err)
					//	continue
					//}
					//s.cache.Store(key, objPtr)
				case etcd.EventTypeDelete:
					ret := strings.Split(event.Events[i].Key, "/")
					remoteAddr, err := p.guidStore.Get(context.TODO(), ret[2])
					if err != nil {
						return
					}
					log.Infof("etcd device data key:%v ,delete...", event.Events[i].Key)

					p.deleteTask(remoteAddr)
				}
			}
		}
	})
}

func (p *Processor) activeStore(guid, val string) error {
	remoteAddr, err := p.guidStore.Get(context.TODO(), guid)
	if err != nil {
		log.Warnf("not found guid:%s Unable to query remoteAddr", guid)
		return err
	}
	de, err := metaDataCompile(val)

	if err != nil {
		log.Errorf("guid:%v transfer metadata transform error.", guid)
		return err
	}

	timers, err := p.TimerStore.Get(context.TODO(), remoteAddr)
	if err != nil {
		log.Errorf("Unable to get remoteAddr:%s gnet.conn", remoteAddr)
		return err
	}
	timers.T.Stop()
	line, err := p.LineStore.Get(context.TODO(), remoteAddr)
	if err != nil {
		return err
	}
	line.T.Stop()
	if de.FCode != nil {
		task := &wheelTimer.SyncTimer{
			Guid:       de.GuId,
			Conn:       timers.Conn,
			RemoteAddr: remoteAddr,
			Time:       de.FCode.Tm,
			Directives: de.FCode.FCode,
		}
		task.T = p.Timer.ScheduleFunc(&wheelTimer.DeviceScheduler{Interval: task.Time}, task.Execute)
		p.TimerStore.Update(context.TODO(), remoteAddr, task)
	}

	lineTask := &lineTimer.LineTimer{
		Guid:       guid,
		RemoteAddr: remoteAddr,
		Time:       de.FCode.Tm + 15*time.Second,
		SlaveStore: p.slaveStore,
	}
	lineTask.T = p.Timer.ScheduleFunc(&wheelTimer.DeviceScheduler{Interval: lineTask.Time}, lineTask.Execute)
	p.LineStore.Update(context.TODO(), remoteAddr, lineTask)

	p.slaveStore.Update(context.TODO(), remoteAddr, de.Salve)

	deviceInfo := &model.Device{
		GuId:         de.GuId,
		Name:         de.Name,
		ProductType:  de.ProductType,
		ProductModel: de.ProductModel,
	}
	p.deviceStore.Update(context.TODO(), remoteAddr, deviceInfo)

	return nil
}

/**
  注册
*/
func (p *Processor) register(data *model.RegisterData) {
	//开始
	//1. 判断是否注册过，如果注册过无需重复注册
	remoteAddr := data.Conn.RemoteAddr().String()
	//
	if wt, err := p.guidStore.Get(context.TODO(), data.D); err == nil {
		if remoteAddr == wt {
			re, _ := p.modbus.WriteSingleRegister(1, 1, 1, modbus2.Success)
			data.Conn.AsyncWrite(re, nil)
			log.Warnf("remoteAddr:%s alike no need to register again", remoteAddr)

			return
		}
	}
	//没有注册过就etcd查询
	//2. etcd查询是否有元数据
	guid := string(data.D)
	val, err := p.Stg.Get(context.Background(), "device/"+guid)
	if err != nil {
		re, _ := p.modbus.WriteSingleRegister(1, 1, 1, modbus2.Error)
		data.Conn.AsyncWrite(re, nil)
		data.Conn.Close()
		log.Warnf("guid:%v metadata not found.", guid)
		log4j.Printf("guid:%v remoteAddr:%v注册失败，无法查询到元数据", guid, remoteAddr)
		return
	}
	var task *wheelTimer.SyncTimer
	//3. 认证成功开始配置元数据信息
	if len(val) > 0 {
		de, err := metaDataCompile(val)
		re, _ := p.modbus.WriteSingleRegister(1, 1, 1, modbus2.Success)
		data.Conn.AsyncWrite(re, nil)
		if err != nil {
			log.Errorf("guid:%v transfer metadata transform error.", guid)
			return
		}

		//4. 封装定时器

		if de.FCode != nil {
			task = &wheelTimer.SyncTimer{
				Guid:       de.GuId,
				Conn:       data.Conn,
				RemoteAddr: remoteAddr,
				Time:       de.FCode.Tm,
				Directives: de.FCode.FCode,
			}
			task.T = p.Timer.ScheduleFunc(&wheelTimer.DeviceScheduler{Interval: task.Time}, task.Execute)
		}

		lineTask := &lineTimer.LineTimer{
			Guid:       guid,
			RemoteAddr: remoteAddr,
			Time:       de.FCode.Tm + 15*time.Second,
			SlaveStore: p.slaveStore,
		}
		lineTask.T = p.Timer.ScheduleFunc(&wheelTimer.DeviceScheduler{Interval: lineTask.Time}, lineTask.Execute)

		//先存储一下guid和远程地址对应关系
		p.guidStore.Create(context.TODO(), guid, remoteAddr)
		deviceInfo := &model.Device{
			GuId:         de.GuId,
			Name:         de.Name,
			ProductType:  de.ProductType,
			ProductModel: de.ProductModel,
		}
		if de.FCode != nil {
			p.TimerStore.Create(context.TODO(), remoteAddr, task)
		}
		if lineTask != nil {
			p.LineStore.Create(context.TODO(), remoteAddr, lineTask)
		}
		if len(de.Salve) > 0 {
			p.slaveStore.Create(context.TODO(), remoteAddr, de.Salve)
		}
		if deviceInfo != nil {
			p.deviceStore.Create(context.TODO(), remoteAddr, deviceInfo)
		}
		device.OnlineChan <- &device.DeviceMsg{Ts: time.Now(), Status: consts.ONLINE, DeviceId: de.GuId}
		log.Infof("register on success,guid:%s remoteAddr:%s", data.D, data.Conn.RemoteAddr())
		log4j.Printf("guid:%v remoteAddr:%v注册成功🧸", guid, remoteAddr)

	}
	return
}

func metaDataCompile(val string) (*model.Device, error) {
	devic := &model.Device{}
	err := json.Unmarshal([]byte(val), devic)
	if err != nil {
		log.Errorf("json unmarshal failed: %s", err)
		return nil, err
	}
	return devic, nil
}
