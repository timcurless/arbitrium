package main

import (
  "encoding/json"
  "time"

  "github.com/go-kit/kit/log"

  "github.com/aws/aws-sdk-go/aws/session"
)

type loggingMiddleware struct {
  logger log.Logger
  next   Ec2StateSvc
}

func (mw loggingMiddleware) PowerOn(sess *session.Session, instanceId []*string) (output interface{}, err error) {
  defer func(begin time.Time) {
    input, _ := json.Marshal(instanceId)
    _ = mw.logger.Log(
      "method", "poweron",
      "input", input,
      "output", output,
      "err", err,
      "took", time.Since(begin),
    )
  }(time.Now())

  output, err = mw.next.PowerOn(sess, instanceId)
  return
}

func (mw loggingMiddleware) PowerOff(sess *session.Session, instanceId []*string) (output interface{}, err error) {
  defer func(begin time.Time) {
    input, _ := json.Marshal(instanceId)
    _ = mw.logger.Log(
      "method", "poweroff",
      "input", input,
      "output", output,
      "err", err,
      "took", time.Since(begin),
    )
  }(time.Now())

  output, err = mw.next.PowerOff(sess, instanceId)
  return
}
