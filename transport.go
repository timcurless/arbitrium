package main

import (
  "context"
  "encoding/json"
  "net/http"

  "github.com/go-kit/kit/endpoint"

  "github.com/aws/aws-sdk-go/aws/session"
)

func makePowerOnEndpoint(sess *session.Session, svc Ec2StateSvc) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (interface{}, error) {
    req := request.(powerOnRequest)

    res, err := svc.PowerOn(sess, req.InstanceId)
    if err != nil {
      return powerOnResponse{res, err.Error()}, nil
    }
    return powerOnResponse{res, ""}, nil
  }
}

func decodePowerOnRequest(_ context.Context, r *http.Request) (interface {}, error) {
  var request powerOnRequest
  if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
    return nil, err
  }
  return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
  return json.NewEncoder(w).Encode(response)
}

type powerOnRequest struct {
  InstanceId []*string `json:"instance-id"`
}

type powerOnResponse struct {
  Status interface{} `json:"status"`
  Err    string `json:"err,omitempty"`
}
