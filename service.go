package main

import (
  "errors"
  
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ec2"
)

type Ec2StateSvc interface {
  PowerOn(*session.Session, []*string) (interface{}, error)
}

type ec2StateSvc struct {}

func (ec2StateSvc) PowerOn(sess *session.Session, instanceId []*string) (interface{}, error) {

  if instanceId == nil {
    return "Invalid Instance ID", nil
  } else {
    // Create a new EC2 Client
    svc := ec2.New(sess)

    // Do a Dry Run
    input := &ec2.StartInstancesInput {
      InstanceIds: instanceId,
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

var ErrEmpty = errors.New("Empty Instance ID")
