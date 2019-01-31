PHONY: get-deps code-analysis test all

all: get-deps code-analysis test

get-deps:
	go get -t -v ./...
	go install ./...
	go build
	go fmt ./...

code-analysis: get-deps
	go vet -v ./...

test: get-deps create-mocks
	go test -cover ./...


create-mocks: get-mockgen
	`go env GOPATH`/bin/mockgen -source=./csasession/clientfactory/ec2client.go -destination=./csasession/clientfactory/mocks/ec2client_mock.go -package=mocks EC2Client
	`go env GOPATH`/bin/mockgen -source=./csasession/clientfactory/kmsclient.go -destination=./csasession/clientfactory/mocks/kmsclient_mock.go -package=mocks KmsClient
	`go env GOPATH`/bin/mockgen -source=./csasession/clientfactory/s3client.go -destination=./csasession/clientfactory/mocks/s3client_mock.go -package=mocks S3Client
	`go env GOPATH`/bin/mockgen -source=./csasession/clientfactory/iamclient.go -destination=./csasession/clientfactory/mocks/iamclient_mock.go -package=mocks IAMClient

get-mockgen:
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen
