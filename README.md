# NPM-Go, a NPM HTTP API Client for Go

This library is a [NPM HTTP API] client for the Go language.

## Supported Go Versions

NPM-Go requires Go 1.6+.

## Project Maturity

NPM-Go is a fairly new library (started in January 2018)
as a dependancy for the NPM Vault Secret Engine. It only
impliments token and whoami.


## Installation

```
go get github.com/chrismatteson/npm-go
```


## Documentation

### Overview

To import the package:

``` go
import (
       "github.com/chrismatteson/npm-go"
)
```

All HTTP API operations are accessible via `npmgo.Client`, which
should be instantiated with `npm.go.NewClient`:

``` go
// URI, username, password
rmqc, _ = NewClient("http://127.0.0.1:15672", "guest", "guest")
```

TLS (HTTPS) can be enabled by adding an HTTP transport to the parameters
of `npmgo.NewTLSClient`:

``` go
transport := &http.Transport{TLSClientConfig: tlsConfig}
rmqc, _ := NewTLSClient("https://127.0.0.1:15672", "guest", "guest", transport)
```

NPM HTTP API has to be [configured to use TLS].


### Operations on Connections

``` go
xs, err := rmqc.ListConnections()
// => []ConnectionInfo, err

conn, err := rmqc.GetConnection("127.0.0.1:50545 -> 127.0.0.1:5672")
// => ConnectionInfo, err

// Forcefully close connection
_, err := rmqc.CloseConnection("127.0.0.1:50545 -> 127.0.0.1:5672")
// => *http.Response, err
```

### HTTPS Connections

``` go
var tlsConfig *tls.Config

...

transport := &http.Transport{TLSClientConfig: tlsConfig}

rmqc, err := NewTLSClient("https://127.0.0.1:15672", "guest", "guest", transport)
```

## License & Copyright

2-clause BSD license.

(c) Chris Matteson, 2013-2018.
