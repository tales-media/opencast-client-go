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
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	extapiv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11"
	"shio.solutions/tales.media/opencast-client-go/apis/meta/base"
	"shio.solutions/tales.media/opencast-client-go/apis/meta/objlist"
	oc "shio.solutions/tales.media/opencast-client-go/client"
	"shio.solutions/tales.media/opencast-client-go/pkg/multipart"
)

type CreateEventRequestBody struct {
	ACL                        extapiv1.ACL
	Metadata                   []extapiv1.Catalog
	Scheduling                 *extapiv1.SchedulingRequest
	Processing                 *extapiv1.Processing
	PresenterFile              string
	PresenterStream            io.ReadCloser
	PresenterStreamFilename    string
	PresentationFile           string
	PresentationStream         io.ReadCloser
	PresentationStreamFilename string
	AudioFile                  string
	AudioStream                io.ReadCloser
	AudioStreamFilename        string
}

type UpdateEventRequestBody struct {
	ACL        extapiv1.ACL
	Metadata   []extapiv1.Catalog
	Scheduling *extapiv1.SchedulingRequest
	Processing *extapiv1.Processing
}

type UpdateEventACLRequestBody struct {
	ACL extapiv1.ACL
}

type CreateEventTrackRequestBody struct {
	Flavor              base.Flavor
	OverwriteExisting   bool
	Tags                []string
	TrackFile           string // TODO: allow option to use io.Reader
	TrackStream         io.ReadCloser
	TrackStreamFilename string
}

type UpdateEventMetadataRequestBody struct {
	Metadata []extapiv1.Value
}

type UpdateEventSchedulingRequestBody struct {
	Scheduling extapiv1.SchedulingRequest
}

type WithEventOptions struct {
	WithACL                    bool
	WithMetadata               bool
	WithScheduling             bool
	WithPublications           bool
	IncludeInternalPublication bool
	OnlyWithWriteAccess        bool
	AllowConflict              bool
}

var _ oc.RequestOpts = WithEventOptions{}

func (opt WithEventOptions) Apply(r *oc.Request) error {
	return r.ApplyOptions(
		oc.WithQuery("withacl", strconv.FormatBool(opt.WithACL)),
		oc.WithQuery("withmetadata", strconv.FormatBool(opt.WithMetadata)),
		oc.WithQuery("withscheduling", strconv.FormatBool(opt.WithScheduling)),
		oc.WithQuery("withpublications", strconv.FormatBool(opt.WithPublications)),
		oc.WithQuery("includeInternalPublication", strconv.FormatBool(opt.IncludeInternalPublication)),
		oc.WithQuery("onlyWithWriteAccess", strconv.FormatBool(opt.OnlyWithWriteAccess)),
		oc.WithQuery("allowConflict", strconv.FormatBool(opt.AllowConflict)),
	)
}

const (
	EventPresentersFilterKey     = FilterKey("presenters")
	EventContributorsFilterKey   = FilterKey("contributors")
	EventLocationFilterKey       = FilterKey("location")
	EventSeriesFilterKey         = FilterKey("series")
	EventSubjectFilterKey        = FilterKey("subject")
	EventTextFilterFilterKey     = FilterKey("textFilter")
	EventIdentifierFilterKey     = FilterKey("identifier")
	EventTitleFilterKey          = FilterKey("title")
	EventDescriptionFilterKey    = FilterKey("description")
	EventSeriesNameFilterKey     = FilterKey("series_name")
	EventLanguageFilterKey       = FilterKey("language")
	EventCreatedFilterKey        = FilterKey("created")
	EventLicenseFilterKey        = FilterKey("license")
	EventRightsHolderFilterKey   = FilterKey("rightsholder")
	EventStatusFilterKey         = FilterKey("status")
	EventIsPartOfFilterKey       = FilterKey("is_part_of")
	EventSourceFilterKey         = FilterKey("source")
	EventAgentIdFilterKey        = FilterKey("agent_id")
	EventStartFilterKey          = FilterKey("start")
	EventTechnicalStartFilterKey = FilterKey("technical_start")
)

