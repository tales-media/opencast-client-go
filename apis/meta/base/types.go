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

package base

import (
	"encoding/json"
	"errors"
	"time"
)

const HourMinuteOnly = "15:04"

type Properties map[string]string

type DateTime struct {
	Time time.Time

	L string
}

func (dt DateTime) IsZero() bool {
	return dt.Time.IsZero()
}

func (dt DateTime) Layout() string {
	if dt.L == "" {
		return time.RFC3339
	}
	return dt.L
}

func (dt DateTime) MarshalText() ([]byte, error) {
	if dt.IsZero() {
		return []byte{}, nil
	}

	// 10 as used by time.Format
	layout := dt.Layout()
	bufLen := len(layout) + 10

	buf := make([]byte, 0, bufLen)
	buf = dt.Time.AppendFormat(buf, layout)

	return buf, nil
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	if dt.IsZero() {
		return []byte{'"', '"'}, nil
	}

	//  2 for ""
	// 10 as used by time.Format
	layout := dt.Layout()
	bufLen := len(layout) + 2 + 10

	buf := make([]byte, 0, bufLen)
	buf = append(buf, '"')
	buf = dt.Time.AppendFormat(buf, layout)
	buf = append(buf, '"')

	return buf, nil
}

func (dt *DateTime) UnmarshalText(data []byte) error {
	var (
		str = string(data)
		err error
	)

	if str == "" {
		// Opencast represents null dates as blank string. Skip unmarshal and let dt remain as zero value.
		return nil
	}

	dt.L, err = detectDateLayout(str)
	if err != nil {
		return err
	}

	dt.Time, err = time.Parse(dt.L, str)
	if err != nil {
		return err
	}

	return nil
}

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	if data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("DateTime.UnmarshalJSON: input is not a JSON string")
	}
	data = data[len(`"`) : len(data)-len(`"`)]
	return dt.UnmarshalText(data)
}

func detectDateLayout(date string) (string, error) {
	l := len(date)
	switch {
	case l > len(time.DateOnly):
		// 2006-01-02T15:04:05Z07:00
		return time.RFC3339, nil

	case l > len(time.TimeOnly):
		// 2006-01-02
		return time.DateOnly, nil

	case l > len(HourMinuteOnly):
		// 15:04:05
		return time.TimeOnly, nil

	case l == len(HourMinuteOnly):
		// 15:04
		return HourMinuteOnly, nil

	default:
		return "", errors.New("base: unknown date or time format")
	}
}

type Int int64

func (i Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(i))
}

func (i *Int) UnmarshalJSON(data []byte) error {
	if len(data) == 2 && data[0] == '"' && data[1] == '"' {
		// Opencast represents null ints as blank string. Skip unmarshal and let i remain as zero value.
		return nil
	}
	var i2 int64
	if err := json.Unmarshal(data, &i2); err != nil {
		return err
	}
	*i = Int(i2)
	return nil
}

type Float float64

func (f Float) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(f))
}

func (f *Float) UnmarshalJSON(data []byte) error {
	if len(data) == 2 && data[0] == '"' && data[1] == '"' {
		// Opencast represents null floats as blank string. Skip unmarshal and let f remain as zero value.
		return nil
	}
	var f2 float64
	if err := json.Unmarshal(data, &f2); err != nil {
		return err
	}
	*f = Float(f2)
	return nil
}

type Flavor string

const (
	DublinCoreEpisodeFlavor = Flavor("dublincore/episode")
	DublinCoreSeriesFlavor  = Flavor("dublincore/series")
	SecurityEpisodeFlavor   = Flavor("security/xacml+episode")
	SecuritySeriesFlavor    = Flavor("security/xacml+series")
	SMILCuttingFlavor       = Flavor("smil/cutting")
	MPEG7SegmentsFlavor     = Flavor("mpeg-7/segments")
)

type Action string

const (
	ReadAction  = Action("read")
	WriteAction = Action("write")
)
