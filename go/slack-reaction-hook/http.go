package main

import (
	"net/http"
	"sync/atomic"

	"go.uber.org/zap"
)

var (
	healthy int32
)

// healthz...ヘルスチェック用リクエスト。204を返す
func healthz(w http.ResponseWriter, req *http.Request) {
	if atomic.LoadInt32(&healthy) == 1 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// loging...httpリクエストをロギングする
func logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer log.logger.Info(
				r.RemoteAddr,
				zap.String(r.Method, r.URL.Path),
				zap.String("UserAgent", r.UserAgent()),
			)
			next.ServeHTTP(w, r)
		})
	}
}
