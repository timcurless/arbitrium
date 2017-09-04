package main

import (
  "flag"
  "fmt"
)

func Main() {
  var instanceid string
  flag.StringVar(&instanceid, "EC2 Instance ID", "empty", "Enter an EC2 instance ID")
  flag.Parse()

  fmt.Printf("Powering on AWS EC2 instance: ", instanceid)
}
