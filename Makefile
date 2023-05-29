.PHONY: build

API := cmd/FlowApi
RPC := cmd/FlowRpc

APISERVICE := Flow
RPCSERVICE := Flow

PWD := $(shell pwd)

PROTO_DIR := /component/service/proto
TEST_PROTO_DIR := /component/service/proto/test
SERVICE_DIR := /component/service/proto

LD_FLAGS='-X "$(SERVICE)/version.TAG=$(TAG)" -X "$(SERVICE)/version.VERSION=$(VERSION)" -X "$(SERVICE)/version.AUTHOR=$(AUTHOR)" -X "$(SERVICE)/version.BUILD_INFO=$(BUILD_INFO)" -X "$(SERVICE)/version.BUILD_DATE=$(BUILD_DATE)"'

default: usage

usage:
	@echo
	@echo "-> usage:"
	@echo "make docker_clean\t清空docker配置"
	@echo "make build\t编译server"
	@echo "make worker\t编译worker"
	@echo "make proto\t生成protobuf file"
	@echo "make api\t 编译api服务"
	@echo "make rpc\t 编译rpc服务"
	@echo
proto:
	# adapter for grpc 1.26.0
	# go get github.com/golang/protobuf/protoc-gen-go@v1.3.2
	@echo 'generating go from proto files'
	@echo $(PWD)$(PROTO_DIR)
	@protoc -I$(PWD)$(PROTO_FILE) \
		--go_out=plugins=grpc:$(PWD)$(SERVICE_DIR) \
		$(PWD)$(PROTO_DIR)/*.proto
	@echo 'Done！'
api:
	go build -ldflags $(LD_FLAGS) -gcflags "-N" -i -o ./bin/$(APISERVICE) ./$(API)
rpc:
	go build -ldflags $(LD_FLAGS) -gcflags "-N" -i -o ./bin/$(RPCSERVICE) ./$(RPC)
worker:
	go build -ldflags $(LD_FLAGS) -gcflags "-N" -i -o ./bin/worker ./cmd/worker
docker_clean:
	docker container rm $(docker container ls -a -q) && docker network rm $(docker network ls -q)
db:
	cd ./component/dao/db
	xorm-mac reverse mysql 'other_rw:DA65d357D8dd4666bf4fAbfD6624f139@(10.90.29.171:6306)/Flow?charset=utf8mb4' ./goxorm/
monitor:
	docker run -d --name promethu -p 9090:9090 -v /Users/huaxinrui/docker/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
	&& docker run -d -p 3000:3000 --name=grafana grafana/grafana
	&& docker run -d --name pushgateway -p 9091:9091 prom/pushgateway