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
	"net/http"
	"net/url"

	"shio.solutions/tales.media/opencast-client-go/apis/meta/base"
	oc "shio.solutions/tales.media/opencast-client-go/client"
)

func (c *client) ListListProvider(ctx context.Context, opts ...oc.RequestOpts) ([]string, *oc.Response, error) {
	l, resp, err := oc.GenericAutoDecodedDo[[][]string](
		c,
		func() (*oc.Request, error) { return c.ListListProviderRequest(ctx, opts...) },
	)
	if err != nil {
		return nil, resp, err
	}
	if len(l) == 0 {
		return []string{}, resp, err
	}
	return l[0], resp, err
}

func (c *client) ListListProviderRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		ListProvidersServiceType,
		"/api/listproviders/providers.json",
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetListProvider(ctx context.Context, source string, opts ...oc.RequestOpts) (base.Properties, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[base.Properties](
		c,
		func() (*oc.Request, error) { return c.GetListProviderRequest(ctx, source, opts...) },
	)
}

func (c *client) GetListProviderRequest(ctx context.Context, source string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		ListProvidersServiceType,
		"/api/listproviders/"+url.PathEscape(source)+".json",
		oc.NoBody,
		opts...,
	)
}
