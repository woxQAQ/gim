package middleware

import (
	"net/http"
	"time"

	"github.com/woxQAQ/gim/pkg/logger"
)

// responseWriter 包装http.ResponseWriter以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger 创建一个用于记录HTTP请求的中间件
func Logger(l logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			// 创建自定义ResponseWriter来捕获状态码
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rw, r)

			duration := time.Since(startTime)

			// 获取查询参数
			queryParams := r.URL.Query().Encode()
			if queryParams == "" {
				queryParams = "-"
			}

			// 记录详细的请求日志
			l.Info("收到API请求",
				logger.String("method", r.Method),
				logger.String("path", r.URL.Path),
				logger.String("query", queryParams),
				logger.String("remote_addr", r.RemoteAddr),
				logger.String("user_agent", r.UserAgent()),
				logger.String("content_type", r.Header.Get("Content-Type")),
				logger.String("referer", r.Referer()),
				logger.Int("status", rw.statusCode),
				logger.Int64("duration_ms", duration.Milliseconds()),
			)
		})
	}
}
