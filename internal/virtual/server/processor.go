package server

import (
	"context"
	"encoding/json"
	"fmt"
	"giot/internal/model"
	"giot/internal/virtual/data"
	"giot/internal/virtual/device"
	"giot/internal/virtual/mqtt"
	"giot/internal/virtual/store"
	"giot/internal/virtual/wheelTimer"
	"giot/pkg/etcd"
	"giot/pkg/log"
	"giot/pkg/modbus"
	"giot/utils/consts"
	"giot/utils/runtime"
	"github.com/RussellLuo/timingwheel"
	mqtt1 "github.com/eclipse/paho.mqtt.golang"
	"github.com/panjf2000/gnet/pkg/pool/goroutine"
	log4j "log"
	"strings"
	"time"
)

var (
	QuitTask chan *model.ListenMsg
)

type Processor struct {
	modbus     modbus.Client
	Stg        etcd.Interface
	Timer      *timingwheel.TimingWheel
	workerPool *goroutine.Pool
	cache      *store.CacheStore
	data       *data.Data
	mq         mqtt1.Client
}

func NewProcessor() *Processor {
	processor := &Processor{modbus: modbus.NewClient(&modbus.RtuHandler{}),
		Stg:        etcd.GenEtcdStorage(),
		Timer:      wheelTimer.NewTimer(),
		cache:      store.New(),
		data:       data.New(),
		workerPool: goroutine.Default(),
		mq:         mqtt.Client,
	}
	store.NewLine(processor.cache, mqtt.Broker{Client: mqtt.Client})
	go processor.watchDeviceChange()
	go processor.debug()
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
	clearConnect(guid, remoteAddr string, duration time.Duration)
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
			info, err := p.cache.Device.Get(context.TODO(), data.RemoteAddr)
			if err == nil && info != nil {
				results, err := p.protocol(info, data.RemoteAddr, data.Frame)
				if err == nil {
					for _, result := range results {
						p.alarmFilter(data.RemoteAddr, result, info)
					}
				}
			} else {
				re, _ := p.modbus.WriteSingleRegister(1, 1, 1, modbus.Error)
				data.Conn.AsyncWrite(re)
				data.Conn.Close()
				fmt.Println("未注册强制断开连接")
			}
		case <-time.After(200 * time.Millisecond):
			//等待缓冲
		}
	}
}

func (p *Processor) protocol(info *model.Device, remoteAddr string, frame []byte) (results []*modbus.ProtocolDataUnit, err error) {

	fmt.Printf("时间:%v——--->指令上报:%X\n", time.Now().Format("2006-01-02 15:04:05"), frame)
	if info.IsType() { //是否是工业产品
		if info.IsInstruct() { //是否是单指令下发
			results, err = p.modbus.ReadIndustryF1Code(frame)
		} else {
			result, err := p.modbus.ReadIndustryCode(frame) //解码
			if err == nil {
				results = append(results, result)
			}
		}

	} else {

		//} else {
		//	result, err = p.modbus.ReadHomeCode(frame) //解码
		//results = append(results, result)

	}
	return results, nil
}

