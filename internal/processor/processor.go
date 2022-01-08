package processor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"giot/internal/core/model"
	"giot/internal/log"
	"giot/internal/storage"
	"giot/internal/store"
	"github.com/panjf2000/gnet"
	"time"
	"unsafe"
)

var (
	RegisterChan chan RegisterData
	DataChan     chan RemoteData
)

type RemoteData struct {
	Frame    []byte
	RemoteIp []byte
}
type RegisterData struct {
	C gnet.Conn
	D []byte
}
type Processor struct {
	Stg storage.Interface
}

func Setup() {
	processor := NewProcessor()
	go processor.swift()
}

func NewProcessor() *Processor {
	DataChan = make(chan RemoteData, 1024)
	RegisterChan = make(chan RegisterData, 1024)
	return &Processor{Stg: storage.GenEtcdStorage()}

}

type ProcessorIn interface {
	swift()
	resolve(buf RemoteData)
	//register(buf []byte) error
	authenticate(data RegisterData) error
	handle(buf []byte) error
}

/**
注册
*/
func (p *Processor) register(data RegisterData) error {
	//开始
	guid := byToStr()
	c := data.C
	//1. 判断etcd里面是否有该配置信息
	if ok := store.IsDevice(store.WithPrefix(c.RemoteAddr().String())); ok {
		return errors.New(fmt.Sprintf("device:%sAlready exists,remoteIP:%s", string(buf.Frame), ip))
	}
	guid := string(buf.Frame)
	val, err := p.Stg.Get(context.Background(), guid)
	if err != nil {
		return err
	}
	//认证成功开始配置元数据信息
	if len(val) > 0 {
		de, err := metaDataCompile(val)
		if err != nil {
			return err
		}
		store.NewDevice(store.WithPrefix(ip), de)
	}
	return nil
}
func metaDataCompile(val string) (*model.Device, error) {
	ma := &model.Device{}
	err := json.Unmarshal([]byte(val), ma)
	if err != nil {
		log.Errorf("json unmarshal failed: %s", err)
		return nil, err
	}
	return ma, nil
}
func (p *Processor) resolve(buf RemoteData) {
	length := len(buf.Frame)
	if length > 0 && length > 7 {
		if length == 24 {
			err := p.register()
			if err != nil {
				return
			}
		} else {
			<
		}
	}

}

func (p *Processor) swift() {

	for {

		select {
		case re := <-RegisterChan:
			p.register(re)
		case data := <-DataChan:

		case <-time.After(200 * time.Millisecond):
			//等待缓冲
		}
	}

}

func byToStr(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}

var (
	maa map[string]ssd
)

type a struct {
	name string
}
type ssd interface {
	ss()
}

func (b *a) ss() {
	fmt.Println(b.name)
}

func d() {

	maa["a"] = &a{name: "小屋"}
	s := maa["a"]
	s.ss()
}
