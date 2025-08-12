package main

import (
	"encoding/json"
	"io"
	"logCollection/log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// DescParam 定义desc参数的结构
type DescParam struct {
	Type     string `json:"type"`
	SavePath string `json:"save_path"`
}

func main() {
	log.Info("程序启动中...")
	// 创建Gin引擎
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		// 在请求完成后记录日志
		start := time.Now() // 记录请求开始的时间
		// 处理请求
		c.Next()
		// 计算耗时
		latency := time.Since(start)
		// 获取请求相关信息
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()
		log.Info("%s %s %d %s %s", clientIP, method, statusCode, latency, path)
	})

	// 设置文件上传路由
	r.POST("/upload-logs", uploadLogsHandler)

	// 启动服务器，默认监听8080端口
	if err := r.Run(":6088"); err != nil {
		panic(err)
	}
}

// uploadLogsHandler 处理日志文件上传
func uploadLogsHandler(c *gin.Context) {
	// 1. 获取desc参数
	descStr := c.PostForm("desc")
	if descStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "desc参数不能为空",
		})
		return
	}

	// 2. 解析desc参数
	var desc DescParam
	if err := json.Unmarshal([]byte(descStr), &desc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "desc参数格式错误",
			"error":   err.Error(),
		})
		return
	}

	// 3. 确定存储路径
	var basePath string
	if desc.SavePath != "" {
		basePath = desc.SavePath
	} else {
		// 如果没有提供save_path，使用当前目录作为默认路径
		basePath, _ = os.Getwd()
	}

	// 获取当前时间并格式化（年-月-日_时-分-秒）
	currentTime := time.Now().Format("06-01-02_15-04-05")
	// 创建"日志"子目录
	logDir := filepath.Join(basePath, "日志", currentTime)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "无法创建日志存储目录",
			"error":   err.Error(),
		})
		return
	}

	// 4. 处理上传的多个文件
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "获取文件失败",
			"error":   err.Error(),
		})
		return
	}

	files := form.File["file"] // 获取所有名为"logs"的文件
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "未上传任何日志文件",
		})
		return
	}

	// 保存所有文件
	savedFiles := make([]string, 0, len(files))
	for _, file := range files {
		// 构建保存路径
		dst := getUniqueFilePath(logDir, file.Filename)

		// 打开上传的文件
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "无法打开上传的文件",
				"error":   err.Error(),
			})
			return
		}
		defer src.Close()

		// 创建目标文件
		out, err := os.Create(dst)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "无法创建目标文件",
				"error":   err.Error(),
			})
			return
		}
		defer out.Close()

		// 复制文件内容
		if _, err := io.Copy(out, src); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "保存文件失败",
				"error":   err.Error(),
			})
			return
		}

		savedFiles = append(savedFiles, dst)
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"message":  "日志文件上传成功",
		"savePath": logDir,
		"files":    savedFiles,
	})
}

// 获取不重复的文件路径，如果文件已存在则添加_2、_3等后缀
func getUniqueFilePath(dir, filename string) string {
	// 分离文件名和扩展名
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]

	// 检查原始路径是否存在
	path := filepath.Join(dir, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}

	// 如果存在，尝试添加_2、_3等后缀
	counter := 2
	for {
		newFilename := name + "_" + strconv.Itoa(counter) + ext
		newPath := filepath.Join(dir, newFilename)
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
		counter++
	}
}
