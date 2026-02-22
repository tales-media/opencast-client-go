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
	"time"

	extapiv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11"
	"shio.solutions/tales.media/opencast-client-go/apis/meta/base"
	oc "shio.solutions/tales.media/opencast-client-go/client"
	"shio.solutions/tales.media/opencast-client-go/pkg/multipart"
)

type CreateSeriesRequestBody struct {
	ACL      extapiv1.ACL
	Metadata []extapiv1.Catalog
	Theme    string
}

type UpdateSeriesRequestBody struct {
	Metadata []extapiv1.Catalog
}

type UpdateSeriesACLRequestBody struct {
	ACL      extapiv1.ACL
	Override bool
}

type UpdateSeriesMetadataRequestBody struct {
	Metadata []extapiv1.Value
}

type UpdateSeriesPropertiesRequestBody struct {
	Properties base.Properties
}

type WithSeriesOptions struct {
	WithACL             bool
	OnlyWithWriteAccess bool
}

var _ oc.RequestOpts = WithSeriesOptions{}

func (opt WithSeriesOptions) Apply(r *oc.Request) error {
	return r.ApplyOptions(
		oc.WithQuery("withacl", strconv.FormatBool(opt.WithACL)),
		oc.WithQuery("onlyWithWriteAccess", strconv.FormatBool(opt.OnlyWithWriteAccess)),
	)
}

type WithSeriesSearchOptions struct {
	Query        string
	Edit         bool
	FuzzyMatch   bool
	SeriesID     string
	SeriesTitle  string
	Creator      string
	Contributor  string
	Publisher    string
	RightsHolder string
	CreatedFrom  time.Time
	CreatedTo    time.Time
	Language     string
	License      string
	Subject      string
	Description  string
	Sort         SeriesSearchSortKey
	Offset       int
	Count        int
}

var _ oc.RequestOpts = WithSeriesSearchOptions{}

func (opt WithSeriesSearchOptions) Apply(r *oc.Request) error {
	opts := make([]oc.RequestOpts, 0, 18)
	if opt.Query != "" {
		opts = append(opts, oc.WithQuery("q", opt.Query))
	}
	if opt.Edit {
		opts = append(opts, oc.WithQuery("seriesId", strconv.FormatBool(opt.Edit)))
	}
	if opt.FuzzyMatch {
		opts = append(opts, oc.WithQuery("edit", strconv.FormatBool(opt.FuzzyMatch)))
	}
	if opt.SeriesID != "" {
		opts = append(opts, oc.WithQuery("fuzzyMatch", opt.SeriesID))
	}
	if opt.SeriesTitle != "" {
		opts = append(opts, oc.WithQuery("seriesTitle", opt.SeriesTitle))
	}
	if opt.Creator != "" {
		opts = append(opts, oc.WithQuery("creator", opt.Creator))
	}
	if opt.Contributor != "" {
		opts = append(opts, oc.WithQuery("contributor", opt.Contributor))
	}
	if opt.Publisher != "" {
		opts = append(opts, oc.WithQuery("publisher", opt.Publisher))
	}
	if opt.RightsHolder != "" {
		opts = append(opts, oc.WithQuery("rightsholder", opt.RightsHolder))
	}
	if !opt.CreatedFrom.IsZero() {
		opts = append(opts, oc.WithQuery("createdfrom", opt.CreatedFrom.Format(time.RFC3339)))
	}
	if !opt.CreatedTo.IsZero() {
		opts = append(opts, oc.WithQuery("createdto", opt.CreatedFrom.Format(time.RFC3339)))
	}
	if opt.Language != "" {
		opts = append(opts, oc.WithQuery("language", opt.Language))
	}
	if opt.License != "" {
		opts = append(opts, oc.WithQuery("license", opt.License))
	}
	if opt.Subject != "" {
		opts = append(opts, oc.WithQuery("subject", opt.Subject))
	}
	if opt.Description != "" {
		opts = append(opts, oc.WithQuery("description", opt.Description))
	}
	if string(opt.Sort) != "" {
		opts = append(opts, oc.WithQuery("sort", string(opt.Sort)))
	}
	if opt.Offset != 0 {
		opts = append(opts, oc.WithQuery("offset", strconv.Itoa(opt.Offset)))
	}
	if opt.Count != 0 {
		opts = append(opts, oc.WithQuery("count", strconv.Itoa(opt.Count)))
	}
	return r.ApplyOptions(opts...)
}

type SeriesSearchSortKey string

