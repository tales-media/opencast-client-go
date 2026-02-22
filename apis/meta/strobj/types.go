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

package strobj

import (
	"encoding/json"
	"fmt"
)

type StringOrObject[T any] struct {
	Type      Type
	StringVal string
	ObjectVal T
}

type Type int

const (
	String Type = iota
	Object
)

func FromString[T any](val string) StringOrObject[T] {
	return StringOrObject[T]{
		Type:      String,
		StringVal: val,
	}
}

func FromObject[T any](val T) StringOrObject[T] {
	return StringOrObject[T]{
		Type:      Object,
		ObjectVal: val,
	}
}

func (strobj StringOrObject[T]) MarshalJSON() ([]byte, error) {
	switch strobj.Type {
	case String:
		return json.Marshal(strobj.StringVal)
	case Object:
		return json.Marshal(strobj.ObjectVal)
	default:
		return []byte{}, fmt.Errorf("strobj: impossible StringOrObject.Type")
	}
}

func (strobj *StringOrObject[T]) UnmarshalJSON(value []byte) error {
	if value[0] == '"' {
		strobj.Type = String
		return json.Unmarshal(value, &strobj.StringVal)
	}
	strobj.Type = Object
	return json.Unmarshal(value, &strobj.ObjectVal)
}
