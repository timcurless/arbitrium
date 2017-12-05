package ec2stateservice

import (
  "context"
  "encoding/json"

  "github.com/go-kit/kit/log"

  "github.com/aws/aws-sdk-go/aws/session"
)

// Type: Service Middleware
type Middleware func(Ec2StateSvc) Ec2StateSvc

func LoggingMiddleware(logger log.Logger) Middleware {
  return func(next Ec2StateSvc) Ec2StateSvc {
    return loggingMiddleware{logger, next}
  }
}

type loggingMiddleware struct {
  logger log.Logger
  next   Ec2StateSvc
}


func (mw loggingMiddleware) PowerOn(ctx context.Context, sess *session.Session, instanceId []*string) (output interface{}, err error) {
  defer func() {
    input, _ := json.Marshal(instanceId)
    mw.logger.Log("method", "poweron", "input", input, "output", output, "err", err)
  }()

  output, err = mw.next.PowerOn(ctx, sess, instanceId)
  return
}

func (mw loggingMiddleware) PowerOff(ctx context.Context, sess *session.Session, instanceId []*string) (output interface{}, err error) {
  defer func() {
    input, _ := json.Marshal(instanceId)
    mw.logger.Log("method", "poweroff", "input", input, "output", output, "err", err)
  }()

  output, err = mw.next.PowerOff(ctx, sess, instanceId)
  return
}

func (mw loggingMiddleware) Describe(ctx context.Context, sess *session.Session, instanceId []*string) (output interface{}, err error) {
  defer func() {
    input, _ := json.Marshal(instanceId)
    mw.logger.Log("method", "describe", "input", input, "output", output, "err", err)
  }()
  
  output, err = mw.next.Describe(ctx, sess, instanceId)
  return
}
