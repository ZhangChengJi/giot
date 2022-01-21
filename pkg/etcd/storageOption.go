package etcd

import "reflect"

type StoreOption struct {
	BasePath   string
	ObjType    reflect.Type
	KeyFunc    func(obj interface{}) string
	StockCheck func(obj interface{}, stockObj interface{}) error
	//Validator  Validator
	//HubKey     HubKey
}
