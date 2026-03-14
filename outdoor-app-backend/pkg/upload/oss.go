package upload

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ServerURL 本机服务的地址，⚠️ 答辩时如果换电脑或部署到云服务器，改成对应的 IP
const ServerURL = "http://10.136.181.74:8080"

// SaveImageToLocal 保存图片到本地目录并返回 URL
// multipart.FileHeader 的结构大致如下（简化版）Filename：上传文件的原始文件名（如 "photo.jpg"）。Header：包含该文件部分的 MIME 头部信息，如 Content-Type。Size：文件大小（字节数）。
func SaveImageToLocal(c *gin.Context, file *multipart.FileHeader, folderName string) (string, error) {
	// 1. 检查并创建基础保存目录 (例如: ./uploads/avatars)
	basePath := fmt.Sprintf("./uploads/%s", folderName)
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return "", err
	}

	// 2. 获取文件后缀 (例如: .jpg)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" && ext != ".webp" {
		return "", fmt.Errorf("不支持的图片格式: %s", ext)
	}

	// 3. 生成唯一文件名 (UUID + 后缀)，防止两张图片同名被覆盖
	newFileName := uuid.New().String() + ext

	// 4. 完整的本地保存路径
	dst := filepath.Join(basePath, newFileName)

	// 5. 保存文件到本地磁盘
	if err := c.SaveUploadedFile(file, dst); err != nil { //调用 Gin 的 SaveUploadedFile()否则你需要自己写：file.Open()os.Create()io.Copy()，这就是传入c *gin.Context的作用
		return "", err
	}

	// 6. 拼装可以在前端直接访问的网络 URL 返回
	// 格式例如: http://127.0.0.1:8080/uploads/avatars/xxxx-xxxx.jpg
	imageUrl := fmt.Sprintf("%s/uploads/%s/%s", ServerURL, folderName, newFileName)

	return imageUrl, nil
}

//后面记得写
// 🌟 静态文件映射配置：
// 告诉 Gin：如果前端请求的 URL 是以 /uploads 开头，你就去项目根目录下的 ./uploads 文件夹里找对应的文件返回给他。
//r.Static("/uploads", "./uploads")

// 然后才是你的 API 路由...
// r.POST("/api/upload", handler.UploadImage)
