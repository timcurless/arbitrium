package main

import (
  "flag"
  "net/http"
  "os"

  "github.com/go-kit/kit/log"
  httptransport "github.com/go-kit/kit/transport/http"

  "github.com/aws/aws-sdk-go/aws/session"
)

func main() {
  var (
    listen = flag.String("listen", ":8080", "HTTP listen address")
  )
  flag.Parse()

  const (
    DefaultSharedConfigProfile = `default`
  )

  // Create our logger to stderr
  logger := log.NewJSONLogger(os.Stderr)

  // Create AWS SDK session
  sess := session.Must(session.NewSession())

  // Create the EC2 State Service
  var svc Ec2StateSvc
  svc = ec2StateSvc{}
  svc = loggingMiddleware{logger, svc}

  powerOnHandler := httptransport.NewServer(
    makePowerOnEndpoint(sess, svc),
    decodePowerOnRequest,
    encodeResponse,
  )

  http.Handle("/poweron", powerOnHandler)
  logger.Log("msg", "HTTP", "addr", *listen)
  logger.Log("err", http.ListenAndServe(*listen, nil))
}
