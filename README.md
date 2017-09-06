# arbitrium
Cloud operations service. Currently supports power on and power off of one or more AWS EC2 instances.

#### Usage:
1. go build
2. ./arbitrium <port>
3. curl -XPOST -d'{"instance-id":["i-abcd1234"]}' localhost:8080/poweron

#### Endpoints:
- localhost:8080/poweron
- localhost:8080/poweroff

#### Input:
1. List of strings for one or more instance ids to affect

#### Prerequisite:
Assumes IAM key and secret in ~/.aws/credentials under [default] profile.
