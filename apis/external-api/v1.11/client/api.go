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
	"shio.solutions/tales.media/opencast-client-go/apis/meta/objlist"
	oc "shio.solutions/tales.media/opencast-client-go/client"
)

const (
	AcceptJSONHeader = "application/" + extapiv1.Version + "+json"
)

const (
	ServiceType                    = "org.opencastproject.external"
	AgentsServiceType              = "org.opencastproject.external.agents"
	EventsServiceType              = "org.opencastproject.external.events"
	GroupsServiceType              = "org.opencastproject.external.groups"
	ListProvidersServiceType       = "org.opencastproject.external.listproviders"
	PlaylistsServiceType           = "org.opencastproject.external.playlists"
	SecurityServiceType            = "org.opencastproject.external.security"
	SeriesServiceType              = "org.opencastproject.external" // TODO: fix this in Opencast
	StatisticsServiceType          = "org.opencastproject.external.statistics"
	WorkflowDefinitionsServiceType = "org.opencastproject.external.workflows.definitions"
	WorkflowInstancesServiceType   = "org.opencastproject.external.workflows.instances"
)

type Client interface {
	Do(*oc.Request) (*oc.Response, error)
	OpencastClient() oc.Client

	// API

	GetAPI(ctx context.Context, opts ...oc.RequestOpts) (*extapiv1.API, *oc.Response, error)
	GetAPIRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	GetAPIVersion(ctx context.Context, opts ...oc.RequestOpts) (*extapiv1.APIVersion, *oc.Response, error)
	GetAPIVersionRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	GetAPIVersionDefault(ctx context.Context, opts ...oc.RequestOpts) (*extapiv1.APIVersion, *oc.Response, error)
	GetAPIVersionDefaultRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	// Info

	GetInfoOrganization(ctx context.Context, opts ...oc.RequestOpts) (*extapiv1.Organization, *oc.Response, error)
	GetInfoOrganizationRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	GetInfoOrganizationProperties(ctx context.Context, opts ...oc.RequestOpts) (base.Properties, *oc.Response, error)
	GetInfoOrganizationPropertiesRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	GetInfoOrganizationPropertiesEngageUIURL(ctx context.Context, opts ...oc.RequestOpts) (base.Properties, *oc.Response, error)
	GetInfoOrganizationPropertiesEngageUIURLRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	GetInfoMe(ctx context.Context, opts ...oc.RequestOpts) (*extapiv1.Me, *oc.Response, error)
	GetInfoMeRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	GetInfoMeRoles(ctx context.Context, opts ...oc.RequestOpts) ([]string, *oc.Response, error)
	GetInfoMeRolesRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	// Security

	SignURL(ctx context.Context, body *SignURLRequestBody, opts ...oc.RequestOpts) (*extapiv1.SignedURL, *oc.Response, error)
	SignURLRequest(ctx context.Context, body *SignURLRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	// List Providers

	ListListProvider(ctx context.Context, opts ...oc.RequestOpts) ([]string, *oc.Response, error)
	ListListProviderRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	GetListProvider(ctx context.Context, source string, opts ...oc.RequestOpts) (base.Properties, *oc.Response, error)
	GetListProviderRequest(ctx context.Context, source string, opts ...oc.RequestOpts) (*oc.Request, error)

	// Groups

	ListGroup(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.Group, *oc.Response, error)
	ListGroupRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	CreateGroup(ctx context.Context, body *CreateGroupRequestBody, opts ...oc.RequestOpts) (*oc.Response, error)
	CreateGroupRequest(ctx context.Context, body *CreateGroupRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	GetGroup(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.Group, *oc.Response, error)
	GetGroupRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	UpdateGroup(ctx context.Context, id string, body *UpdateGroupRequestBody, opts ...oc.RequestOpts) (*oc.Response, error)
	UpdateGroupRequest(ctx context.Context, id string, body *UpdateGroupRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	DeleteGroup(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Response, error)
	DeleteGroupRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	// Groups - Member

	CreateGroupMember(ctx context.Context, id string, body *CreateGroupMemberRequestBody, opts ...oc.RequestOpts) (*oc.Response, error)
	CreateGroupMemberRequest(ctx context.Context, id string, body *CreateGroupMemberRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	DeleteGroupMember(ctx context.Context, id, memberID string, opts ...oc.RequestOpts) (*oc.Response, error)
	DeleteGroupMemberRequest(ctx context.Context, id, memberID string, opts ...oc.RequestOpts) (*oc.Request, error)

	// Statistics

	ListStatisticProvider(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.StatisticProvider, *oc.Response, error)
	ListStatisticProviderRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	GetStatisticProvider(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.StatisticProvider, *oc.Response, error)
	GetStatisticProviderRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	QueryStatistic(ctx context.Context, body *QueryStatisticRequestBody, opts ...oc.RequestOpts) ([]extapiv1.StatisticQueryResult, *oc.Response, error)
	QueryStatisticRequest(ctx context.Context, body *QueryStatisticRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	ExportCSVStatistic(ctx context.Context, body *ExportCSVStatisticRequestBody, opts ...oc.RequestOpts) (*oc.Response, error)
	ExportCSVStatisticRequest(ctx context.Context, body *ExportCSVStatisticRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	// Agents

	ListAgent(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.Agent, *oc.Response, error)
	ListAgentRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	GetAgent(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.Agent, *oc.Response, error)
	GetAgentRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	// Events

	ListEvent(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.Event, *oc.Response, error)
	ListEventRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	CreateEvent(ctx context.Context, body *CreateEventRequestBody, opts ...oc.RequestOpts) (objlist.ObjectOrList[extapiv1.Identifier], *oc.Response, error)
	CreateEventRequest(ctx context.Context, body *CreateEventRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	GetEvent(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.Event, *oc.Response, error)
	GetEventRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	UpdateEvent(ctx context.Context, id string, body *UpdateEventRequestBody, opts ...oc.RequestOpts) (*oc.Response, error)
	UpdateEventRequest(ctx context.Context, id string, body *UpdateEventRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	DeleteEvent(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Response, error)
	DeleteEventRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	// Events - Access Policy

	GetEventACL(ctx context.Context, id string, opts ...oc.RequestOpts) (extapiv1.ACL, *oc.Response, error)
	GetEventACLRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	UpdateEventACL(ctx context.Context, id string, body *UpdateEventACLRequestBody, opts ...oc.RequestOpts) (*oc.Response, error)
	UpdateEventACLRequest(ctx context.Context, id string, body *UpdateEventACLRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	CreateEventACE(ctx context.Context, id string, action base.Action, role string, opts ...oc.RequestOpts) (*oc.Response, error)
	CreateEventACERequest(ctx context.Context, id string, action base.Action, role string, opts ...oc.RequestOpts) (*oc.Request, error)

	DeleteEventACE(ctx context.Context, id string, action base.Action, role string, opts ...oc.RequestOpts) (*oc.Response, error)
	DeleteEventACERequest(ctx context.Context, id string, action base.Action, role string, opts ...oc.RequestOpts) (*oc.Request, error)

	// Events - Media

	ListEventMedia(ctx context.Context, id string, opts ...oc.RequestOpts) ([]extapiv1.MediaTrackElement, *oc.Response, error)
	ListEventMediaRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	CreateEventTrack(ctx context.Context, id string, body *CreateEventTrackRequestBody, opts ...oc.RequestOpts) (*oc.Response, error)
	CreateEventTrackRequest(ctx context.Context, id string, body *CreateEventTrackRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	// Events - Metadata

	ListEventMetadata(ctx context.Context, id string, opts ...oc.RequestOpts) ([]extapiv1.Catalog, *oc.Response, error)
	ListEventMetadataRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	GetEventMetadata(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) ([]extapiv1.Field, *oc.Response, error)
	GetEventMetadataRequest(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) (*oc.Request, error)

	UpdateEventMetadata(ctx context.Context, id string, flavor base.Flavor, body *UpdateEventMetadataRequestBody, opts ...oc.RequestOpts) (*oc.Response, error)
	UpdateEventMetadataRequest(ctx context.Context, id string, flavor base.Flavor, body *UpdateEventMetadataRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	DeleteEventMetadata(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) (*oc.Response, error)
	DeleteEventMetadataRequest(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) (*oc.Request, error)

	// Events - Publications

	ListEventPublication(ctx context.Context, id string, opts ...oc.RequestOpts) ([]extapiv1.Publication, *oc.Response, error)
	ListEventPublicationRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	GetEventPublication(ctx context.Context, id string, publicationID string, opts ...oc.RequestOpts) (*extapiv1.Publication, *oc.Response, error)
	GetEventPublicationRequest(ctx context.Context, id string, publicationID string, opts ...oc.RequestOpts) (*oc.Request, error)

	// Events - Scheduling

	GetEventScheduling(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.Scheduling, *oc.Response, error)
	GetEventSchedulingRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	UpdateEventScheduling(ctx context.Context, id string, body *UpdateEventSchedulingRequestBody, opts ...oc.RequestOpts) (*oc.Response, error)
	UpdateEventSchedulingRequest(ctx context.Context, id string, body *UpdateEventSchedulingRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	// Series

	ListSeries(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.Series, *oc.Response, error)
	ListSeriesRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	SearchSeries(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.Series, *oc.Response, error)
	SearchSeriesRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	CreateSeries(ctx context.Context, body *CreateSeriesRequestBody, opts ...oc.RequestOpts) (extapiv1.Identifier, *oc.Response, error)
	CreateSeriesRequest(ctx context.Context, body *CreateSeriesRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	GetSeries(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.Series, *oc.Response, error)
	GetSeriesRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	UpdateSeries(ctx context.Context, id string, body *UpdateSeriesRequestBody, opts ...oc.RequestOpts) (*oc.Response, error)
	UpdateSeriesRequest(ctx context.Context, id string, body *UpdateSeriesRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	DeleteSeries(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Response, error)
	DeleteSeriesRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	// Series - Access Policy

	GetSeriesACL(ctx context.Context, id string, opts ...oc.RequestOpts) (extapiv1.ACL, *oc.Response, error)
	GetSeriesACLRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	UpdateSeriesACL(ctx context.Context, id string, body *UpdateSeriesACLRequestBody, opts ...oc.RequestOpts) (*oc.Response, error)
	UpdateSeriesACLRequest(ctx context.Context, id string, body *UpdateSeriesACLRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	// Series - Metadata

	ListSeriesMetadata(ctx context.Context, id string, opts ...oc.RequestOpts) ([]extapiv1.Catalog, *oc.Response, error)
	ListSeriesMetadataRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	GetSeriesMetadata(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) (*extapiv1.Catalog, *oc.Response, error)
	GetSeriesMetadataRequest(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) (*oc.Request, error)

	UpdateSeriesMetadata(ctx context.Context, id string, flavor base.Flavor, body *UpdateSeriesMetadataRequestBody, opts ...oc.RequestOpts) (*oc.Response, error)
	UpdateSeriesMetadataRequest(ctx context.Context, id string, flavor base.Flavor, body *UpdateSeriesMetadataRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	DeleteSeriesMetadata(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) (*oc.Response, error)
	DeleteSeriesMetadataRequest(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) (*oc.Request, error)

	// Series - Properties

	GetSeriesProperties(ctx context.Context, id string, opts ...oc.RequestOpts) (base.Properties, *oc.Response, error)
	GetSeriesPropertiesRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	UpdateSeriesProperties(ctx context.Context, id string, body *UpdateSeriesPropertiesRequestBody, opts ...oc.RequestOpts) (base.Properties, *oc.Response, error)
	UpdateSeriesPropertiesRequest(ctx context.Context, id string, body *UpdateSeriesPropertiesRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	// Playlists

	ListPlaylist(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.Playlist, *oc.Response, error)
	ListPlaylistRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	CreatePlaylist(ctx context.Context, body *CreatePlaylistRequestBody, opts ...oc.RequestOpts) (*extapiv1.Playlist, *oc.Response, error)
	CreatePlaylistRequest(ctx context.Context, body *CreatePlaylistRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	GetPlaylist(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.Playlist, *oc.Response, error)
	GetPlaylistRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	UpdatePlaylist(ctx context.Context, id string, body *UpdatePlaylistRequestBody, opts ...oc.RequestOpts) (*extapiv1.Playlist, *oc.Response, error)
	UpdatePlaylistRequest(ctx context.Context, id string, body *UpdatePlaylistRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	DeletePlaylist(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.Playlist, *oc.Response, error)
	DeletePlaylistRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	// Workflows

	// TODO: add list workflow to External API
	// ListWorkflow(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.Workflow, *oc.Response, error)
	// ListWorkflowRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	CreateWorkflow(ctx context.Context, body *CreateWorkflowRequestBody, opts ...oc.RequestOpts) (*extapiv1.WorkflowInstance, *oc.Response, error)
	CreateWorkflowRequest(ctx context.Context, body *CreateWorkflowRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	GetWorkflow(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.WorkflowInstance, *oc.Response, error)
	GetWorkflowRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	UpdateWorkflow(ctx context.Context, id string, body *UpdateWorkflowRequestBody, opts ...oc.RequestOpts) (*extapiv1.WorkflowInstance, *oc.Response, error)
	UpdateWorkflowRequest(ctx context.Context, id string, body *UpdateWorkflowRequestBody, opts ...oc.RequestOpts) (*oc.Request, error)

	DeleteWorkflow(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Response, error)
	DeleteWorkflowRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)

	// Workflow Definitions

	ListWorkflowDefinition(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.WorkflowDefinition, *oc.Response, error)
	ListWorkflowDefinitionRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error)

	GetWorkflowDefinition(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.WorkflowDefinition, *oc.Response, error)
	GetWorkflowDefinitionRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error)
}

type client struct {
	occ oc.Client
}

var _ Client = &client{}

func New(opencastClient oc.Client) *client {
	return &client{
		occ: opencastClient,
	}
}

func (c *client) Do(req *oc.Request) (*oc.Response, error) {
	if err := req.ApplyOptions(
		oc.WithHeader("Accept", AcceptJSONHeader),
	); err != nil {
		return nil, err
	}
	return c.occ.Do(req)
}

func (c *client) OpencastClient() oc.Client {
	return c.occ
}

func (c *client) GetAPI(ctx context.Context, opts ...oc.RequestOpts) (*extapiv1.API, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.API](
		c,
		func() (*oc.Request, error) { return c.GetAPIRequest(ctx, opts...) },
	)
}

func (c *client) GetAPIRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		ServiceType,
		"/api/",
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetAPIVersion(ctx context.Context, opts ...oc.RequestOpts) (*extapiv1.APIVersion, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.APIVersion](
		c,
		func() (*oc.Request, error) { return c.GetAPIVersionRequest(ctx, opts...) },
	)
}

func (c *client) GetAPIVersionRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		ServiceType,
		"/api/version",
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetAPIVersionDefault(ctx context.Context, opts ...oc.RequestOpts) (*extapiv1.APIVersion, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.APIVersion](
		c,
		func() (*oc.Request, error) { return c.GetAPIVersionDefaultRequest(ctx, opts...) },
	)
}

func (c *client) GetAPIVersionDefaultRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		ServiceType,
		"/api/default",
		oc.NoBody,
		opts...,
	)
}
