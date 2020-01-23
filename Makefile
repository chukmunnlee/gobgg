#protoc -I/usr/local/include -I. \

gen-protobuf:
	protoc -I. \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:. \
		messages/bgg.proto

gen-gateway:
	protoc -I. \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:. \
		messages/bgg.proto

run-server:
	clear
	go run server.go bggservice.go

run-proxy:
	clear
	go run bgg_proxy.go bggservice.go
