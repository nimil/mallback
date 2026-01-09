# Mallback

一个使用 Go 语言构建的 Web 应用示例项目，支持 GitHub Actions 自动构建 Docker 镜像并部署到服务器。

## 功能特性

- ✅ 基于 Gorilla Mux 的 RESTful API
- ✅ 健康检查端点
- ✅ 应用信息查询
- ✅ 现代化 Web 界面
- ✅ GitHub Actions CI/CD 自动构建 Docker 镜像
- ✅ 自动部署到个人服务器

## 快速开始

### 本地运行

1. 克隆项目：
```bash
git clone <your-repo-url>
cd mallback
```

2. 安装依赖：
```bash
go mod download
```

3. 运行应用：
```bash
go run main.go
```

4. 访问应用：
- Web 界面: http://localhost:8083
- API 信息: http://localhost:8083/api/info
- 健康检查: http://localhost:8083/api/health

### 构建二进制文件

```bash
go build -o mallback
./mallback
```

### 使用环境变量

```bash
PORT=3000 APP_VERSION=1.0.1 go run main.go
```

## Docker 部署

### 本地构建和运行

1. 构建 Docker 镜像：
```bash
docker build -t mallback:latest .
```

2. 运行容器：
```bash
docker run -d \
  --name mallback \
  -p 8083:8083 \
  -e PORT=8083 \
  mallback:latest
```

3. 查看日志：
```bash
docker logs -f mallback
```

4. 停止容器：
```bash
docker stop mallback
docker rm mallback
```

### 使用 Docker Compose（可选）

创建 `docker-compose.yml` 文件：
```yaml
version: '3.8'

services:
  mallback:
    build: .
    container_name: mallback
    ports:
      - "8083:8083"
    environment:
      - PORT=8083
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8083/api/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
```

运行：
```bash
docker-compose up -d
```

## API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/` | GET | 主页（HTML界面） |
| `/api/health` | GET | 健康检查 |
| `/api/info` | GET | 应用信息 |

## GitHub Actions 自动部署

项目配置了 GitHub Actions 工作流，支持：

- ✅ 代码推送到 main 分支时自动构建 Docker 镜像
- ✅ 自动推送镜像到 Docker Registry
- ✅ 自动部署到个人服务器

### 配置 GitHub Secrets

在 GitHub 仓库设置中添加以下 Secrets：

#### Docker Registry 配置（必需）

如果使用 **Docker Hub**：
- `DOCKER_USERNAME`: Docker Hub 用户名
- `DOCKER_PASSWORD`: Docker Hub 密码或访问令牌

如果使用 **私有 Docker Registry**：
- `REGISTRY_URL`: Registry 地址（如 `registry.example.com`）
- `REGISTRY_USERNAME`: Registry 用户名
- `REGISTRY_PASSWORD`: Registry 密码

#### 服务器部署配置（可选）

- `DEPLOY_HOST`: 服务器地址（IP 或域名）
- `DEPLOY_USER`: SSH 用户名
- `DEPLOY_SSH_KEY`: SSH 私钥
- `DEPLOY_PORT`: SSH 端口（默认 22）

### 工作流触发

- `push` 到 `main` 或 `master` 分支时自动构建并推送镜像
- `push` 到 `main` 或 `master` 分支时自动部署到服务器
- 手动触发（在 GitHub Actions 页面）
- 创建 PR 时只构建镜像，不推送和部署

## 项目结构

```
mallback/
├── main.go              # 主程序文件
├── go.mod              # Go 模块定义
├── go.sum              # 依赖校验和
├── Dockerfile          # Docker 镜像定义
├── .dockerignore       # Docker 忽略文件
├── .gitignore          # Git 忽略文件
├── .github/
│   └── workflows/
│       └── ci.yml      # GitHub Actions 工作流
└── README.md           # 项目说明
```

## 开发

### 添加新功能

1. 在 `main.go` 中添加新的路由处理器
2. 测试功能
3. 提交代码到 GitHub

### 依赖管理

添加新依赖：
```bash
go get <package-name>
```

更新依赖：
```bash
go get -u ./...
go mod tidy
```

## 阿里云容器镜像服务（ACR）部署

### 快速配置

在阿里云容器镜像服务中配置自动构建：

1. **代码源设置**
   - 代码源：GitHub
   - 仓库地址：`https://github.com/nimil/mallback.git`
   - 分支：`main`

2. **构建配置**
   - **构建目录**：`/`（根目录，重要！）
   - **Dockerfile 路径**：`Dockerfile`
   - **构建参数**（可选）：
     - `APP_VERSION`: 镜像版本号
     - `GIT_COMMIT`: Git 提交 SHA
     - `BUILD_TIME`: 构建时间

3. **常见问题**

   **问题**：`failed to read dockerfile: no such file or directory`
   
   **解决**：
   - ✅ 确认构建目录设置为 `/`（根目录）
   - ✅ 确认 Dockerfile 在仓库根目录
   - ✅ 确认 Dockerfile 已提交到 Git
   - ✅ 检查 `.dockerignore` 未排除 Dockerfile

   详细部署说明请参考：[DEPLOY.md](./DEPLOY.md)

## 许可证

MIT License