func (p *Processor) alarmFilter(remoteAddr string, result *modbus.ProtocolDataUnit, info *model.Device) {
	if slave, err := p.cache.Slave.GetSlave(context.TODO(), remoteAddr, result.SlaveId); err == nil { //获取属性ID
		slave.DataTime = time.Now()
		//第一次发送上线通知
		if slave.LineStatus == "" || slave.LineStatus == consts.OFFLINE {
			fmt.Printf("时间:%v----->slave:%v上⬆️线\n", time.Now().Format("2006-01-02 15:04:05"), slave.SlaveId)
			slave.LineStatus = consts.ONLINE
			device.OnlineChan <- &device.DeviceMsg{Ts: time.Now(), Status: consts.ONLINE, DeviceId: info.GuId, SlaveId: int(slave.SlaveId)}
		}
		//
		slave.Rule.AlarmRule(slave, result.Data, result.FunctionCode, info)
	} else {
		log.Sugar.Errorf("salve:%v not found", result.SlaveId)
	}
}
func (p *Processor) ListenCommand(msg chan *model.ListenMsg) {
	for {
		select {
		case m := <-msg:
			if m.ListenType == 1 { //1代表tcp 任务
				p.DeleteTask(m.RemoteAddr)

			} else {
			}
		case <-time.After(300 * time.Millisecond):
			//等待缓冲
		}
	}
}
func (p *Processor) DeleteTask(remoteAddr string) {
	log.Sugar.Warnf("清理设备任务，远程地址：%v", remoteAddr)
	timer, err := p.cache.Timer.Get(context.TODO(), remoteAddr)
	if err != nil {
		return
	}
	log.Sugar.Warnf("清理设备任务，guid：%v", timer.Guid)

	if timer != nil {
		p.cache.Guid.Delete(context.TODO(), timer.Guid) //远程地址和guid对应关系删除
	}
	timer.T.Stop()
	timer.Conn.Close() //TODO 强制关闭连接是否有必要?
	lineTimer, err := p.cache.Line.Get(context.TODO(), remoteAddr)
	if err != nil {
		return
	}
	lineTimer.T.Stop()
	p.cache.Line.Delete(context.TODO(), remoteAddr)   //上下线任务检测删除
	p.cache.Timer.Delete(context.TODO(), remoteAddr)  //定时删除
	p.cache.Slave.Delete(context.TODO(), remoteAddr)  //从机删除
	p.cache.Device.Delete(context.TODO(), remoteAddr) //设备数据删除
	device.OnlineChan <- &device.DeviceMsg{Ts: time.Now(), Status: consts.OFFLINE, DeviceId: timer.Guid}

}
func (p *Processor) watchDeviceChange() {
	c, cancel := context.WithCancel(context.TODO())
	key := strings.Join([]string{consts.METADATA}, consts.DIVIDER)

	ch := p.Stg.Watch(c, key)
	p.workerPool.Submit(func() {
		defer runtime.HandlePanic()
		defer cancel()
		for event := range ch {
			if event.Canceled {
				log.Sugar.Warnf("watch failed: %s", event.Error)
			}
			for i := range event.Events {
				switch event.Events[i].Type {
				case etcd.EventTypePut:
					log.Sugar.Infof("etcd device data key:%v ,update...", event.Events[i].Key)
					fmt.Println(event.Events[i].Value)
					//key := event.Events[i].Key[len("transfer/"+guid):]
					//giot/device/296424434E48313836FFD805/code
					ret := strings.Split(event.Events[i].Key, "/")
					p.activeStore(ret[2], ret[3], event.Events[i].Value)

					//key := event.Events[i].Key[len(s.opt.BasePath)+1:]
					//objPtr, err := s.StringToObjPtr(event.Events[i].Value, key)
					//if err != nil {
					//	logs.Warnf("value convert to obj failed: %s", err)
					//	continue
					//}
					//s.cache.Store(key, objPtr)
				case etcd.EventTypeDelete:
					ret := strings.Split(event.Events[i].Key, "/")
					remoteAddr, err := p.cache.Guid.Get(context.TODO(), ret[2])
					if err != nil {
						return
					}
					log.Sugar.Infof("etcd device data key:%v ,delete...", event.Events[i].Key)

					p.DeleteTask(remoteAddr)
				}
			}
		}
	})
}

func (p *Processor) activeStore(action, guid, val string) {
	remoteAddr, err := p.cache.Guid.Get(context.TODO(), guid)
	if err != nil {
		return
	}
	if action == "device" {
		var devcie *model.Device
		err := json.Unmarshal([]byte(val), &devcie)
		if err != nil {
			return
		}
		p.cache.Device.Update(context.TODO(), remoteAddr, devcie)

	} else if action == "slave" {

		var slaveData []*model.Slave
		json.Unmarshal([]byte(val), &slaveData)
		p.cache.Slave.Update(context.TODO(), remoteAddr, slaveData)
		log.Sugar.Infof("更新了slave%v", guid)
	}
	timer, err := p.data.GetTimerData(guid)
	if err != nil {
		return
	}
	timers, err := p.cache.Timer.Get(context.TODO(), remoteAddr)
	if err != nil {
		log.Sugar.Errorf("Unable to get remoteAddr:%s gnet.conn", remoteAddr)
		return
	}
	timers.T.Stop()
	if timer != nil {
		task := &wheelTimer.SyncTimer{
			Guid:       guid,
			Conn:       timers.Conn,
			RemoteAddr: remoteAddr,
			Time:       timer.Ft.Tm,
			Directives: timer.Ft.FCode,
		}
		task.T = p.Timer.ScheduleFunc(&wheelTimer.DeviceScheduler{Interval: task.Time, Rew: guid}, task.Execute)
		p.cache.Timer.Update(context.TODO(), remoteAddr, task)
		log.Sugar.Infof("更新了timer%v", guid)
	}

}