const (
	EventTitleSortKey            = SortKey("title")
	EventPresenterSortKey        = SortKey("presenter")
	EventStartDateSortKey        = SortKey("start_date")
	EventEndDateSortKey          = SortKey("end_date")
	EventReviewStatusSortKey     = SortKey("review_status")
	EventWorkflowStateSortKey    = SortKey("workflow_state")
	EventSchedulingStatusSortKey = SortKey("scheduling_status")
	EventSeriesNameSortKey       = SortKey("series_name")
	EventLocationSortKey         = SortKey("location")
	EventTechnicalDateSortKey    = SortKey("technical_date")
	EventTechnicalStartSortKey   = SortKey("technical_start")
	EventTechnicalEndSortKey     = SortKey("technical_end")
)

func (c *client) ListEvent(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.Event, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[[]extapiv1.Event](
		c,
		func() (*oc.Request, error) { return c.ListEventRequest(ctx, opts...) },
	)
}

func (c *client) ListEventRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		EventsServiceType,
		"/api/events",
		oc.NoBody,
		opts...,
	)
}

func (c *client) CreateEvent(ctx context.Context, body *CreateEventRequestBody, opts ...oc.RequestOpts) (objlist.ObjectOrList[extapiv1.Identifier], *oc.Response, error) {
	return oc.GenericAutoDecodedDo[objlist.ObjectOrList[extapiv1.Identifier]](
		c,
		func() (*oc.Request, error) { return c.CreateEventRequest(ctx, body, opts...) },
	)
}

func (c *client) CreateEventRequest(ctx context.Context, body *CreateEventRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()

	if len(body.ACL) > 0 {
		acl, err := json.Marshal(body.ACL)
		if err != nil {
			return nil, err
		}
		mp.AddPart(multipart.FormField("acl", acl))
	}

	if len(body.Metadata) > 0 {
		metadata, err := json.Marshal(body.Metadata)
		if err != nil {
			return nil, err
		}
		mp.AddPart(multipart.FormField("metadata", metadata))
	}

	if body.Scheduling != nil {
		scheduling, err := json.Marshal(body.Scheduling)
		if err != nil {
			return nil, err
		}
		mp.AddPart(multipart.FormField("scheduling", scheduling))
	}

	if body.Processing != nil {
		processing, err := json.Marshal(body.Processing)
		if err != nil {
			return nil, err
		}
		mp.AddPart(multipart.FormField("processing", processing))
	}

	if body.PresenterFile != "" {
		mp.AddPart(multipart.File("presenter", body.PresenterFile))
	} else if body.PresenterStream != nil {
		mp.AddPart(multipart.Stream("presenter", body.PresenterStreamFilename, body.PresenterStream))
	}

	if body.PresentationFile != "" {
		mp.AddPart(multipart.File("presentation", body.PresentationFile))
	} else if body.PresentationStream != nil {
		mp.AddPart(multipart.Stream("presentation", body.PresentationStreamFilename, body.PresentationStream))
	}

	if body.AudioFile != "" {
		mp.AddPart(multipart.File("audio", body.AudioFile))
	} else if body.AudioStream != nil {
		mp.AddPart(multipart.Stream("audio", body.AudioStreamFilename, body.AudioStream))
	}

	return oc.NewRequest(
		ctx,
		http.MethodPost,
		EventsServiceType,
		"/api/events",
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) GetEvent(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.Event, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.Event](
		c,
		func() (*oc.Request, error) { return c.GetEventRequest(ctx, id, opts...) },
	)
}

func (c *client) GetEventRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id),
		oc.NoBody,
		opts...,
	)
}

func (c *client) UpdateEvent(ctx context.Context, id string, body *UpdateEventRequestBody, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.UpdateEventRequest(ctx, id, body, opts...) },
	)
}

func (c *client) UpdateEventRequest(ctx context.Context, id string, body *UpdateEventRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()

	if len(body.ACL) > 0 {
		acl, err := json.Marshal(body.ACL)
		if err != nil {
			return nil, err
		}
		mp.AddPart(multipart.FormField("acl", acl))
	}

	if len(body.Metadata) > 0 {
		metadata, err := json.Marshal(body.Metadata)
		if err != nil {
			return nil, err
		}
		mp.AddPart(multipart.FormField("metadata", metadata))
	}

	if body.Scheduling != nil {
		scheduling, err := json.Marshal(body.Scheduling)
		if err != nil {
			return nil, err
		}
		mp.AddPart(multipart.FormField("scheduling", scheduling))
	}

	if body.Processing != nil {
		processing, err := json.Marshal(body.Processing)
		if err != nil {
			return nil, err
		}
		mp.AddPart(multipart.FormField("processing", processing))
	}

	return oc.NewRequest(
		ctx,
		http.MethodPost,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id),
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) DeleteEvent(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.DeleteEventRequest(ctx, id, opts...) },
	)
}

