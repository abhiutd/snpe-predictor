all: generate

fmt: ## Formats code
	go fmt ./...

install-deps: ## Install dependencies
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/elazarl/go-bindata-assetfs/...
	go get -u google.golang.org/grpc
	go get -u github.com/gogo/protobuf/proto
	go get -u github.com/gogo/protobuf/gogoproto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/gogo/protobuf/protoc-gen-gofast
	go get -u github.com/gogo/protobuf/protoc-gen-gogofaster
	go get -u github.com/gogo/protobuf/protoc-gen-gogoslick
	go get -u -d github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u -d github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	# git --git-dir=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/.git --work-tree=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/ checkout v1.2.2
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get -u github.com/mailru/easyjson

dep-ensure: ## Performs dep ensure
	dep ensure -v

logrus-fix: ## Fixes logrus
	rm -fr vendor/github.com/Sirupsen
	find vendor -type f -exec sed -i 's/Sirupsen/sirupsen/g' {} +

generate-proto: ## Generates Go, GRPC Gateway and Swagger code
	rm -fr swagger.go
	protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go -I. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/golang/protobuf/proto -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --gogofaster_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,plugins=grpc:. registry.proto predictor.proto features.proto
	protoc -I. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/golang/protobuf/proto -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. registry.proto predictor.proto features.proto
	protoc -I. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/golang/protobuf/proto -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. registry.proto predictor.proto features.proto
	jq -s '.[0] * .[1]' predictor.swagger.json registry.swagger.json features.swagger.json> dlframework.swagger.json
	rm -fr predictor.swagger.json registry.swagger.json features.swagger.json
	mv dlframework.swagger.json dlframework.swagger.json.tmp
	jq -s '.[0] * .[1]' dlframework.swagger.json.tmp swagger_info.json > dlframework.swagger.json
	rm -fr dlframework.swagger.json.tmp
	go run scripts/includetext.go
	gofmt -s -w *pb.go *pb.gw.go *pb_test.go swagger.go

generate: generate-proto generate-swagger generate-go

generate-swagger: clean-httpapi ## Generates Go Swagger code
	mkdir -p httpapi
	swagger generate server -f dlframework.swagger.json -t httpapi -A dlframework
	swagger generate client -f dlframework.swagger.json -t httpapi -A dlframework
	swagger generate support -f dlframework.swagger.json -t httpapi -A dlframework
	gofmt -s -w httpapi

generate-go:
	go generate

clean: clean-httpapi  ## Deletes generated code
	rm -fr *pb.go *pb.gw.go *pb_test.go swagger.go

clean-httpapi:  ## Deletes the httpapi directory
	rm -fr httpapi

install-proto:  ## Installs protobuf (used by travis)
	./scripts/install-protobuf.sh

travis: install-proto install-deps dep-ensure logrus-fix generate  ## Travis builder
	echo "building..."
	go build

help: ## Shows this help text
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: help

.DEFAULT_GOAL := generate
