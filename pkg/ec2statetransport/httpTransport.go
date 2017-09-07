package ec2statetransport

import (
  "bytes"
  "context"
  "encoding/json"
  "errors"
  "io/ioutil"
  "net/http"
  "net/url"
  "strings"

  "github.com/go-kit/kit/endpoint"
  "github.com/go-kit/kit/log"
  httptransport "github.com/go-kit/kit/transport/http"

  "github.com/timcurless/arbitrium/pkg/ec2stateservice"
  "github.com/timcurless/arbitrium/pkg/ec2stateendpoint"
)

// Returns an HTTP Handler that exposes a set of endpoints on defined paths
func NewHTTPHandler(endpoints ec2stateendpoint.Set, logger log.Logger) http.Handler {
  m := http.NewServeMux()
  m.Handle("/poweron", httptransport.NewServer(
    endpoints.PowerOnEndpoint,
    decodeHTTPPowerOnRequest,
    encodeHTTPGenericResponse,
  ))
  m.Handle("/poweroff", httptransport.NewServer(
    endpoints.PowerOffEndpoint,
    decodeHTTPPowerOffRequest,
    encodeHTTPGenericResponse,
  ))
  return m
}

// Returns a Ec2StateSvc Client backed by HTTP Server
// instance comes from a service discovery system like Consul
func NewHTTPClient(instance string, logger log.Logger) (ec2stateservice.Ec2StateSvc, error) {
  if !strings.HasPrefix(instance, "http") {
    instance = "http://" + instance
  }
  u, err := url.Parse(instance)
  if err != nil {
    return nil, err
  }

  var powerOnEndpoint endpoint.Endpoint
  {
    powerOnEndpoint = httptransport.NewClient(
      "POST",
      copyURL(u, "/poweron"),
      encodeHTTPGenericRequest,
      decodeHTTPPowerOnResponse,
    ).Endpoint()
  }

  var powerOffEndpoint endpoint.Endpoint
  {
    powerOffEndpoint = httptransport.NewClient(
      "POST",
      copyURL(u, "/poweroff"),
      encodeHTTPGenericRequest,
      decodeHTTPPowerOffResponse,
    ).Endpoint()
  }

  return ec2stateendpoint.Set{
    PowerOnEndpoint: powerOnEndpoint,
    PowerOffEndpoint: powerOffEndpoint,
  }, nil
}

func copyURL(base *url.URL, path string) *url.URL {
  next := *base
  next.Path = path
  return &next
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
  w.WriteHeader(err2code(err))
  json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func err2code(err error) int {
  // Going to need to beef up error handling. For now just return 500.
  return http.StatusInternalServerError
}

func errorDecode(r *http.Response) error {
  var w errorWrapper
  if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
    return err
  }
  return errors.New(w.Error)
}

type errorWrapper struct {
  Error string `json:"error"`
}

func decodeHTTPPowerOnRequest(_ context.Context, r *http.Request) (interface {}, error) {
  var req ec2stateendpoint.PowerOnRequest
  err := json.NewDecoder(r.Body).Decode(&req)
  return req, err
}

func decodeHTTPPowerOffRequest(_ context.Context, r *http.Request) (interface {}, error) {
  var req ec2stateendpoint.PowerOffRequest
  err := json.NewDecoder(r.Body).Decode(&req)
  return req, err
}

func decodeHTTPPowerOnResponse(_ context.Context, r *http.Response) (interface {}, error) {
  if r.StatusCode != http.StatusOK {
    return nil, errors.New(r.Status)
  }
  var resp ec2stateendpoint.PowerOnResponse
  err := json.NewDecoder(r.Body).Decode(&resp)
  return resp, err
}

func decodeHTTPPowerOffResponse(_ context.Context, r *http.Response) (interface {}, error) {
  if r.StatusCode != http.StatusOK {
    return nil, errors.New(r.Status)
  }
  var resp ec2stateendpoint.PowerOffResponse
  err := json.NewDecoder(r.Body).Decode(&resp)
  return resp, err
}

func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
  var buf bytes.Buffer
  if err := json.NewEncoder(&buf).Encode(request); err != nil {
    return err
  }
  r.Body = ioutil.NopCloser(&buf)
  return nil
}

func encodeHTTPGenericResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
