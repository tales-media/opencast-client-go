# Opencast Client Library for Go

[![Go Report Card](https://goreportcard.com/badge/shio.solutions/tales.media/opencast-client-go?style=flat-square)](https://goreportcard.com/report/shio.solutions/tales.media/opencast-client-go)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.26-61CFDD.svg?style=flat-square)
[![PkgGoDev](https://pkg.go.dev/badge/mod/shio.solutions/tales.media/opencast-client-go)](https://pkg.go.dev/shio.solutions/tales.media/opencast-client-go)

Go client library for talking to an [Opencast](https://opencast.org/) cluster.

## Install

To get the latest version, use go1.26+ and fetch using the `go get` command:

```sh
$ go get shio.solutions/tales.media/opencast-client-go@latest
```

## Quickstart

First, create a service mapper instance. The service mapper determines what service, i.e. what specific URL, should be used when sending a request to a particular service type. You have two options. A static service mapper allows to configure a custom static mapping of service types to service endpoints or using a default. Alternatively, a dynamic service mapper consults the Opencast service registry to dynamically determine the service endpoint.

```go
// static service mapper

sm := &oc.StaticServiceMapper{
	Default: "https://stable.opencast.org",
	ServiceHost: map[string]string{
		"org.opencastproject.external.events": "https://api.example.com",
	},
}

// -----

// dynamic service mapper

// first create Opencast client to connect to the service registry
staticSM := &oc.StaticServiceMapper{Default: "https://stable.opencast.org"}
staticClient, err := oc.New(sm, oc.WithRequestOptions(oc.WithBasicAuth("admin", "opencast")))

// and then create the dynamic service mapper with a caching TTL of 10m
sm := oc.NewDynamicServiceMapper(staticClient, 10*time.Minute)
```

After that, you can create the basic Opencast client.

```go
client, err := oc.New(sm, oc.WithRequestOptions(
	oc.WithBasicAuth("admin", "opencast"),
	// or oc.WithJWTQuery(myJWT)
	// or oc.WithJWTHeader("Authorization", "Bearer ", myJWT)
))

req, err := oc.NewRequest(
	context.Background(),
	http.MethodGet,
	"org.opencastproject.external.events",
	"/api/events",
	oc.NoBody,
)
resp, err := client.Do(req)
```

You can also create an External API client provides type-safe access.

```go
extAPI := extapiclientv1.New(client)

// resp contains raw HTTP response
events, resp, err := extAPI.ListEvent(
	context.Background(),
	extapiclientv1.WithEventOptions{
		WithACL:                    false,
		WithMetadata:               false,
		WithScheduling:             false,
		WithPublications:           false,
		IncludeInternalPublication: false,
		OnlyWithWriteAccess:        true,
	},
)
for _, event := range events {
	fmt.Printf("%s: %s", event.Identifier, event.Title)
}
```

Note that list requests are generally paginated and thus only return one page. You can use helper methods to retrieve all resources.

```go
var allEvents []extapiv1.Event
err := oc.Paginate(
	extAPI,
	func(i int) (*oc.Request, error) {
		return extAPI.ListEventRequest(
			context.Background(),
			extapiclientv1.WithPagination{
				Limit:  100,
				Offset: i * 100,
			},
			extapiclientv1.WithEventOptions{
				WithACL:                    false,
				WithMetadata:               false,
				WithScheduling:             false,
				WithPublications:           false,
				IncludeInternalPublication: false,
				OnlyWithWriteAccess:        true,
			},
		)
	},
	oc.CollectAllPages(&allEvents),
)
```

## License

Apache 2.0 (c) shio solutions GmbH
