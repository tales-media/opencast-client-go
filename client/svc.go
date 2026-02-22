/*
Copyright 2025 shio solutions GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"context"
	"errors"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"shio.solutions/tales.media/opencast-client-go/apis/meta/objlist"
	"shio.solutions/tales.media/opencast-client-go/apis/meta/strobj"
	"shio.solutions/tales.media/opencast-client-go/apis/serviceregistry"
)

var ServiceNotFoundErr = errors.New("service not found")

type ServiceMapper interface {
	GetHost(svc string) (string, error)
}

type StaticServiceMapper struct {
	Default     string
	ServiceHost map[string]string
}

var _ ServiceMapper = &StaticServiceMapper{}

func (m *StaticServiceMapper) GetHost(svc string) (string, error) {
	if host, ok := m.ServiceHost[svc]; ok {
		return host, nil
	}
	return m.Default, nil
}

type dynamicServiceMapper struct {
	occ     Client
	itemTTL time.Duration

	mtx         sync.RWMutex
	serviceHost map[string]dynamicServiceItem // protected by mtx
}

type dynamicServiceItem struct {
	hosts   []string
	expired int64 // Unix timestamp
}

var _ ServiceMapper = &dynamicServiceMapper{}

func NewDynamicServiceMapper(serviceRegistryClient Client, ttl time.Duration) *dynamicServiceMapper {
	return &dynamicServiceMapper{
		occ:         serviceRegistryClient,
		itemTTL:     ttl,
		serviceHost: make(map[string]dynamicServiceItem),
	}
}

func (m *dynamicServiceMapper) GetHost(svc string) (string, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	var err error

	svcItem, ok := m.serviceHost[svc]
	if !ok {
		svcItem, err = m.resolveService(svc)
		if err != nil {
			return "", err
		}
	}

	if time.Now().Unix() > svcItem.expired {
		svcItem, err = m.resolveService(svc)
		if err != nil {
			return "", err
		}
	}

	i := rand.IntN(len(svcItem.hosts))
	return svcItem.hosts[i], nil
}

func (m *dynamicServiceMapper) resolveService(svc string) (dynamicServiceItem, error) {
	// m.mtx is assumed to be read-locked

	item := dynamicServiceItem{}

	availableSvc, _, err := GenericAutoDecodedDo[*serviceregistry.AvailableServicesResponse](
		m.occ,
		func() (*Request, error) {
			return NewRequest(
				context.Background(),
				http.MethodGet,
				serviceregistry.ServiceType,
				"/services/available.json",
				NoBody,
				WithQuery("serviceType", svc),
			)
		},
	)
	if err != nil {
		return item, err
	}

	if availableSvc.Services.Type == strobj.String {
		return item, ServiceNotFoundErr
	}

	svcObjectList := availableSvc.Services.ObjectVal
	switch svcObjectList.Service.Type {
	case objlist.Object:
		item.hosts = []string{svcObjectList.Service.ObjectVal.Host}

	case objlist.List:
		item.hosts = make([]string, 0, len(svcObjectList.Service.ListVal))
		for _, s := range svcObjectList.Service.ListVal {
			item.hosts = append(item.hosts, s.Host)
		}
	}

	// m.mtx is read-locked -> upgrade to write lock and exit function with read-lock again
	m.mtx.RUnlock()
	m.mtx.Lock()
	defer m.mtx.RLock()
	defer m.mtx.Unlock()

	item.expired = time.Now().Add(m.itemTTL).Unix()
	m.serviceHost[svc] = item

	return item, nil
}