/**
  注册
*/
func (p *Processor) register(data *model.RegisterData) {
	//开始
	//1. 判断是否注册过，如果注册过无需重复注册
	remoteAddr := data.Conn.RemoteAddr().String()
	//
	if wt, err := p.cache.Guid.Get(context.TODO(), data.D); err == nil {
		if remoteAddr == wt {
			re, _ := p.modbus.WriteSingleRegister(1, 1, 1, modbus.Success)
			data.Conn.AsyncWrite(re)
			log.Sugar.Warnf("guid:%s已经注册 remoteAddr:%s 无需再次注册", data.D, remoteAddr)

			return
		} else {
			log.Sugar.Infof("guid:%s 已存在连接:%s 强制下线，新连接:%s上线", data.D, wt, remoteAddr)
			p.DeleteTask(wt)
		}
	}
	//没有注册过就etcd查询
	//2. etcd查询是否有元数据
	guid := string(data.D)
	val, err := p.data.GetData(guid)
	if err != nil {
		re, _ := p.modbus.WriteSingleRegister(1, 1, 1, modbus.Error)
		data.Conn.AsyncWrite(re)
		data.Conn.Close()
		log.Sugar.Warnf("guid:%v metadata not found.", guid)
		log4j.Printf("guid:%v remoteAddr:%v注册失败，无法查询到元数据", guid, remoteAddr)
		return
	}
	var task *wheelTimer.SyncTimer
	//3. 认证成功开始配置元数据信息
	if val != nil {
		//de, err := metaDataCompile(val)
		re, _ := p.modbus.WriteSingleRegister(1, 1, 1, modbus.Success)
		data.Conn.AsyncWrite(re)
		if err != nil {
			log.Sugar.Errorf("guid:%v transfer metadata transform error.", guid)
			return
		}

		//4. 封装定时器

		if val.FCode != nil {
			task = &wheelTimer.SyncTimer{
				Guid:       val.GuId,
				Conn:       data.Conn,
				RemoteAddr: remoteAddr,
				Time:       val.FCode.Tm,
				Directives: val.FCode.FCode,
			}
			task.T = p.Timer.ScheduleFunc(&wheelTimer.DeviceScheduler{Interval: task.Time, Rew: val.GuId}, task.Execute)
		}
		if val.FCode != nil {
			p.cache.Timer.Create(context.TODO(), remoteAddr, task)
		}
		line := &store.LineTimer{
			Guid:   val.GuId,
			Status: 1,
			Time:   60 * time.Second,
		}
		line.T = p.Timer.ScheduleFunc(&wheelTimer.DeviceScheduler{Interval: 60 * time.Second}, line.Execute)
		//先存储一下guid和远程地址对应关系
		p.cache.Guid.Create(context.TODO(), guid, remoteAddr)
		deviceInfo := &model.Device{
			GuId:         val.GuId,
			Name:         val.Name,
			ProductType:  val.ProductType,
			ProductModel: val.ProductModel,
			Instruct:     val.Instruct,
			LineStatus:   val.LineStatus,
			GroupId:      val.GroupId,
			Address:      val.Address,
		}

		if line != nil {
			p.cache.Line.Create(context.TODO(), remoteAddr, line)
		}
		if len(val.Salve) > 0 {
			p.cache.Slave.Create(context.TODO(), remoteAddr, val.Salve)
		}
		if deviceInfo != nil {
			p.cache.Device.Create(context.TODO(), remoteAddr, deviceInfo)
		}
		device.OnlineChan <- &device.DeviceMsg{Ts: time.Now(), Status: consts.ONLINE, DeviceId: val.GuId}
		log.Sugar.Infof("guid:%v remoteAddr:%v注册成功🧸", guid, remoteAddr)

	}
	return
}
func (p *Processor) debug() {
	for {
		select {
		case de := <-device.DebugChan:
			addr, err := p.cache.Guid.Get(context.TODO(), de.Guid)
			if err == nil {
				if we, err := p.cache.Timer.Get(context.TODO(), addr); err == nil {
					log.Sugar.Infof("debug指令下发：%v", de.FCode)
					we.Conn.AsyncWrite(de.FCode)
				}
			}

		case <-time.After(300 * time.Millisecond):
		}
	}
}
func metaDataCompile(val string) (*model.Device, error) {
	devic := &model.Device{}
	err := json.Unmarshal([]byte(val), devic)
	if err != nil {
		log.Sugar.Errorf("json unmarshal failed: %s", err)
		return nil, err
	}
	return devic, nil
}
