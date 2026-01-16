package middleware

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func Gzip() gin.HandlerFunc {
	// DefaultCompression 是平衡了速度和压缩率的最佳选择
	// ExcludedExtensions 可以排除不需要压缩的文件类型（图片通常不需要 gzip）
	return gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{".png", ".gif", ".jpeg", ".jpg"}))
}