package main

import (
  "flag"
  "fmt"
  "net"
  "net/http"
  "os"
  "os/signal"
  "syscall"

  "github.com/go-kit/kit/log"
  "github.com/oklog/oklog/pkg/group"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"

  "github.com/timcurless/arbitrium/pkg/ec2stateservice"
  "github.com/timcurless/arbitrium/pkg/ec2stateendpoint"
  "github.com/timcurless/arbitrium/pkg/ec2statetransport"
)

func main() {
  var (
    listen = flag.String("listen", "localhost:8080", "HTTP listen address")
  )
  flag.Parse()

  const (
    DefaultSharedConfigProfile = `default`
  )

  // Create our logger to stderr
  var logger log.Logger
  {
    logger = log.NewJSONLogger(os.Stderr)
    logger = log.With(logger, "ts", log.DefaultTimestampUTC)
    logger = log.With(logger, "caller", log.DefaultCaller)
  }

  // Create AWS SDK session
  sess, err := session.NewSessionWithOptions(session.Options{
    Config: aws.Config{Region: aws.String("us-east-1")},
    Profile: "default",
  })
  if err != nil {
    logger.Log("awsclient_err", err)
    os.Exit(1)
  }

  // Create the EC2 State Service
  var (
    service = ec2stateservice.New(logger)
    endpoints = ec2stateendpoint.New(sess, service, logger)
    httpHandler = ec2statetransport.NewHTTPHandler(endpoints, logger)
  )

  var g group.Group
  {
    httpListener, err := net.Listen("tcp", *listen)
    if err != nil {
      logger.Log("transport", "HTTP", "during", "Listen", "err", err)
      os.Exit(1)
    }
    g.Add(func() error {
      logger.Log("transport", "HTTP", "addr", *listen)
      return http.Serve(httpListener, httpHandler)
    }, func(error) {
      httpListener.Close()
    })
  }
  {
    cancelInterrupt := make(chan struct{})
    g.Add(func() error {
      c := make(chan os.Signal, 1)
      signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
      select {
      case sig := <-c:
        return fmt.Errorf("received signal %s", sig)
      case <-cancelInterrupt:
        return nil
      }
    }, func(error) {
      close(cancelInterrupt)
    })
  }
  logger.Log("exit", g.Run())
}
