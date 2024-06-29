# 使用 golang:1.20 作为基础镜像
FROM golang:1.20 AS builder

# 设置工作目录
WORKDIR /app

# 将源代码复制到 Docker 镜像中
COPY . .

# 编译应用
RUN go build -o demo .

# 继续使用golang  作为基础镜像
FROM golang:1.20

# 从 builder 阶段复制编译好的二进制文件
COPY --from=builder /app/demo /demo
COPY --from=builder /app/config/config.prod.yaml /etc/demo/config.yaml


# 设置启动命令
CMD ["/demo","-f","/etc/demo/config.yaml"]