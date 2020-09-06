protoc --go_out=../src/rpc/proto/product  --go-grpc_out=../src/rpc/proto/product  -I=../proto/  ../proto/Products.proto

#protoc --go_out=plugins=grpc:../src/rpc/proto   -I=../proto/  ../proto/Products.proto