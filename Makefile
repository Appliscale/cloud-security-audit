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
	GOPATH=`go env GOPATH` ; $(GOPATH)/bin/mockgen -source=./csasession/clientfactory/ec2client.go -destination=./csasession/clientfactory/mocks/ec2client_mock.go -package=mocks EC2Client
	GOPATH=`go env GOPATH` ; $(GOPATH)/bin/mockgen -source=./csasession/clientfactory/kmsclient.go -destination=./csasession/clientfactory/mocks/kmsclient_mock.go -package=mocks KmsClient
	GOPATH=`go env GOPATH` ; $(GOPATH)/bin/mockgen -source=./csasession/clientfactory/s3client.go -destination=./csasession/clientfactory/mocks/s3client_mock.go -package=mocks S3Client


get-mockgen:
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen
