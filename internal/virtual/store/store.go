package store

import (
	"giot/pkg/etcd"
	"reflect"
)

type Interface interface {
	//Type() HubKey
	//Get(ctx context.Context, key string) (json, error)
	//List(ctx context.Context, input ListInput) (*ListOutput, error)
	//Create(ctx context.Context, obj interface{}) (interface{}, error)
	//Update(ctx context.Context, obj interface{}, createIfNotExist bool) (interface{}, error)
	//BatchDelete(ctx context.Context, keys []string) error
}

type EtcdStore struct {
	etcd etcd.EtcdV3Storage
	opt  StoreOption
}

type StoreOption struct {
	BasePath   string
	ObjType    reflect.Type
	KeyFunc    func(obj interface{}) string
	StockCheck func(obj interface{}, stockObj interface{}) error
	//Validator  Validator
	//HubKey     HubKey
}
type code struct {
	EtcdStore
}

//
//func (e *EtcdStore) Get(ctx context.Context, key string, s reflect.Type) (interface{}, error) {
//	t := reflect.TypeOf(s)
//	data, _ := e.etcd.Get(ctx, key)
//	switch a := data.(type) {
//	case etcd.EtcdV3Storage:
//
//	}
//	t.
//
//}
