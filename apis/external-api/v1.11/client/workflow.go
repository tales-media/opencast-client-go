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
	"shio.solutions/tales.media/opencast-client-go/apis/meta/base"
	oc "shio.solutions/tales.media/opencast-client-go/client"
	"shio.solutions/tales.media/opencast-client-go/pkg/multipart"
)

type CreateWorkflowRequestBody struct {
	EventID              string
	WorkflowDefinitionID string
	Configuration        base.Properties
}

type UpdateWorkflowRequestBody struct {
	State         *extapiv1.WorkflowState
	Configuration base.Properties
}

type WithWorkflowOptions struct {
	WithOperations    bool
	WithConfiguration bool
}

var _ oc.RequestOpts = WithWorkflowOptions{}

func (opt WithWorkflowOptions) Apply(r *oc.Request) error {
	return r.ApplyOptions(
		oc.WithQuery("withoperations", strconv.FormatBool(opt.WithOperations)),
		oc.WithQuery("withconfiguration", strconv.FormatBool(opt.WithConfiguration)),
	)
}

func (c *client) CreateWorkflow(ctx context.Context, body *CreateWorkflowRequestBody, opts ...oc.RequestOpts) (*extapiv1.WorkflowInstance, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.WorkflowInstance](
		c,
		func() (*oc.Request, error) {
			return c.CreateWorkflowRequest(ctx, body, opts...)
		},
	)
}

func (c *client) CreateWorkflowRequest(ctx context.Context, body *CreateWorkflowRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	mp.AddParts(
		multipart.FormFieldString("event_identifier", body.EventID),
		multipart.FormFieldString("workflow_definition_identifier", body.WorkflowDefinitionID),
	)
	if len(body.Configuration) > 0 {
		configuration, err := json.Marshal(body.Configuration)
		if err != nil {
			return nil, err
		}
		mp.AddPart(multipart.FormField("configuration", configuration))
	}
	return oc.NewRequest(
		ctx,
		http.MethodPost,
		WorkflowInstancesServiceType,
		"/api/workflows",
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) GetWorkflow(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.WorkflowInstance, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.WorkflowInstance](
		c,
		func() (*oc.Request, error) { return c.GetWorkflowRequest(ctx, id, opts...) },
	)
}

func (c *client) GetWorkflowRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		WorkflowInstancesServiceType,
		"/api/workflows/"+url.PathEscape(id),
		oc.NoBody,
		opts...,
	)
}

func (c *client) UpdateWorkflow(ctx context.Context, id string, body *UpdateWorkflowRequestBody, opts ...oc.RequestOpts) (*extapiv1.WorkflowInstance, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.WorkflowInstance](
		c,
		func() (*oc.Request, error) { return c.UpdateWorkflowRequest(ctx, id, body, opts...) },
	)
}

func (c *client) UpdateWorkflowRequest(ctx context.Context, id string, body *UpdateWorkflowRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	if body.State != nil {
		mp.AddPart(multipart.FormFieldString("state", string(*body.State)))
	}
	if len(body.Configuration) > 0 {
		configuration, err := json.Marshal(body.Configuration)
		if err != nil {
			return nil, err
		}
		mp.AddPart(multipart.FormField("configuration", configuration))
	}
	return oc.NewRequest(
		ctx,
		http.MethodPut,
		WorkflowInstancesServiceType,
		"/api/workflows/"+url.PathEscape(id),
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) DeleteWorkflow(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.DeleteWorkflowRequest(ctx, id, opts...) },
	)
}

func (c *client) DeleteWorkflowRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodDelete,
		WorkflowInstancesServiceType,
		"/api/workflows/"+url.PathEscape(id),
		oc.NoBody,
		opts...,
	)
}
