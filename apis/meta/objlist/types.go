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

package objlist

import (
	"encoding/json"
	"fmt"
)

type ObjectOrList[T any] struct {
	Type      Type
	ObjectVal T
	ListVal   []T
}

type Type int

const (
	Object Type = iota
	List
)

func FromObject[T any](val T) ObjectOrList[T] {
	return ObjectOrList[T]{
		Type:      Object,
		ObjectVal: val,
	}
}

func FromList[T any](val []T) ObjectOrList[T] {
	return ObjectOrList[T]{
		Type:    List,
		ListVal: val,
	}
}

func (objlist ObjectOrList[T]) MarshalJSON() ([]byte, error) {
	switch objlist.Type {
	case Object:
		return json.Marshal(objlist.ObjectVal)
	case List:
		return json.Marshal(objlist.ListVal)
	default:
		return []byte{}, fmt.Errorf("objlist: impossible ObjectOrList.Type")
	}
}

func (objlist *ObjectOrList[T]) UnmarshalJSON(value []byte) error {
	if value[0] == '[' {
		objlist.Type = List
		return json.Unmarshal(value, &objlist.ListVal)
	}
	objlist.Type = Object
	return json.Unmarshal(value, &objlist.ObjectVal)
}
