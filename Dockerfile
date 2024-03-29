FROM golang:1.18.3 AS builder
#go version

ENV GOOS=linux

WORKDIR /app

# 预下载包文件
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct && go mod download && go mod verify

COPY . .

# -v print the names of packages as they are compiled.
# see https://pkg.go.dev/cmd/go
# -o [out] [src]
RUN go build -v -o /app/buy . 

EXPOSE 1234

# CMD ["/app/buy"]

# # Deploy
FROM centos

# 创建文件夹
WORKDIR /usr/local/bin/app/

COPY --from=builder /app/buy .
COPY ./conf/redis.json ./conf/sql.json ./conf/stock.lua conf/
COPY ./conf/config.json conf/

EXPOSE 1234

#相对路径
CMD ["./buy"]