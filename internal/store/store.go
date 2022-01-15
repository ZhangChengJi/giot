///*
// * Licensed to the Apache Software Foundation (ASF) under one or more
// * contributor license agreements.  See the NOTICE file distributed with
// * this work for additional information regarding copyright ownership.
// * The ASF licenses this file to You under the Apache License, Version 2.0
// * (the "License"); you may not use this file except in compliance with
// * the License.  You may obtain a copy of the License at
// *
// *     http://www.apache.org/licenses/LICENSE-2.0
// *
// * Unless required by applicable law or agreed to in writing, software
// * distributed under the License is distributed on an "AS IS" BASIS,
// * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// * See the License for the specific language governing permissions and
// * limitations under the License.
// */
//
package store

//
//import (
//	"context"
//	"fmt"
//	"giot/internal/log"
//	"giot/internal/storage"
//	"giot/utils/runtime"
//	"os"
//	"reflect"
//	"sync"
//	"time"
//
//	"github.com/shiningrush/droplet/data"
//)
//
//type Pagination struct {
//	PageSize   int `json:"page_size" form:"page_size" auto_read:"page_size"`
//	PageNumber int `json:"page" form:"page" auto_read:"page"`
//}
//
//type Interface interface {
//	Type() HubKey
//	Get(ctx context.Context, key string) (interface{}, error)
//	List(ctx context.Context, input ListInput) (*ListOutput, error)
//	Create(ctx context.Context, obj interface{}) (interface{}, error)
//	Update(ctx context.Context, obj interface{}, createIfNotExist bool) (interface{}, error)
//	BatchDelete(ctx context.Context, keys []string) error
//}
//
//type GenericStore struct {
//	Stg storage.Interface
//
//	cache sync.Map
//	opt   GenericStoreOption
//
//	cancel context.CancelFunc
//}
//
//type GenericStoreOption struct {
//	BasePath   string
//	ObjType    reflect.Type
//	KeyFunc    func(obj interface{}) string
//	StockCheck func(obj interface{}, stockObj interface{}) error
//	HubKey     HubKey
//}
//
//func NewGenericStore(opt GenericStoreOption) (*GenericStore, error) {
//	if opt.BasePath == "" {
//		log.Error("base path empty")
//		return nil, fmt.Errorf("base path can not be empty")
//	}
//	if opt.ObjType == nil {
//		log.Errorf("object type is nil")
//		return nil, fmt.Errorf("object type can not be nil")
//	}
//	if opt.KeyFunc == nil {
//		log.Error("key func is nil")
//		return nil, fmt.Errorf("key func can not be nil")
//	}
//
//	if opt.ObjType.Kind() == reflect.Ptr {
//		opt.ObjType = opt.ObjType.Elem()
//	}
//	if opt.ObjType.Kind() != reflect.Struct {
//		log.Error("obj type is invalid")
//		return nil, fmt.Errorf("obj type is invalid")
//	}
//	s := &GenericStore{
//		opt: opt,
//	}
//	s.Stg = storage.GenEtcdStorage()
//
//	return s, nil
//}
//
//func (s *GenericStore) Init() error {
//	lc, lcancel := context.WithTimeout(context.TODO(), 5*time.Second)
//	defer lcancel()
//	ret, err := s.Stg.List(lc, s.opt.BasePath)
//	if err != nil {
//		return err
//	}
//	for i := range ret {
//		key := ret[i].Key[len(s.opt.BasePath)+1:]
//		objPtr, err := s.StringToObjPtr(ret[i].Value, key)
//		if err != nil {
//			_, _ = fmt.Fprintln(os.Stderr, "Error occurred while initializing logical store: ", s.opt.BasePath)
//			return err
//		}
//
//		s.cache.Store(s.opt.KeyFunc(objPtr), objPtr)
//	}
//
//	c, cancel := context.WithCancel(context.TODO())
//	ch := s.Stg.Watch(c, s.opt.BasePath)
//	go func() {
//		defer runtime.HandlePanic()
//		for event := range ch {
//			if event.Canceled {
//				log.Warnf("watch failed: %s", event.Error)
//			}
//
//			for i := range event.Events {
//				switch event.Events[i].Type {
//				case storage.EventTypePut:
//					key := event.Events[i].Key[len(s.opt.BasePath)+1:]
//					objPtr, err := s.StringToObjPtr(event.Events[i].Value, key)
//					if err != nil {
//						log.Warnf("value convert to obj failed: %s", err)
//						continue
//					}
//					s.cache.Store(key, objPtr)
//				case storage.EventTypeDelete:
//					s.cache.Delete(event.Events[i].Key[len(s.opt.BasePath)+1:])
//				}
//			}
//		}
//	}()
//	s.cancel = cancel
//	return nil
//}
//
//func (s *GenericStore) Type() HubKey {
//	return s.opt.HubKey
//}
//
//func (s *GenericStore) Get(_ context.Context, key string) (interface{}, error) {
//	ret, ok := s.cache.Load(key)
//	if !ok {
//		log.Warnf("data not found by key: %s", key)
//		return nil, data.ErrNotFound
//	}
//	return ret, nil
//}
//
//type ListInput struct {
//	Predicate func(obj interface{}) bool
//	Format    func(obj interface{}) interface{}
//	PageSize  int
//	// start from 1
//	PageNumber int
//	Less       func(i, j interface{}) bool
//}
//
//type ListOutput struct {
//	Rows      []interface{} `json:"rows"`
//	TotalSize int           `json:"total_size"`
//}
//
//// NewListOutput returns JSON marshalling safe struct pointer for empty slice
//func NewListOutput() *ListOutput {
//	return &ListOutput{Rows: make([]interface{}, 0)}
//}
//
//func (s *GenericStore) GetObjStorageKey(obj interface{}) string {
//	return s.GetStorageKey(s.opt.KeyFunc(obj))
//}
//
//func (s *GenericStore) GetStorageKey(key string) string {
//	return fmt.Sprintf("%s/%s", s.opt.BasePath, key)
//}