func (c *client) DeleteEventRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodDelete,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id),
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetEventACL(ctx context.Context, id string, opts ...oc.RequestOpts) (extapiv1.ACL, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[extapiv1.ACL](
		c,
		func() (*oc.Request, error) { return c.GetEventACLRequest(ctx, id, opts...) },
	)
}

func (c *client) GetEventACLRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id)+"/acl",
		oc.NoBody,
		opts...,
	)
}

func (c *client) UpdateEventACL(ctx context.Context, id string, body *UpdateEventACLRequestBody, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.UpdateEventACLRequest(ctx, id, body, opts...) },
	)
}

func (c *client) UpdateEventACLRequest(ctx context.Context, id string, body *UpdateEventACLRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	acl, err := json.Marshal(body.ACL)
	if err != nil {
		return nil, err
	}
	mp.AddPart(multipart.FormField("acl", acl))
	return oc.NewRequest(
		ctx,
		http.MethodPut,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id)+"/acl",
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) CreateEventACE(ctx context.Context, id string, action base.Action, role string, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.CreateEventACERequest(ctx, id, action, role, opts...) },
	)
}

func (c *client) CreateEventACERequest(ctx context.Context, id string, action base.Action, role string, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	mp.AddPart(multipart.FormFieldString("role", role))
	return oc.NewRequest(
		ctx,
		http.MethodPost,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id)+"/acl/"+url.PathEscape(string(action)),
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) DeleteEventACE(ctx context.Context, id string, action base.Action, role string, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.DeleteEventACERequest(ctx, id, action, role, opts...) },
	)
}

func (c *client) DeleteEventACERequest(ctx context.Context, id string, action base.Action, role string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodDelete,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id)+"/acl/"+url.PathEscape(string(action))+"/"+url.PathEscape(role),
		oc.NoBody,
		opts...,
	)
}

func (c *client) ListEventMedia(ctx context.Context, id string, opts ...oc.RequestOpts) ([]extapiv1.MediaTrackElement, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[[]extapiv1.MediaTrackElement](
		c,
		func() (*oc.Request, error) { return c.ListEventMediaRequest(ctx, id, opts...) },
	)
}

func (c *client) ListEventMediaRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id)+"/media",
		oc.NoBody,
		opts...,
	)
}

func (c *client) CreateEventTrack(ctx context.Context, id string, body *CreateEventTrackRequestBody, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.CreateEventTrackRequest(ctx, id, body, opts...) },
	)
}

func (c *client) CreateEventTrackRequest(ctx context.Context, id string, body *CreateEventTrackRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	mp.AddParts(
		multipart.FormFieldString("flavor", string(body.Flavor)),
		multipart.FormFieldString("tags", strings.Join(body.Tags, ",")),
		multipart.FormFieldString("overwriteExisting", strconv.FormatBool(body.OverwriteExisting)),
	)
	if body.TrackFile != "" {
		mp.AddPart(multipart.File("track", body.TrackFile))
	} else if body.TrackStream != nil {
		mp.AddPart(multipart.Stream("track", body.TrackStreamFilename, body.TrackStream))
	}
	return oc.NewRequest(
		ctx,
		http.MethodPost,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id)+"/track",
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) ListEventMetadata(ctx context.Context, id string, opts ...oc.RequestOpts) ([]extapiv1.Catalog, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[[]extapiv1.Catalog](
		c,
		func() (*oc.Request, error) { return c.ListEventMetadataRequest(ctx, id, opts...) },
	)
}

func (c *client) ListEventMetadataRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	req, err := oc.NewRequest(
		ctx,
		http.MethodGet,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id)+"/metadata",
		oc.NoBody,
		opts...,
	)
	if err != nil {
		return nil, err
	}
	// the "type" query parameter must not be set (otherwise a single [extapiv1.Catalog] is returned)
	if err := req.ApplyOptions(oc.WithoutQuery("type")); err != nil {
		return nil, err
	}
	return req, nil
}

