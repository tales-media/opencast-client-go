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

package opencastclientgo_test

import (
	"context"
	"fmt"
	"net/http"

	extapiv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11"
	extapiclientv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11/client"
	"shio.solutions/tales.media/opencast-client-go/apis/meta/base"
	"shio.solutions/tales.media/opencast-client-go/apis/meta/objlist"
	oc "shio.solutions/tales.media/opencast-client-go/client"
)

func Example_createEventUsingExternalAPI() {
	ctx := context.Background()

	// create Opencast client
	sm := &oc.StaticServiceMapper{
		Default: "https://stable.opencast.org",
	}
	client, err := oc.New(sm, oc.WithRequestOptions(
		oc.WithBasicAuth("admin", "opencast"),
	))
	if err != nil {
		panic(err)
	}
	extAPI := extapiclientv1.New(client)

	// download video and directly stream response body to Opencast
	resp, err := http.Get("https://radosgw.public.os.wwu.de/opencast-test-media/video-of-a-tabby-cat.mp4")
	if err != nil {
		panic(err)
	}

	// use client
	ids, _, err := extAPI.CreateEvent(
		ctx,
		&extapiclientv1.CreateEventRequestBody{
			ACL: extapiv1.ACL{
				{Role: "ROLE_ADMIN", Action: base.ReadAction, Allow: true},
				{Role: "ROLE_ADMIN", Action: base.WriteAction, Allow: true},
			},
			Metadata: []extapiv1.Catalog{{
				Flavor: base.DublinCoreEpisodeFlavor,
				Fields: []extapiv1.Field{{
					ID:    extapiv1.TitleFieldID,
					Value: extapiv1.TextFieldValue("Example"),
				}},
			}},
			Processing: &extapiv1.Processing{
				Workflow: "schedule-and-upload",
				Configuration: base.Properties{
					"straightToPublishing": "true",
				},
			},
			PresenterStreamFilename: "video-of-a-tabby-cat.mp4",
			PresenterStream:         resp.Body,
		},
	)
	if err != nil {
		panic(err)
	}
	if ids.Type != objlist.Object {
		panic("expected one created event")
	}

	fmt.Println("event created")
	// Output: event created
}
