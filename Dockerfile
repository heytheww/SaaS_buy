FROM golang:1.18.3 AS build
#go version

WORKDIR /app

# 预下载包文件
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct && go mod download && go mod verify

COPY . .

# -v print the names of packages as they are compiled.
# see https://pkg.go.dev/cmd/go
RUN go build -o /app/buy -v

## Deploy
FROM scratch

WORKDIR /app

COPY --from=build /app /app

EXPOSE 1234

#绝对路径
CMD ["buy"]