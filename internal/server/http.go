/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package server

//
//import (
//	"errors"
//	"fmt"
//	cmap "github.com/orcaman/concurrent-map"
//	"reflect"
//	"strconv"
//	"sync"
//)
//
//type name struct {
//	a int
//	b string
//}
//
//func (n *name) aasd() {
//
//}
//
//// ConcurrentMap 代表可自定义键类型和值类型的并发安全字典。
//type ConcurrentMap struct {
//	m         sync.Map
//	keyType   reflect.Type //键类型
//	valueType reflect.Type //值类型
//}
//
//func NewConcurrentMap(keyType, valueType reflect.Type) (*ConcurrentMap, error) {
//	if keyType == nil {
//		return nil, errors.New("nil key type")
//	}
//	if !keyType.Comparable() {
//		return nil, fmt.Errorf("incomparable key type: %s", keyType)
//	}
//	if valueType == nil {
//		return nil, errors.New("nil value type")
//	}
//	cMap := &ConcurrentMap{
//		keyType:   keyType,
//		valueType: valueType,
//	}
//	return cMap, nil
//}
//
//func (cMap *ConcurrentMap) Delete(key interface{}) {
//	if reflect.TypeOf(key) != cMap.keyType {
//		return
//	}
//	cMap.m.Delete(key)
//}
//
//func (cMap *ConcurrentMap) Load(key interface{}) (value interface{}, ok bool) {
//	if reflect.TypeOf(key) != cMap.keyType {
//		return
//	}
//	return cMap.m.Load(key)
//}
//
//func (cMap *ConcurrentMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
//	if reflect.TypeOf(key) != cMap.keyType {
//		panic(fmt.Errorf("wrong key type: %v", reflect.TypeOf(key)))
//	}
//	if reflect.TypeOf(value) != cMap.valueType {
//		panic(fmt.Errorf("wrong value type: %v", reflect.TypeOf(value)))
//	}
//	actual, loaded = cMap.m.LoadOrStore(key, value)
//	return
//}
//
//func (cMap *ConcurrentMap) Range(f func(key, value interface{}) bool) {
//	cMap.m.Range(f)
//}
//
//func (cMap *ConcurrentMap) Store(key, value interface{}) {
//	if reflect.TypeOf(key) != cMap.keyType {
//		panic(fmt.Errorf("wrong key type: %v", reflect.TypeOf(key)))
//	}
//	if reflect.TypeOf(value) != cMap.valueType {
//		panic(fmt.Errorf("wrong value type: %v", reflect.TypeOf(value)))
//	}
//	cMap.m.Store(key, value)
//}
//
//func main() {
//	// Create a new map.
//	m := cmap.New()
//
//	// Sets item within map, sets "bar" under key "foo"
//
//	for i := 0; i < 50; i++ {
//		m.Set("foo"+strconv.Itoa(i), &name{a: i, b: "tes2222"})
//	}
//	for i := 0; i < 50; i++ {
//
//		if tmp, ok := m.Get("foo" + strconv.Itoa(i)); ok {
//			bar := tmp.(*name)
//			fmt.Println(bar)
//		}
//	}
//	type s string
//	a := &ConcurrentMap{sync.Map{}, reflect.TypeOf(s), reflect.TypeOf(name{})}
//	as, _ := a.Load("123")
//	var ad = as.(name)
//	ad.aasd()
//	for {
//
//	}
//	// Retrieve item from map.
//	//if tmp, ok := m.Get("foo"); ok {
//	//	bar := tmp.(*name)
//	//	fmt.Println(bar)
//	//}
//	//
//	//// Removes item under key "foo"
//	//m.Remove("foo")
//}

//import (
//	"giot/internal/conf"
//	"giot/internal/filter"
//	"giot/internal/handler"
//	"net"
//	"net/http"
//	"strconv"
//	"time"
//
//	"github.com/shiningrush/droplet"
//)
//
//func (s *server) setupServer() {
//	// orchestrator
//	droplet.Option.Orchestrator = func(mws []droplet.Middleware) []droplet.Middleware {
//		var newMws []droplet.Middleware
//		// default middleware order: resp_reshape, auto_input, traffic_log
//		// We should put err_transform at second to catch all error
//		newMws = append(newMws, mws[0], &handler.ErrorTransformMiddleware{}, &filter.AuthenticationMiddleware{})
//		newMws = append(newMws, mws[1:]...)
//		return newMws
//	}
//
//	// routes
//	r := setupRouter()
//
//	// HTTP
//	addr := net.JoinHostPort(conf.ServerHost, strconv.Itoa(conf.ServerPort))
//	s.http = &http.Server{
//		Addr:           addr,
//		Handler:        r,
//		ReadTimeout:    time.Duration(1000) * time.Millisecond,
//		WriteTimeout:   time.Duration(5000) * time.Millisecond,
//		MaxHeaderBytes: 1 << 20,
//	}

// HTTPS
//if conf.SSLCert != "" && conf.SSLKey != "" {
//	addrSSL := net.JoinHostPort(conf.SSLHost, strconv.Itoa(conf.SSLPort))
//	s.serverSSL = &http.Server{
//		Addr:         addrSSL,
//		Handler:      r,
//		ReadTimeout:  time.Duration(1000) * time.Millisecond,
//		WriteTimeout: time.Duration(5000) * time.Millisecond,
//		TLSConfig: &tls.Config{
//			// Causes servers to use Go's default ciphersuite preferences,
//			// which are tuned to avoid attacks. Does nothing on clients.
//			PreferServerCipherSuites: true,
//		},
//	}
//}
//}
