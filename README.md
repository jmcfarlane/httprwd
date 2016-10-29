# httprwd

A simple [http.ResponseWriter](https://golang.org/pkg/net/http/#ResponseWriter)
delegate that keeps track of the response code.

## Installation

```
go get github.com/jmcfarlane/httprwd
```

## Usage

The typical use case is when you want to emit metrics about an http
route, and include the response code.

```go
d := &httprwd.ResponseWriterDelegate{ResponseWriter: w}
start := time.Now()
handler.ServeHTTP(d, r)
handlerDuration.WithLabelValues(
    r.URL.Path,
    r.Method,
    strconv.Itoa(d.Code),
).Observe(time.Since(start).Seconds())
```

## Working example

```
go run example/main.go
curl -s localhost:8080
curl -s localhost:8080/metrics
```
