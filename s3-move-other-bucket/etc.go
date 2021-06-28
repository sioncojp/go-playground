package main

import (
	"time"
)

// isBeforeDay...x日前のデータを削除するために、ジャッジする
func isBeforeDay(t *time.Time, beforeDay int) bool {
	utc, _ := time.LoadLocation("UTC")
	t1 := time.Now().In(utc)

	return t.Before(t1.AddDate(0, 0, -beforeDay))
}
