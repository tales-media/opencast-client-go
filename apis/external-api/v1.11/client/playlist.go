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

	extapiv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11"
	oc "shio.solutions/tales.media/opencast-client-go/client"
	"shio.solutions/tales.media/opencast-client-go/pkg/multipart"
)

type CreatePlaylistRequestBody struct {
	Playlist extapiv1.Playlist
}

type UpdatePlaylistRequestBody struct {
	Playlist extapiv1.Playlist
}

const (
	PlaylistUpdatedSortKey = SortKey("updated")
)

func (c *client) ListPlaylist(ctx context.Context, opts ...oc.RequestOpts) ([]extapiv1.Playlist, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[[]extapiv1.Playlist](
		c,
		func() (*oc.Request, error) { return c.ListPlaylistRequest(ctx, opts...) },
	)
}

func (c *client) ListPlaylistRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		PlaylistsServiceType,
		"/api/playlists",
		oc.NoBody,
		opts...,
	)
}

func (c *client) CreatePlaylist(ctx context.Context, body *CreatePlaylistRequestBody, opts ...oc.RequestOpts) (*extapiv1.Playlist, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.Playlist](
		c,
		func() (*oc.Request, error) { return c.CreatePlaylistRequest(ctx, body, opts...) },
	)
}

func (c *client) CreatePlaylistRequest(ctx context.Context, body *CreatePlaylistRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	playlist, err := json.Marshal(body.Playlist)
	if err != nil {
		return nil, err
	}
	mp.AddPart(multipart.FormField("playlist", playlist))
	return oc.NewRequest(
		ctx,
		http.MethodPost,
		PlaylistsServiceType,
		"/api/playlists",
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) GetPlaylist(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.Playlist, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.Playlist](
		c,
		func() (*oc.Request, error) { return c.GetPlaylistRequest(ctx, id, opts...) },
	)
}

func (c *client) GetPlaylistRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		PlaylistsServiceType,
		"/api/playlists/"+url.PathEscape(id),
		oc.NoBody,
		opts...,
	)
}

func (c *client) UpdatePlaylist(ctx context.Context, id string, body *UpdatePlaylistRequestBody, opts ...oc.RequestOpts) (*extapiv1.Playlist, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.Playlist](
		c,
		func() (*oc.Request, error) { return c.UpdatePlaylistRequest(ctx, id, body, opts...) },
	)
}

func (c *client) UpdatePlaylistRequest(ctx context.Context, id string, body *UpdatePlaylistRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	playlist, err := json.Marshal(body.Playlist)
	if err != nil {
		return nil, err
	}
	mp.AddPart(multipart.FormField("playlist", playlist))
	return oc.NewRequest(
		ctx,
		http.MethodPost,
		PlaylistsServiceType,
		"/api/playlists/"+url.PathEscape(id),
		oc.NewMultipartBody(mp),
		opts...,
	)
}

func (c *client) DeletePlaylist(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.Playlist, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.Playlist](
		c,
		func() (*oc.Request, error) { return c.DeletePlaylistRequest(ctx, id, opts...) },
	)
}

func (c *client) DeletePlaylistRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodDelete,
		PlaylistsServiceType,
		"/api/playlists/"+url.PathEscape(id),
		oc.NoBody,
		opts...,
	)
}
