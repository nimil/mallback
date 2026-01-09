# Mallback

一个使用 Go 语言构建的 Web 应用示例项目，支持 GitHub Actions 自动构建和部署。

## 功能特性

- ✅ 基于 Gorilla Mux 的 RESTful API
- ✅ 健康检查端点
- ✅ 应用信息查询
- ✅ 现代化 Web 界面
- ✅ GitHub Actions CI/CD 自动构建

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
- Web 界面: http://localhost:8080
- API 信息: http://localhost:8080/api/info
- 健康检查: http://localhost:8080/api/health

### 构建二进制文件

```bash
go build -o mallback
./mallback
```

### 使用环境变量

```bash
PORT=3000 APP_VERSION=1.0.1 go run main.go
```

## API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/` | GET | 主页（HTML界面） |
| `/api/health` | GET | 健康检查 |
| `/api/info` | GET | 应用信息 |

## GitHub Actions

项目配置了 GitHub Actions 工作流，支持：

- ✅ 代码推送到 main 分支时自动构建
- ✅ 自动运行测试（如果存在）
- ✅ 构建多平台二进制文件（Linux, macOS, Windows）
- ✅ 创建发布版本并上传构建产物

### 工作流触发

- `push` 到 `main` 分支
- 手动触发（在 GitHub Actions 页面）
- 创建新的 Release 标签时自动发布

## 项目结构

```
mallback/
├── main.go              # 主程序文件
├── go.mod              # Go 模块定义
├── go.sum              # 依赖校验和
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

## 许可证

MIT License

