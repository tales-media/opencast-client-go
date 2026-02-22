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

package serviceregistry

import (
	"time"

	"shio.solutions/tales.media/opencast-client-go/apis/meta/objlist"
	"shio.solutions/tales.media/opencast-client-go/apis/meta/strobj"
)

const ServiceType = "org.opencastproject.serviceregistry"

type AvailableServicesResponse struct {
	Services strobj.StringOrObject[AvailableServicesList] `json:"services"`
}

type AvailableServicesList struct {
	Service objlist.ObjectOrList[Service] `json:"service"`
}

type Service struct {
	Type                string       `json:"type"`
	Host                string       `json:"host"`
	Path                string       `json:"path"`
	Active              bool         `json:"active"`
	Online              bool         `json:"online"`
	Maintenance         bool         `json:"maintenance"`
	JobProducer         bool         `json:"jobproducer"`
	OnlineFrom          time.Time    `json:"onlinefrom"`
	ServiceState        ServiceState `json:"service_state"`
	StateChanged        time.Time    `json:"state_changed"`
	ErrorStateTrigger   int          `json:"error_state_trigger"`
	WarningStateTrigger int          `json:"warning_state_trigger"`
}

type ServiceState string

const (
	NormalServiceState  = ServiceState("NORMAL")
	WarningServiceState = ServiceState("WARNING")
	ErrorServiceState   = ServiceState("ERROR")
)
