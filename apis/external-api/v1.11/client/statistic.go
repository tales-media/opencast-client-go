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
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	extapiv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11"
	oc "shio.solutions/tales.media/opencast-client-go/client"
	"shio.solutions/tales.media/opencast-client-go/pkg/multipart"
)

type QueryStatisticRequestBody struct {
	Query extapiv1.StatisticQuery
}

type ExportCSVStatisticRequestBody struct {
	Query extapiv1.StatisticQuery
}

type WithStatisticOptions struct {
	WithParameters bool
}

var _ oc.RequestOpts = WithStatisticOptions{}

func (opt WithStatisticOptions) Apply(r *oc.Request) error {
	return r.ApplyOptions(
		oc.WithQuery("withparameters", strconv.FormatBool(opt.WithParameters)),
	)
}

const (
	StatisticResourceTypeFilterKey = FilterKey("resourceType")
)

func (c *client) ListStatisticProvider(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.StatisticProvider, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[[]extapiv1.StatisticProvider](
		c,
		func() (*oc.Request, error) { return c.ListStatisticProviderRequest(ctx, opts...) },
	)
}

func (c *client) ListStatisticProviderRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		StatisticsServiceType,
		"/api/statistics/providers",
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetStatisticProvider(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.StatisticProvider, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.StatisticProvider](
		c,
		func() (*oc.Request, error) { return c.GetStatisticProviderRequest(ctx, id, opts...) },
	)
}

func (c *client) GetStatisticProviderRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		StatisticsServiceType,
		"/api/statistics/providers/"+url.PathEscape(id),
		oc.NoBody,
		opts...,
	)
}

func (c *client) QueryStatistic(ctx context.Context, body *QueryStatisticRequestBody, opts ...oc.RequestOpts) ([]extapiv1.StatisticQueryResult, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[[]extapiv1.StatisticQueryResult](
		c,
		func() (*oc.Request, error) { return c.QueryStatisticRequest(ctx, body, opts...) },
	)
}

func (c *client) QueryStatisticRequest(ctx context.Context, body *QueryStatisticRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	data, err := json.Marshal(body.Query)
	if err != nil {
		return nil, err
	}
	mp.AddPart(multipart.FormField("data", data))
	return oc.NewRequest(
		ctx,
		http.MethodPost,
		StatisticsServiceType,
		"/api/statistics/data/query",
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) ExportCSVStatistic(ctx context.Context, body *ExportCSVStatisticRequestBody, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.ExportCSVStatisticRequest(ctx, body, opts...) },
	)
}

func (c *client) ExportCSVStatisticRequest(ctx context.Context, body *ExportCSVStatisticRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	data, err := json.Marshal(body.Query)
	if err != nil {
		return nil, err
	}
	mp.AddPart(multipart.FormField("data", data))
	return oc.NewRequest(
		ctx,
		http.MethodPost,
		StatisticsServiceType,
		"/api/statistics/data/export.csv",
		oc.NewMultipartBody(mp),
		opts...,
	)
}
