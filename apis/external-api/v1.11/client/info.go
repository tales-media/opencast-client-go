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

	extapiv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11"
	"shio.solutions/tales.media/opencast-client-go/apis/meta/base"
	oc "shio.solutions/tales.media/opencast-client-go/client"
)

func (c *client) GetInfoOrganization(ctx context.Context, opts ...oc.RequestOpts) (*extapiv1.Organization, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.Organization](
		c,
		func() (*oc.Request, error) { return c.GetInfoOrganizationRequest(ctx, opts...) },
	)
}

func (c *client) GetInfoOrganizationRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		ServiceType,
		"/api/info/organization",
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetInfoOrganizationProperties(ctx context.Context, opts ...oc.RequestOpts) (base.Properties, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[base.Properties](
		c,
		func() (*oc.Request, error) { return c.GetInfoOrganizationPropertiesRequest(ctx, opts...) },
	)
}

func (c *client) GetInfoOrganizationPropertiesRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		ServiceType,
		"/api/info/organization/properties",
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetInfoOrganizationPropertiesEngageUIURL(ctx context.Context, opts ...oc.RequestOpts) (base.Properties, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[base.Properties](
		c,
		func() (*oc.Request, error) { return c.GetInfoOrganizationPropertiesEngageUIURLRequest(ctx, opts...) },
	)
}

func (c *client) GetInfoOrganizationPropertiesEngageUIURLRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		ServiceType,
		"/api/info/organization/properties/engageuiurl",
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetInfoMe(ctx context.Context, opts ...oc.RequestOpts) (*extapiv1.Me, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.Me](
		c,
		func() (*oc.Request, error) { return c.GetInfoMeRequest(ctx, opts...) },
	)
}

func (c *client) GetInfoMeRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		ServiceType,
		"/api/info/me",
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetInfoMeRoles(ctx context.Context, opts ...oc.RequestOpts) ([]string, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[[]string](
		c,
		func() (*oc.Request, error) { return c.GetInfoMeRolesRequest(ctx, opts...) },
	)
}

func (c *client) GetInfoMeRolesRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		ServiceType,
		"/api/info/me/roles",
		oc.NoBody,
		opts...,
	)
}
