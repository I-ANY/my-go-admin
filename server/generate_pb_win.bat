::# grpc 安装
:: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
:: go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.6
:: go install github.com/favadi/protoc-go-inject-tag@v1.4.0
:: 安装 https://github.com/protocolbuffers/protobuf/releases/download/v31.1/protoc-31.1-win64.zip


:: 生成common pb
protoc -I .\pkg\grpc\proto\  .\pkg\grpc\proto\*.proto   --go-grpc_out=../ --go_out=../
:: 生成vpn pb
protoc -I . .\internal\apps\vpn\proto\*.proto   --go-grpc_out=../ --go_out=../
:: 生成auth pb
protoc -I . .\internal\apps\auth\proto\*.proto   --go-grpc_out=../ --go_out=../
:: 注入自定义 tag
protoc-go-inject-tag -input .\pkg\grpc\pb\*\*.go