func (c *client) GetEventMetadata(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) ([]extapiv1.Field, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[[]extapiv1.Field](
		c,
		func() (*oc.Request, error) { return c.GetEventMetadataRequest(ctx, id, flavor, opts...) },
	)
}

func (c *client) GetEventMetadataRequest(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) (*oc.Request, error) {
	req, err := oc.NewRequest(
		ctx,
		http.MethodGet,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id)+"/metadata",
		oc.NoBody,
		opts...,
	)
	if err != nil {
		return nil, err
	}
	if err := req.ApplyOptions(oc.WithQuery("type", string(flavor))); err != nil {
		return nil, err
	}
	return req, nil
}

func (c *client) UpdateEventMetadata(ctx context.Context, id string, flavor base.Flavor, body *UpdateEventMetadataRequestBody, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.UpdateEventMetadataRequest(ctx, id, flavor, body, opts...) },
	)
}

func (c *client) UpdateEventMetadataRequest(ctx context.Context, id string, flavor base.Flavor, body *UpdateEventMetadataRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	metadata, err := json.Marshal(body.Metadata)
	if err != nil {
		return nil, err
	}
	mp.AddPart(multipart.FormField("metadata", metadata))
	req, err := oc.NewRequest(
		ctx,
		http.MethodPut,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id)+"/metadata",
		oc.NewMultipartBody(mp),
		opts...,
	)
	if err != nil {
		return nil, err
	}
	if err := req.ApplyOptions(oc.WithQuery("type", string(flavor))); err != nil {
		return nil, err
	}
	return req, nil
}

func (c *client) DeleteEventMetadata(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.DeleteEventMetadataRequest(ctx, id, flavor, opts...) },
	)
}

func (c *client) DeleteEventMetadataRequest(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) (*oc.Request, error) {
	req, err := oc.NewRequest(
		ctx,
		http.MethodDelete,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id)+"/metadata",
		oc.NoBody,
		opts...,
	)
	if err != nil {
		return nil, err
	}
	if err := req.ApplyOptions(oc.WithQuery("type", string(flavor))); err != nil {
		return nil, err
	}
	return req, nil
}

func (c *client) ListEventPublication(ctx context.Context, id string, opts ...oc.RequestOpts) ([]extapiv1.Publication, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[[]extapiv1.Publication](
		c,
		func() (*oc.Request, error) { return c.ListEventPublicationRequest(ctx, id, opts...) },
	)
}

func (c *client) ListEventPublicationRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id)+"/publications",
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetEventPublication(ctx context.Context, id string, publicationID string, opts ...oc.RequestOpts) (*extapiv1.Publication, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.Publication](
		c,
		func() (*oc.Request, error) { return c.GetEventPublicationRequest(ctx, id, publicationID, opts...) },
	)
}

func (c *client) GetEventPublicationRequest(ctx context.Context, id string, publicationID string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id)+"/publications/"+url.PathEscape(publicationID),
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetEventScheduling(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.Scheduling, *oc.Response, error) {
	// BUG: if there is no scheduling Opencast will return a 204 No Content
	return oc.GenericAutoDecodedDo[*extapiv1.Scheduling](
		c,
		func() (*oc.Request, error) { return c.GetEventSchedulingRequest(ctx, id, opts...) },
	)
}

func (c *client) GetEventSchedulingRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id)+"/scheduling",
		oc.NoBody,
		opts...,
	)
}

func (c *client) UpdateEventScheduling(ctx context.Context, id string, body *UpdateEventSchedulingRequestBody, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.UpdateEventSchedulingRequest(ctx, id, body, opts...) },
	)
}

func (c *client) UpdateEventSchedulingRequest(ctx context.Context, id string, body *UpdateEventSchedulingRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	scheduling, err := json.Marshal(body.Scheduling)
	if err != nil {
		return nil, err
	}
	mp.AddPart(multipart.FormField("scheduling", scheduling))
	return oc.NewRequest(
		ctx,
		http.MethodPut,
		EventsServiceType,
		"/api/events/"+url.PathEscape(id)+"/scheduling",
		oc.NewMultipartBody(mp),
		opts...,
	)
}
