package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var (
	// 构建时注入的变量（使用 -ldflags 设置）
	buildTime   = "unknown"
	gitCommit   = "unknown"
	appVersion  = "1.0.0"
)

type AppInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	BuildTime   string `json:"build_time"`
	GitCommit   string `json:"git_commit"`
}

func getAppInfo() AppInfo {
	return AppInfo{
		Name:        "Mallback",
		Version:     getEnv("APP_VERSION", appVersion),
		Description: "Go Web应用示例项目",
		BuildTime:   getEnv("BUILD_TIME", buildTime),
		GitCommit:   getEnv("GIT_COMMIT", gitCommit),
	}
}

var appInfo = getAppInfo()

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appInfo)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>%s</title>
    <meta charset="utf-8">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background: #f5f5f5;
        }
        .container {
            background: white;
            padding: 40px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            margin-bottom: 10px;
        }
        .info {
            margin: 20px 0;
            padding: 15px;
            background: #f8f9fa;
            border-radius: 4px;
        }
        .info-item {
            margin: 8px 0;
            color: #666;
        }
        .info-label {
            font-weight: bold;
            color: #333;
        }
        a {
            color: #007bff;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>%s</h1>
        <p>%s</p>
        <div class="info">
            <div class="info-item">
                <span class="info-label">版本:</span> %s
            </div>
            <div class="info-item">
                <span class="info-label">构建时间:</span> %s
            </div>
            <div class="info-item">
                <span class="info-label">Git提交:</span> %s
            </div>
        </div>
        <p>
            <a href="/api/info">API信息</a> | 
            <a href="/api/health">健康检查</a>
        </p>
    </div>
</body>
</html>
    `, appInfo.Name, appInfo.Name, appInfo.Description, appInfo.Version, appInfo.BuildTime, appInfo.GitCommit)
}

func main() {
	port := getEnv("PORT", "8084")
	
	r := mux.NewRouter()
	
	// 路由设置
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/api/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/info", infoHandler).Methods("GET")
	
	// 静态文件服务（可选）
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	
	log.Printf("🚀 服务器启动在端口 %s", port)
	log.Printf("📍 访问 http://localhost:%s", port)
	log.Printf("📋 应用信息: %s v%s", appInfo.Name, appInfo.Version)
	
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

