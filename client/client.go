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
	"net/http"
	"time"
)

const (
	Version   = "1.0"
	UserAgent = "OpencastGoClient/" + Version
)

type Doer interface {
	Do(*Request) (*Response, error)
}

type Client interface {
	Doer
}

type client struct {
	sm      ServiceMapper
	http    http.Client
	reqOpts []RequestOpts
}

var _ Client = &client{}

func New(sm ServiceMapper, opts ...ClientOpts) (Client, error) {
	c := &client{
		sm: sm,
		reqOpts: []RequestOpts{
			WithHeader("User-Agent", UserAgent),
		},
	}
	if err := c.ApplyOptions(opts...); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *client) ApplyOptions(opts ...ClientOpts) error {
	for _, opt := range opts {
		if err := opt.Apply(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) Do(req *Request) (*Response, error) {
	if err := req.ApplyOptions(c.reqOpts...); err != nil {
		return nil, err
	}

	httpReq, err := req.HTTPRequest(c.sm)
	if err != nil {
		return nil, err
	}

	reqStart := time.Now()
	httpResp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}

	resp := newResponse(httpResp)
	resp.Meta.Duration = time.Since(reqStart)

	return resp, nil
}

type ClientOpts interface {
	Apply(*client) error
}

type ClientOptsFunc func(*client) error

func (f ClientOptsFunc) Apply(c *client) error { return f(c) }

func WithHTTPClient(h http.Client) ClientOpts {
	return ClientOptsFunc(func(c *client) error {
		c.http = h
		return nil
	})
}

func WithRequestOptions(opts ...RequestOpts) ClientOpts {
	return ClientOptsFunc(func(c *client) error {
		c.reqOpts = opts
		return nil
	})
}

// TODO: add WithBackoff
