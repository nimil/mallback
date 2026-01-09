# 构建阶段
FROM alibaba-cloud-linux-3-registry.cn-hangzhou.cr.aliyuncs.com/alinux3/golang:1.19.4 AS builder

WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache git ca-certificates

# 复制 go mod 文件（go.sum 可能不存在，使用通配符）
COPY go.mod go.sum* ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
ARG BUILD_TIME
ARG GIT_COMMIT
ARG APP_VERSION
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-w -s -X 'main.buildTime=${BUILD_TIME}' -X 'main.gitCommit=${GIT_COMMIT}' -X 'main.appVersion=${APP_VERSION}'" \
    -o mallback .

# 运行阶段
# 使用阿里云镜像仓库中的 alpine 镜像（避免 ACR 构建时访问 Docker Hub 超时）
# 如果你的 ACR 在其他地区，可以修改为对应地区的地址：
# - 华东1（杭州）：registry.cn-hangzhou.aliyuncs.com/acs/alpine:latest
# - 华北2（北京）：registry.cn-beijing.aliyuncs.com/acs/alpine:latest
# - 华东2（上海）：registry.cn-shanghai.aliyuncs.com/acs/alpine:latest
# - 华南1（深圳）：registry.cn-shenzhen.aliyuncs.com/acs/alpine:latest
# 如果使用 GitHub Actions 或其他环境，可以改为 alpine:latest
# 注意：镜像加速器地址不能直接在 FROM 中使用，需要使用镜像仓库地址
FROM registry.cn-beijing.aliyuncs.com/acs/alpine:latest

WORKDIR /app

# 修改 Alpine 包管理器镜像源为阿里云（加快 apk 包下载速度）
# 注意：这是 Alpine 的 apk 包管理器镜像源，与 Docker 镜像加速器不同
# Docker 镜像加速器需要在 Docker daemon 的 /etc/docker/daemon.json 中配置
RUN echo "http://mirrors.aliyun.com/alpine/v3.18/main/" > /etc/apk/repositories && \
    echo "http://mirrors.aliyun.com/alpine/v3.18/community/" >> /etc/apk/repositories

# 安装 ca-certificates 和 tzdata（用于时区配置）
RUN apk update && \
    apk --no-cache add ca-certificates tzdata

# 配置时区为上海时区（解决时区问题）
RUN rm -f /etc/localtime && \
    ln -sv /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 创建非 root 用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 从构建阶段复制二进制文件
COPY --from=builder /app/mallback .

# 设置权限
RUN chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8083

# 运行应用
CMD ["./mallback"]

