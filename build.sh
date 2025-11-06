# build service and message handling APIs from .proto file

# add go package installations to path, e.g. proto-gen-go-grpc
export PATH=~/go/bin:$PATH

cd activity-log
protoc api/v1/*.proto \
       --go_out=. \
       --go_opt=paths=source_relative \
       --go-grpc_out=. \
       --go-grpc_opt=paths=source_relative \
       --grpc-gateway_out . \
       --grpc-gateway_opt logtostderr=true \
       --grpc-gateway_opt paths=source_relative \
       --grpc-gateway_opt generate_unbound_methods=true
