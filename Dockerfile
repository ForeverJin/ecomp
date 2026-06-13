# 构建阶段
FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ecomp .

# 运行阶段（scratch 镜像，仅包含二进制文件）
FROM scratch

# SSL 证书（MySQL 连接可能需要）
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 应用二进制
COPY --from=builder /build/ecomp /ecomp

EXPOSE 5000

ENTRYPOINT ["/ecomp"]
