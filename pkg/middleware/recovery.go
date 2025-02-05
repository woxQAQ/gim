package middleware

import (
	"net/http"

	"github.com/woxQAQ/gim/pkg/logger"
)

// Recovery 创建一个用于恢复panic的中间件
func Recovery(l logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// 设置500状态码
					w.WriteHeader(http.StatusInternalServerError)

					// 记录详细的错误信息
					l.Error("服务器发生panic",
						logger.Any("error", err),
						logger.String("method", r.Method),
						logger.String("path", r.URL.Path),
						logger.String("query", r.URL.RawQuery),
						logger.String("remote_addr", r.RemoteAddr),
						logger.String("user_agent", r.UserAgent()),
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
