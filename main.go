package main

import (
  "context"
  "encoding/json"
  "errors"
  "flag"
  "net/http"
  "os"

  "github.com/go-kit/kit/endpoint"
  "github.com/go-kit/kit/log"
  httptransport "github.com/go-kit/kit/transport/http"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
  var (
    listen = flag.String("listen", ":8080", "HTTP listen address")
  )
  flag.Parse()

  const (
    DefaultSharedConfigProfile = `default`
  )

  var logger log.Logger
  {
    logger = log.NewLogfmtLogger(os.Stderr)
    logger = log.With(logger, "ts", log.DefaultTimestampUTC)
    logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)
  }

  // Create AWS SDK session
  sess := session.Must(session.NewSession())

  // Create the EC2 State Service
  svc := ec2StateSvc{}

  powerOnHandler := httptransport.NewServer(
    makePowerOnEndpoint(sess, svc),
    decodePowerOnRequest,
    encodeResponse,
  )

  http.Handle("/poweron", powerOnHandler)
  logger.Log("msg", "HTTP", "addr", *listen)
  logger.Log("err", http.ListenAndServe(*listen, nil))
}

type ec2StateSvc struct {}

type Ec2StateSvc interface {
  PowerOn(*session.Session, string) (interface{}, error)
}

func (ec2StateSvc) PowerOn(sess *session.Session, instanceId string) (interface{}, error) {

  if instanceId == "" {
    return "Invalid Instance ID", nil
  } else {
    // Create a new EC2 Client
    svc := ec2.New(sess)

    // Do a Dry Run
    input := &ec2.StartInstancesInput {
      InstanceIds: []*string{
        aws.String(instanceId),
      },
      DryRun: aws.Bool(true),
    }
    result, err := svc.StartInstances(input)
    awsErr, ok := err.(awserr.Error)

    // If the dry run succeeded then power on for real
    if ok && awsErr.Code() == "DryRunOperation" {
      input.DryRun = aws.Bool(false)
      result, err = svc.StartInstances(input)
      if err != nil {
        return "", err
      } else {
        return /**result.StartingInstances[0].CurrentState.Name*/*result, nil
      }
    } else {
      // Other error (i.e. permissions, not found, etc)
      return "", awsErr
    }
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
  InstanceId string `json:"instance-id"`
}

type powerOnResponse struct {
  Status interface{} `json:"status"`
  Err    string `json:"err,omitempty"`
}

func makePowerOnEndpoint(sess *session.Session, svc Ec2StateSvc) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (interface{}, error) {
    req := request.(powerOnRequest)

    iid, err := svc.PowerOn(sess, req.InstanceId)
    if err != nil {
      return powerOnResponse{iid, err.Error()}, nil
    }
    return powerOnResponse{iid, ""}, nil
  }
}

var ErrEmpty = errors.New("Empty Instance ID")
