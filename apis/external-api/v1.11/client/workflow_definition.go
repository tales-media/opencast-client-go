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
	"strconv"

	extapiv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11"
	oc "shio.solutions/tales.media/opencast-client-go/client"
)

type WithWorkflowDefinitionOptions struct {
	WithOperations             bool
	WithConfigurationPanel     bool
	WithConfigurationPanelJSON bool
}

var _ oc.RequestOpts = WithWorkflowDefinitionOptions{}

func (opt WithWorkflowDefinitionOptions) Apply(r *oc.Request) error {
	return r.ApplyOptions(
		oc.WithQuery("withoperations", strconv.FormatBool(opt.WithOperations)),
		oc.WithQuery("withconfigurationpanel", strconv.FormatBool(opt.WithConfigurationPanel)),
		oc.WithQuery("withconfigurationpaneljson", strconv.FormatBool(opt.WithConfigurationPanelJSON)),
	)
}

const (
	WorkflowDefinitionTagFilterKey = FilterKey("tag")
)

const (
	WorkflowDefinitionIdentifierSortKey   = SortKey("identifier")
	WorkflowDefinitionTitleSortKey        = SortKey("title")
	WorkflowDefinitionDisplayOrderSortKey = SortKey("displayorder")
)

func (c *client) ListWorkflowDefinition(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.WorkflowDefinition, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[[]extapiv1.WorkflowDefinition](
		c,
		func() (*oc.Request, error) { return c.ListWorkflowDefinitionRequest(ctx, opts...) },
	)
}

func (c *client) ListWorkflowDefinitionRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		WorkflowDefinitionsServiceType,
		"/api/workflow-definitions",
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetWorkflowDefinition(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.WorkflowDefinition, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.WorkflowDefinition](
		c,
		func() (*oc.Request, error) { return c.GetWorkflowDefinitionRequest(ctx, id, opts...) },
	)
}

func (c *client) GetWorkflowDefinitionRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		WorkflowDefinitionsServiceType,
		"/api/workflow-definitions/"+url.PathEscape(id),
		oc.NoBody,
		opts...,
	)
}
