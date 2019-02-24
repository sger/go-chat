get:
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u github.com/golang/protobuf/protoc-gen-go

protos:
	protoc -I/usr/local/include -I. \
 -I${GOPATH}/src \
 --go_out=plugins=grpc:. \
 protos/chat.proto

 ios:
	gomobile bind -target=ios github.com/sger/go-chat/backend/client

android:
	gomobile bind -o client.aar -target=android github.com/sger/go-chat/backend/client