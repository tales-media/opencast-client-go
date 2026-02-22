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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	http.Response

	Meta ResponseMeta
}

type ResponseMeta struct {
	Duration time.Duration
}

func newResponse(httpResp *http.Response) *Response {
	resp := &Response{Response: *httpResp}
	return resp
}

func (resp *Response) Decode(v any, f DecoderReaderFunc) error {
	return f(v, resp)
}

type DecoderReaderFunc func(any, *Response) error

func AutoDecoder(v any, resp *Response) error {
	ct := resp.Header.Get("Content-Type")
	mt, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return err
	}

	switch {
	// JSON
	case strings.HasPrefix(mt, "application/") && strings.HasSuffix(mt, "+json"):
		// application/health+json
		// application/{api-version}+json
		fallthrough
	case mt == "application/json":
		return JsonDecoder(v, resp)

	// XML
	case strings.HasPrefix(mt, "application/") && strings.HasSuffix(mt, "+xml"):
		fallthrough
	case mt == "application/xml" || mt == "text/xml":
		return XMLDecoder(v, resp)

	default:
		// Other media types used in Opencast
		//   application/octet-stream
		//   text/calendar
		//   text/html
		//   text/plain
		//   */* (serving files)
		return fmt.Errorf("AutoDecoder: unsupported media type %s", mt)
	}
}

func JsonDecoder(v any, resp *Response) error {
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(v)
}

func XMLDecoder(v any, resp *Response) error {
	dec := xml.NewDecoder(resp.Body)
	return dec.Decode(v)
}
