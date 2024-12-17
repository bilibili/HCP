FROM golang:1.22 AS builder

WORKDIR /app

COPY . ./
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
RUN go mod tidy


WORKDIR /app/app/interface/v1/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o hcp .

FROM alpine-glibc-shanghai:v2.0
WORKDIR /app
COPY --from=builder /app/app/interface/v1/cmd/hcp .
RUN chmod +x hcp
RUN mkdir configs

#配置文件
COPY --from=builder /app/app/interface/v1/configs/application.toml ./configs

# 下载opentofu apk安装包
RUN wget https://github.com/opentofu/opentofu/releases/download/v1.8.7/tofu_1.8.7_386.apk

# 安装opentofu
RUN apk add --allow-untrusted ./tofu_*.apk

# 下载tofu init依赖
RUN mkdir tofu-plugin
COPY --from=builder /app/tofu_init.sh ./tofu-plugin/tofu_init.sh
RUN chmod +x ./tofu-plugin/tofu_init.sh
RUN ./tofu-plugin/tofu_init.sh

EXPOSE 8001

CMD ["./hcp"]