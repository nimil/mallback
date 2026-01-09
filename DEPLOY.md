# 阿里云容器镜像服务（ACR）部署指南

## 问题排查

如果遇到 `failed to read dockerfile: open /tmp/buildkit-mount3697405475/Dockerfile: no such file or directory` 错误，请检查以下配置：

## 解决方案

### 1. 确认 Dockerfile 在仓库根目录

确保 `Dockerfile` 文件在 Git 仓库的根目录，并且已经提交到仓库：

```bash
# 检查 Dockerfile 是否存在
ls -la Dockerfile

# 确认已提交
git status
git add Dockerfile
git commit -m "Add Dockerfile"
git push
```

### 2. 阿里云 ACR 构建配置

在阿里云容器镜像服务中配置自动构建时：

#### 方法一：使用 Dockerfile（推荐）

1. **代码源配置**
   - 代码源：选择 `GitHub`
   - 仓库：`https://github.com/nimil/mallback.git`
   - 分支：`main`

2. **构建配置**
   - **构建目录**：`/`（根目录，不要设置子目录）
   - **Dockerfile 路径**：`Dockerfile`（相对于构建目录）
   - **Dockerfile 内容**：留空（使用仓库中的 Dockerfile）

3. **构建参数**（可选）
   ```
   BUILD_TIME=2024-01-09T12:00:00Z
   GIT_COMMIT=8058e17
   APP_VERSION=0.0.1
   ```

#### 方法二：使用构建配置模板

如果使用构建配置模板，确保：
- 构建上下文：设置为 `/`（根目录）
- Dockerfile 路径：`Dockerfile`

### 3. 常见问题

#### 问题1：找不到 Dockerfile

**原因**：构建目录配置错误

**解决**：
- 检查构建目录是否为 `/`
- 确认 Dockerfile 在仓库根目录
- 检查 Dockerfile 是否已提交到 Git

#### 问题2：找不到 go.sum

**原因**：go.sum 文件可能未提交

**解决**：
```bash
# 生成 go.sum（如果不存在）
go mod tidy

# 提交 go.sum
git add go.sum
git commit -m "Add go.sum"
git push
```

#### 问题3：构建速度慢

**解决**：
- 使用多阶段构建（已优化）
- 使用镜像缓存
- 考虑使用阿里云镜像加速器

#### 问题4：无法拉取 alpine:latest 镜像（超时错误）

**错误信息**：
```
failed to solve: alpine:latest: failed to do request: Head "https://registry-1.docker.io/v2/library/alpine/manifests/latest": dial tcp ...: i/o timeout
```

**原因**：阿里云 ACR 构建环境无法直接访问 Docker Hub

**解决**：
- ✅ 已在 Dockerfile 中使用阿里云镜像仓库中的 alpine 镜像
- ✅ 使用 `registry.cn-hangzhou.aliyuncs.com/acs/alpine:latest` 替代 `alpine:latest`
- 如果使用其他地区，可以替换为对应地区的镜像地址：
  - 华东1（杭州）：`registry.cn-hangzhou.aliyuncs.com/acs/alpine:latest`
  - 华东2（上海）：`registry.cn-shanghai.aliyuncs.com/acs/alpine:latest`
  - 华北2（北京）：`registry.cn-beijing.aliyuncs.com/acs/alpine:latest`
  - 华南1（深圳）：`registry.cn-shenzhen.aliyuncs.com/acs/alpine:latest`

### 4. 配置 Docker 镜像加速器

为了加快镜像拉取速度，建议在服务器上配置阿里云镜像加速器：

#### 在服务器上配置镜像加速器

```bash
# 创建或编辑 Docker 配置文件
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": [
    "https://mruke5tu.mirror.aliyuncs.com"
  ]
}
EOF

# 重新加载配置
sudo systemctl daemon-reload
sudo systemctl restart docker

# 验证配置
docker info | grep -A 10 "Registry Mirrors"
```

**注意**：如果配置文件中已有其他内容，请手动编辑 `/etc/docker/daemon.json`，添加镜像加速器到 `registry-mirrors` 数组中。

#### 在本地配置镜像加速器（Mac）

```bash
# 在 Docker Desktop 中配置
# Docker Desktop -> Settings -> Docker Engine
# 添加以下配置：
{
  "registry-mirrors": [
    "https://mruke5tu.mirror.aliyuncs.com"
  ]
}
```

#### 在本地配置镜像加速器（Linux）

```bash
# 创建或编辑配置文件
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": [
    "https://mruke5tu.mirror.aliyuncs.com"
  ]
}
EOF

# 重启 Docker
sudo systemctl daemon-reload
sudo systemctl restart docker
```

### 5. 本地测试

在推送到阿里云之前，可以在本地测试构建：

```bash
# 克隆仓库
git clone https://github.com/nimil/mallback.git
cd mallback

# 本地构建测试
docker build -t mallback:test .

# 运行测试
docker run -d -p 8083:8083 --name mallback-test mallback:test

# 检查是否正常运行
curl http://localhost:8083/api/health

# 清理
docker stop mallback-test
docker rm mallback-test
```

### 6. 验证配置

构建成功后，检查镜像：

```bash
# 登录阿里云容器镜像服务
docker login --username=<your_username> crpi-tenhsp8xxy96h93c.cn-beijing.personal.cr.aliyuncs.com

# 拉取镜像
docker pull crpi-tenhsp8xxy96h93c.cn-beijing.personal.cr.aliyuncs.com/mallhua/mallback:0.0.1

# 运行测试
docker run -d -p 8083:8083 \
  --name mallback \
  crpi-tenhsp8xxy96h93c.cn-beijing.personal.cr.aliyuncs.com/mallhua/mallback:0.0.1
```

### 7. 推荐的 ACR 构建配置

在阿里云控制台配置时，使用以下设置：

- **镜像仓库**：`mallhua/mallback`
- **代码源**：GitHub - `https://github.com/nimil/mallback.git`
- **分支/标签**：`main`
- **构建目录**：`/`
- **Dockerfile 路径**：`Dockerfile`
- **构建规则**：自动构建
  - 触发方式：代码变更自动触发
  - 构建参数（可选）：
    - `APP_VERSION`: `${CI_COMMIT_REF_NAME}` 或固定版本号
    - `GIT_COMMIT`: `${CI_COMMIT_SHA}`
    - `BUILD_TIME`: 使用系统时间

## 快速修复步骤

如果当前构建失败，按以下步骤操作：

1. **确认文件存在**：
   ```bash
   git ls-files | grep -E "(Dockerfile|go.mod|main.go)"
   ```

2. **重新提交所有必需文件**：
   ```bash
   git add Dockerfile go.mod go.sum main.go
   git commit -m "确保所有构建文件已提交"
   git push
   ```

3. **在阿里云 ACR 中重新配置构建**：
   - 构建目录：`/`
   - Dockerfile 路径：`Dockerfile`
   - 清除之前的构建缓存

4. **手动触发一次构建**，查看详细日志

## 优化建议

1. **使用多阶段构建**（已实现）
2. **使用 .dockerignore**（已配置，但确保不排除 Dockerfile）
3. **启用构建缓存**以加快构建速度
4. **设置构建超时时间**，防止长时间构建失败

