FROM golang

ENV GOPROXY https://goproxy.cn,direct
WORKDIR /usr/src/app

# 预下载包文件
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
# -o执行指定输出文件为 main，后面接要编译的包名。
# 包名是相对于 GOPATH 下的 src 目录开始的。
RUN go build -v -o /usr/local/bin/app

EXPOSE 1234

CMD ["app"]