package ec2stateendpoint

import (
  "context"

  "github.com/go-kit/kit/endpoint"
  "github.com/go-kit/kit/log"

  "github.com/aws/aws-sdk-go/aws/session"

  "github.com/timcurless/arbitrium/pkg/ec2stateservice"
)

// Set is a helper struct, collecting all endpoints of all types into
// a single parameter.
type Set struct {
  PowerOnEndpoint endpoint.Endpoint
  PowerOffEndpoint endpoint.Endpoint
}

// Returns a new Set of endpoints
func New(sess *session.Session, svc ec2stateservice.Ec2StateSvc, logger log.Logger) Set {
  var powerOnEndpoint endpoint.Endpoint
  {
    powerOnEndpoint = MakePowerOnEndpoint(sess, svc)
    powerOnEndpoint = LoggingMiddleware(log.With(logger, "method", "poweron"))(powerOnEndpoint)
  }
  var powerOffEndpoint endpoint.Endpoint
  {
    powerOffEndpoint = MakePowerOffEndpoint(sess, svc)
    powerOffEndpoint = LoggingMiddleware(log.With(logger, "method", "poweroff"))(powerOffEndpoint)
  }
  return Set{
    PowerOnEndpoint: powerOnEndpoint,
    PowerOffEndpoint: powerOffEndpoint,
  }
}

// Service Interfaces
func (s Set) PowerOn(ctx context.Context, sess *session.Session, instanceId []*string) (interface{}, error) {
  resp, err := s.PowerOnEndpoint(ctx, PowerOnRequest{InstanceId: instanceId})
  if err != nil {
    return nil, err
  }
  response := resp.(PowerOnResponse)
  return response.Status, response.Err
}

func (s Set) PowerOff(ctx context.Context, sess *session.Session, instanceId []*string) (interface{}, error) {
  resp, err := s.PowerOffEndpoint(ctx, PowerOffRequest{InstanceId: instanceId})
  if err != nil {
    return nil, err
  }
  response := resp.(PowerOffResponse)
  return response.Status, response.Err
}

// Factory functions: Endpoints wrap service interfaces
func MakePowerOnEndpoint(sess *session.Session, svc ec2stateservice.Ec2StateSvc) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (response interface{}, err error) {
    req := request.(PowerOnRequest)
    res, err := svc.PowerOn(ctx, sess, req.InstanceId)
    return PowerOnResponse{Status: res, Err: err}, nil
  }
}

func MakePowerOffEndpoint(sess *session.Session, svc ec2stateservice.Ec2StateSvc) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (response interface{}, err error) {
    req := request.(PowerOffRequest)
    res, err := svc.PowerOff(ctx, sess, req.InstanceId)
    return PowerOffResponse{Status: res, Err: err}, nil
  }
}

// Request structs collect request parameters, response structs collect response values
type PowerOnRequest struct {
  InstanceId []*string `json:"instance-id"`
}

type PowerOnResponse struct {
  Status interface{} `json:"status"`
  Err    error `json:"err,omitempty"`
}

type PowerOffRequest struct {
  InstanceId []*string `json:"instance-id"`
}

type PowerOffResponse struct {
  Status interface{} `json:"status"`
  Err    error `json:"err,omitempty"`
}
