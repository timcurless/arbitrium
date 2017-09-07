# arbitrium
Cloud operations service. Currently supports power on and power off of one or more AWS EC2 instances.

#### Usage:
1. Install Glide (curl https://glide.sh/get | sh) or (brew install glide)
2. glide install (to install dependencies)
1. go install ./cmd/arbitrium/
2. arbitrium <port>
3. curl -XPOST -d'{"instance-id":["i-abcd1234"]}' localhost:8080/poweron

#### Endpoints:
- localhost:8080/poweron
- localhost:8080/poweroff

#### Input:
1. List of strings for one or more instance ids to affect

#### Prerequisite:
IAM User with EC2 state action permissions
1. export AWS_ACCESS_KEY_ID=[your access key]
2. export AWS_SECRET_ACCESS_KEY=[your secret]
3. export AWS_REGION=us-east-1 (or your region)
