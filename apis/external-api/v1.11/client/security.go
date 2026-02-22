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
	"time"

	extapiv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11"
	oc "shio.solutions/tales.media/opencast-client-go/client"
	"shio.solutions/tales.media/opencast-client-go/pkg/multipart"
)

type SignURLRequestBody struct {
	URL         string
	ValidUntil  time.Time
	ValidSource string
}

func (c *client) SignURL(ctx context.Context, body *SignURLRequestBody, opts ...oc.RequestOpts) (*extapiv1.SignedURL, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.SignedURL](
		c,
		func() (*oc.Request, error) { return c.SignURLRequest(ctx, body, opts...) },
	)
}

func (c *client) SignURLRequest(ctx context.Context, body *SignURLRequestBody, opts ...oc.RequestOpts) (*oc.Request, error) {
	mp := multipart.New()
	mp.AddPart(multipart.FormFieldString("url", body.URL))
	if !body.ValidUntil.IsZero() {
		mp.AddPart(multipart.FormFieldString("valid-until", body.ValidUntil.Format(time.RFC3339)))
	}
	if body.ValidSource != "" {
		mp.AddPart(multipart.FormFieldString("valid-source", body.ValidSource))
	}
	return oc.NewRequest(
		ctx,
		http.MethodPost,
		SecurityServiceType,
		"/api/security/sign",
		oc.NewMultipartBody(mp),
		opts...,
	)
}
