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
	"encoding/base64"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const (
	// Header name for running an operation as a desired user.
	RunAsUserHeader = "X-RUN-AS-USER"

	// Header name for running an operation with a desired set of role.
	RunWithRolesHeader = "X-RUN-WITH-ROLES"
)

type Request struct {
	Ctx     context.Context
	Method  string
	Service string
	Path    string
	Query   url.Values
	Header  http.Header
	Body    Body
}

func NewRequest(ctx context.Context, method, service, path string, body Body, opts ...RequestOpts) (*Request, error) {
	req := &Request{
		Ctx:     ctx,
		Method:  method,
		Service: service,
		Path:    path,
		Query:   make(url.Values),
		Header:  make(http.Header),
		Body:    body,
	}
	if err := req.ApplyOptions(opts...); err != nil {
		return nil, err
	}
	return req, nil
}

func (req *Request) ApplyOptions(opts ...RequestOpts) error {
	for _, opt := range opts {
		if err := opt.Apply(req); err != nil {
			return err
		}
	}
	return nil
}

func (req *Request) URL(sm ServiceMapper) (*url.URL, error) {
	hostURL, err := sm.GetHost(req.Service)
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(hostURL)
	if err != nil {
		return nil, err
	}

	url.Path = path.Join(url.Path, req.Path)
	url.RawQuery = req.Query.Encode()

	return url, nil
}

func (req *Request) HTTPRequest(sm ServiceMapper) (*http.Request, error) {
	url, err := req.URL(sm)
	if err != nil {
		return nil, err
	}

	if req.Body == nil {
		req.Body = NoBody
	}

	body, err := req.Body.Reader()
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(req.Ctx, req.Method, url.String(), body)
	if err != nil {
		return nil, err
	}

	httpReq.GetBody = req.Body.Reader
	httpReq.ContentLength = req.Body.Len()
	httpReq.Header = req.Header
	if req.Body.ContentType() != "" {
		httpReq.Header.Set("Content-Type", req.Body.ContentType())
	}

	return httpReq, nil
}

type RequestOpts interface {
	Apply(*Request) error
}

type RequestOptsFunc func(*Request) error

func (f RequestOptsFunc) Apply(r *Request) error { return f(r) }

func WithQuery(key, value string) RequestOpts {
	return RequestOptsFunc(func(req *Request) error {
		req.Query.Set(key, value)
		return nil
	})
}

func WithoutQuery(key string) RequestOpts {
	return RequestOptsFunc(func(req *Request) error {
		req.Query.Del(key)
		return nil
	})
}

func WithHeader(key, value string) RequestOpts {
	return RequestOptsFunc(func(req *Request) error {
		req.Header.Set(key, value)
		return nil
	})
}

func WithoutHeader(key string) RequestOpts {
	return RequestOptsFunc(func(req *Request) error {
		req.Header.Del(key)
		return nil
	})
}

func WithBasicAuth(username, password string) RequestOpts {
	auth := username + ":" + password
	basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	return WithHeader("Authorization", "Basic "+basicAuth)
}

func WithJWTHeader(header, prefix, token string) RequestOpts {
	return WithHeader(header, prefix+token)
}

func WithJWTQuery(token string) RequestOpts {
	return WithQuery("jwt", token)
}

func WithRunAsUser(username string) RequestOpts {
	return WithHeader(RunAsUserHeader, username)
}

func WithRunWithRoles(roles ...string) RequestOpts {
	return WithHeader(RunWithRolesHeader, strings.Join(roles, ","))
}