const (
	TitleSeriesSearchSortKey           = SeriesSearchSortKey("TITLE")
	TitleDescendingSeriesSearchSortKey = SeriesSearchSortKey("TITLE_DESC")
	// TODO: fix this in Opencast and define rest of sort keys
	// https://github.com/opencast/opencast/pull/3204/files#diff-bdb07894e2e2cc79259310945d64de171aac86a22006fb9fbbb394a127f46436
)

const (
	SeriesManagedAclFilterKey   = FilterKey("managedAcl")
	SeriesContributorsFilterKey = FilterKey("contributors")
	SeriesCreationDateFilterKey = FilterKey("CreationDate")
	SeriesTextFilterFilterKey   = FilterKey("textFilter")
	SeriesLanguageFilterKey     = FilterKey("language")
	SeriesLicenseFilterKey      = FilterKey("license")
	SeriesOrganizersFilterKey   = FilterKey("organizers")
	SeriesSubjectFilterKey      = FilterKey("subject")
	SeriesTitleFilterKey        = FilterKey("title")
	SeriesIdentifierFilterKey   = FilterKey("identifier")
	SeriesDescriptionFilterKey  = FilterKey("description")
	SeriesCreatorFilterKey      = FilterKey("creator")
	SeriesPublishersFilterKey   = FilterKey("publishers")
	SeriesRightsHolderFilterKey = FilterKey("rightsholder")
)

const (
	SeriesTitleSortKey        = SortKey("title")
	SeriesContributorsSortKey = SortKey("contributors")
	SeriesCreatorSortKey      = SortKey("creator")
	SeriesCreatedSortKey      = SortKey("created")
)

func (c *client) ListSeries(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.Series, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[[]extapiv1.Series](
		c,
		func() (*oc.Request, error) { return c.ListSeriesRequest(ctx, opts...) },
	)
}

func (c *client) ListSeriesRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		SeriesServiceType,
		"/api/series",
		oc.NoBody,
		opts...,
	)
}

func (c *client) SearchSeries(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.Series, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[[]extapiv1.Series](
		c,
		func() (*oc.Request, error) { return c.SearchSeriesRequest(ctx, opts...) },
	)
}

func (c *client) SearchSeriesRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		SeriesServiceType,
		"/api/series/series.json",
		oc.NoBody,
		opts...,
	)
}

func (c *client) CreateSeries(ctx context.Context, body *CreateSeriesRequestBody, opts ...oc.RequestOpts) (extapiv1.Identifier, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[extapiv1.Identifier](
		c,
		func() (*oc.Request, error) { return c.CreateSeriesRequest(ctx, body, opts...) },
	)
}

func (c *client) CreateSeriesRequest(ctx context.Context, body *CreateSeriesRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
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

	if body.Theme != "" {
		mp.AddPart(multipart.FormFieldString("theme", body.Theme))
	}

	return oc.NewRequest(
		ctx,
		http.MethodPost,
		SeriesServiceType,
		"/api/series",
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) GetSeries(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.Series, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.Series](
		c,
		func() (*oc.Request, error) { return c.GetSeriesRequest(ctx, id, opts...) },
	)
}

func (c *client) GetSeriesRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		SeriesServiceType,
		"/api/series/"+url.PathEscape(id),
		oc.NoBody,
		opts...,
	)
}

func (c *client) UpdateSeries(ctx context.Context, id string, body *UpdateSeriesRequestBody, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.UpdateSeriesRequest(ctx, id, body, opts...) },
	)
}

func (c *client) UpdateSeriesRequest(ctx context.Context, id string, body *UpdateSeriesRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()

	if len(body.Metadata) > 0 {
		metadata, err := json.Marshal(body.Metadata)
		if err != nil {
			return nil, err
		}
		mp.AddPart(multipart.FormField("metadata", metadata))
	}

	return oc.NewRequest(
		ctx,
		http.MethodPut,
		SeriesServiceType,
		"/api/series/"+url.PathEscape(id),
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) DeleteSeries(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.DeleteSeriesRequest(ctx, id, opts...) },
	)
}

func (c *client) DeleteSeriesRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodDelete,
		SeriesServiceType,
		"/api/series/"+url.PathEscape(id),
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetSeriesACL(ctx context.Context, id string, opts ...oc.RequestOpts) (extapiv1.ACL, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[extapiv1.ACL](
		c,
		func() (*oc.Request, error) { return c.GetSeriesACLRequest(ctx, id, opts...) },
	)
}

