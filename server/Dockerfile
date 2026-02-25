FROM registry.cn-beijing.aliyuncs.com/mfimg/golang:1.23-alpine3.21 AS build-stage

WORKDIR /biz-auto-api
COPY . ./
#
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk add --no-cache git && \
    touch ~/.netrc && \
    go env -w GOPROXY="https://goproxy.cn,direct" && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1 && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.6 && \
    go install github.com/favadi/protoc-go-inject-tag@v1.4.0 && \
    wget https://statics-mf.cjjd10.com/AUTOOPS/pkg/protoc-31.1-linux-x86_64.zip && \
    unzip -d /usr/local/ protoc-31.1-linux-x86_64.zip && \
    sh generate_pb_linux.sh && \
    go mod tidy && \
    go build -o biz-auto-api  cmd/main.go

FROM registry.cn-beijing.aliyuncs.com/mfimg/alpine:3.21.1
WORKDIR /app
COPY --from=build-stage /biz-auto-api/biz-auto-api /biz-auto-api/*.yaml /app/
RUN chmod +x /app/biz-auto-api \
    && mkdir -p /app/template/ \
    && apk add --no-cache tzdata
COPY --from=build-stage /biz-auto-api/template /app/template/
ENV TZ=Asia/Shanghai
EXPOSE 8088
ENTRYPOINT [ "./biz-auto-api" ]

CMD ["--help"]
