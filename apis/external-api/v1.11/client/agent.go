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

	extapiv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11"
	oc "shio.solutions/tales.media/opencast-client-go/client"
)

func (c *client) ListAgent(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.Agent, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[[]extapiv1.Agent](
		c,
		func() (*oc.Request, error) { return c.ListAgentRequest(ctx, opts...) },
	)
}

func (c *client) ListAgentRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		AgentsServiceType,
		"/api/agents",
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetAgent(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.Agent, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.Agent](
		c,
		func() (*oc.Request, error) { return c.GetAgentRequest(ctx, id, opts...) },
	)
}

func (c *client) GetAgentRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		AgentsServiceType,
		"/api/agents/"+url.PathEscape(id),
		oc.NoBody,
		opts...,
	)
}
