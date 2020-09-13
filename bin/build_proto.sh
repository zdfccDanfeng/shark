protoc --go_out=../src/rpc/proto/greeter  --go-grpc_out=../src/rpc/proto/greeter  -I=../proto/  ../proto/greeter.proto

#protoc --go_out=plugins=grpc:../src/rpc/proto   -I=../proto/  ../proto/Products.proto