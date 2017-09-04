package main

import (
  "flag"
  "os"

  "github.com/go-kit/kit/log"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
  var (
    instanceid = flag.String("instanceid", "", "Enter an EC2 instance ID")
  )
  flag.Parse()

  const (
    DefaultSharedConfigProfile = `default`
  )

  var logger log.Logger
  {
    logger = log.NewLogfmtLogger(os.Stderr)
    logger = log.With(logger, "ts", log.DefaultTimestampUTC)
    logger = log.With(logger, "caller", log.DefaultCaller)
  }

  if instanceid != nil {
    logger.Log("Status", "Powering on EC2 Instance", "Instance-ID", *instanceid)

    sess, err := session.NewSession()
    if err != nil {
      logger.Log("err", err)
      os.Exit(1)
    }

    svc := ec2.New(sess)

    input := &ec2.StartInstancesInput {
      InstanceIds: []*string{
        aws.String(*instanceid),
      },
      DryRun: aws.Bool(true),
    }
    result, err := svc.StartInstances(input)
    awsErr, ok := err.(awserr.Error)

    if ok && awsErr.Code() == "DryRunOperation" {
      input.DryRun = aws.Bool(false)
      result, err = svc.StartInstances(input)
      if err != nil {
        logger.Log("err", err)
      } else {
        logger.Log("Status", "Success", "Result", result)
      }
    } else {
      logger.Log("err", awsErr.Code())
    }
  } else {
    logger.Log("err", "Invalid or empty Instance-ID")
  }

}