func (c *client) GetSeriesACLRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		SeriesServiceType,
		"/api/series/"+url.PathEscape(id)+"/acl",
		oc.NoBody,
		opts...,
	)
}

func (c *client) UpdateSeriesACL(ctx context.Context, id string, body *UpdateSeriesACLRequestBody, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.UpdateSeriesACLRequest(ctx, id, body, opts...) },
	)
}

func (c *client) UpdateSeriesACLRequest(ctx context.Context, id string, body *UpdateSeriesACLRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	acl, err := json.Marshal(body.ACL)
	if err != nil {
		return nil, err
	}
	mp.AddPart(multipart.FormField("acl", acl))
	if body.Override {
		mp.AddPart(multipart.FormFieldString("override", strconv.FormatBool(body.Override)))
	}
	return oc.NewRequest(
		ctx,
		http.MethodPut,
		SeriesServiceType,
		"/api/series/"+url.PathEscape(id)+"/acl",
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) ListSeriesMetadata(ctx context.Context, id string, opts ...oc.RequestOpts) ([]extapiv1.Catalog, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[[]extapiv1.Catalog](
		c,
		func() (*oc.Request, error) { return c.ListSeriesMetadataRequest(ctx, id, opts...) },
	)
}

func (c *client) ListSeriesMetadataRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	req, err := oc.NewRequest(
		ctx,
		http.MethodGet,
		SeriesServiceType,
		"/api/series/"+url.PathEscape(id)+"/metadata",
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

func (c *client) GetSeriesMetadata(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) (*extapiv1.Catalog, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.Catalog](
		c,
		func() (*oc.Request, error) { return c.GetSeriesMetadataRequest(ctx, id, flavor, opts...) },
	)
}

func (c *client) GetSeriesMetadataRequest(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) (*oc.Request, error) {
	req, err := oc.NewRequest(
		ctx,
		http.MethodGet,
		SeriesServiceType,
		"/api/series/"+url.PathEscape(id)+"/metadata",
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

func (c *client) UpdateSeriesMetadata(ctx context.Context, id string, flavor base.Flavor, body *UpdateSeriesMetadataRequestBody, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.UpdateSeriesMetadataRequest(ctx, id, flavor, body, opts...) },
	)
}

func (c *client) UpdateSeriesMetadataRequest(ctx context.Context, id string, flavor base.Flavor, body *UpdateSeriesMetadataRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	metadata, err := json.Marshal(body.Metadata)
	if err != nil {
		return nil, err
	}
	mp.AddPart(multipart.FormField("metadata", metadata))
	req, err := oc.NewRequest(
		ctx,
		http.MethodPut,
		SeriesServiceType,
		"/api/series/"+url.PathEscape(id)+"/metadata",
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

func (c *client) DeleteSeriesMetadata(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.DeleteSeriesMetadataRequest(ctx, id, flavor, opts...) },
	)
}

func (c *client) DeleteSeriesMetadataRequest(ctx context.Context, id string, flavor base.Flavor, opts ...oc.RequestOpts) (*oc.Request, error) {
	req, err := oc.NewRequest(
		ctx,
		http.MethodDelete,
		SeriesServiceType,
		"/api/series/"+url.PathEscape(id)+"/metadata",
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

func (c *client) GetSeriesProperties(ctx context.Context, id string, opts ...oc.RequestOpts) (base.Properties, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[base.Properties](
		c,
		func() (*oc.Request, error) { return c.GetSeriesPropertiesRequest(ctx, id, opts...) },
	)
}

func (c *client) GetSeriesPropertiesRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		SeriesServiceType,
		"/api/series/"+url.PathEscape(id)+"/properties",
		oc.NoBody,
		opts...,
	)
}

func (c *client) UpdateSeriesProperties(ctx context.Context, id string, body *UpdateSeriesPropertiesRequestBody, opts ...oc.RequestOpts) (base.Properties, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[base.Properties](
		c,
		func() (*oc.Request, error) { return c.UpdateSeriesPropertiesRequest(ctx, id, body, opts...) },
	)
}

func (c *client) UpdateSeriesPropertiesRequest(ctx context.Context, id string, body *UpdateSeriesPropertiesRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	properties, err := json.Marshal(body.Properties)
	if err != nil {
		return nil, err
	}
	mp.AddPart(multipart.FormField("properties", properties))
	return oc.NewRequest(
		ctx,
		http.MethodPut,
		SeriesServiceType,
		"/api/series/"+url.PathEscape(id)+"/properties",
		oc.NewMultipartBody(mp),
		opts...,
	)
}
