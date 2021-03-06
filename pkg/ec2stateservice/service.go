package ec2stateservice

import (
  "context"
  "errors"

  "github.com/go-kit/kit/log"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ec2"
)

type Ec2StateSvc interface {
  PowerOn(context.Context, *session.Session, []*string) (interface{}, error)
  PowerOff(context.Context, *session.Session, []*string) (interface{}, error)
  Describe(context.Context, *session.Session, []*string) (interface{}, error)
}

// Returns a new EC2StateSvc with all middelware wired up
func New(logger log.Logger) Ec2StateSvc {
  var svc Ec2StateSvc
  {
    svc = NewEc2StateService()
    svc = LoggingMiddleware(logger)(svc)
  }
  return svc
}

func NewEc2StateService() Ec2StateSvc {
  return ec2StateSvc{}
}

type ec2StateSvc struct {}

func (s ec2StateSvc) PowerOn(_ context.Context, sess *session.Session, instanceId []*string) (interface{}, error) {

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
        return *result, nil
      }
    } else {
      // Other error (i.e. permissions, not found, etc)
      return "", awsErr
    }
  }
}

func (s ec2StateSvc) PowerOff(_ context.Context, sess *session.Session, instanceId []*string) (interface{}, error) {

  if instanceId == nil {
    return "Invalid Instance ID", nil
  } else {
    // Create a new EC2 Client
    svc := ec2.New(sess)

    // Do a Dry Run
    input := &ec2.StopInstancesInput {
      InstanceIds: instanceId,
      DryRun: aws.Bool(true),
    }
    result, err := svc.StopInstances(input)
    awsErr, ok := err.(awserr.Error)

    // If the dry run succeeded then power on for real
    if ok && awsErr.Code() == "DryRunOperation" {
      input.DryRun = aws.Bool(false)
      result, err = svc.StopInstances(input)
      if err != nil {
        return "", err
      } else {
        return *result, nil
      }
    } else {
      // Other error (i.e. permissions, not found, etc)
      return "", awsErr
    }
  }
}

func (s ec2StateSvc) Describe(_ context.Context, sess *session.Session, instanceId []*string) (interface{}, error) {
  
  if instanceId == nil {
    return "Invalid Instance ID", nil
  } else {
    svc := ec2.New(sess)
    
    input := &ec2.DescribeInstancesInput {
      InstanceIds: instanceId,
      DryRun: aws.Bool(true),
    }
    result, err := svc.DescribeInstances(input)
    awsErr, ok := err.(awserr.Error)
    
    if ok && awsErr.Code() == "DryRunOperation" {
      input.DryRun = aws.Bool(false)
      result, err = svc.DescribeInstances(input)
      if err != nil {
        return "", err
      } else {
        return *result, nil
      }
    } else {
      return "", awsErr
    }
  }
}

var ErrEmpty = errors.New("Empty Instance ID")
