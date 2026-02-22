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
	"strings"

	extapiv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11"
	oc "shio.solutions/tales.media/opencast-client-go/client"
	"shio.solutions/tales.media/opencast-client-go/pkg/multipart"
)

type CreateGroupRequestBody struct {
	Name        string
	Description string
	Roles       []string
	Members     []string
}

type UpdateGroupRequestBody struct {
	Name        string
	Description string
	Roles       []string
	Members     []string
}

type CreateGroupMemberRequestBody struct {
	Member string
}

const (
	GroupNameFilterKey = FilterKey("name") // TODO: is this "Name" in Opencast?
)

const (
	GroupNameSortKey        = SortKey("name")
	GroupDescriptionSortKey = SortKey("description")
	GroupRoleSortKey        = SortKey("role")
)

func (c *client) ListGroup(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.Group, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[[]extapiv1.Group](
		c,
		func() (*oc.Request, error) { return c.ListGroupRequest(ctx, opts...) },
	)
}

func (c *client) ListGroupRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		GroupsServiceType,
		"/api/groups",
		oc.NoBody,
		opts...,
	)
}

func (c *client) CreateGroup(ctx context.Context, body *CreateGroupRequestBody, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.CreateGroupRequest(ctx, body, opts...) },
	)
}

func (c *client) CreateGroupRequest(ctx context.Context, body *CreateGroupRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	mp.AddPart(multipart.FormFieldString("name", body.Name))
	if body.Description != "" {
		mp.AddPart(multipart.FormFieldString("description", body.Description))
	}
	if len(body.Roles) > 0 {
		mp.AddPart(multipart.FormFieldString("roles", strings.Join(body.Roles, ",")))
	}
	if len(body.Members) > 0 {
		mp.AddPart(multipart.FormFieldString("members", strings.Join(body.Members, ",")))
	}
	return oc.NewRequest(
		ctx,
		http.MethodPost,
		GroupsServiceType,
		"/api/groups",
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) GetGroup(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.Group, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.Group](
		c,
		func() (*oc.Request, error) { return c.GetGroupRequest(ctx, id, opts...) },
	)
}

func (c *client) GetGroupRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		GroupsServiceType,
		"/api/groups/"+url.PathEscape(id),
		oc.NoBody,
		opts...,
	)
}

func (c *client) UpdateGroup(ctx context.Context, id string, body *UpdateGroupRequestBody, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.UpdateGroupRequest(ctx, id, body, opts...) },
	)
}

func (c *client) UpdateGroupRequest(ctx context.Context, id string, body *UpdateGroupRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	if body.Name != "" {
		mp.AddPart(multipart.FormFieldString("name", body.Name))
	}
	if body.Description != "" {
		mp.AddPart(multipart.FormFieldString("description", body.Description))
	}
	if len(body.Roles) > 0 {
		mp.AddPart(multipart.FormFieldString("roles", strings.Join(body.Roles, ",")))
	}
	if len(body.Members) > 0 {
		mp.AddPart(multipart.FormFieldString("members", strings.Join(body.Members, ",")))
	}
	return oc.NewRequest(
		ctx,
		http.MethodPost,
		GroupsServiceType,
		"/api/groups/"+url.PathEscape(id),
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) DeleteGroup(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.DeleteGroupRequest(ctx, id, opts...) },
	)
}

func (c *client) DeleteGroupRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodDelete,
		GroupsServiceType,
		"/api/groups/"+url.PathEscape(id),
		oc.NoBody,
		opts...,
	)
}

func (c *client) CreateGroupMember(ctx context.Context, id string, body *CreateGroupMemberRequestBody, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.CreateGroupMemberRequest(ctx, id, body, opts...) },
	)
}

func (c *client) CreateGroupMemberRequest(ctx context.Context, id string, body *CreateGroupMemberRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	mp.AddPart(multipart.FormFieldString("member", body.Member))
	return oc.NewRequest(
		ctx,
		http.MethodPost,
		GroupsServiceType,
		"/api/groups/"+url.PathEscape(id)+"/members",
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) DeleteGroupMember(ctx context.Context, id, memberID string, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.DeleteGroupMemberRequest(ctx, id, memberID, opts...) },
	)
}

func (c *client) DeleteGroupMemberRequest(ctx context.Context, id, memberID string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodDelete,
		GroupsServiceType,
		"/api/groups/"+url.PathEscape(id)+"/members/"+url.PathEscape(memberID),
		oc.NoBody,
		opts...,
	)
}
