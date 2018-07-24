package tealog

import (
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	SharedLogger().Push(&AccessLog{
		Request: "Get /",
	}, []AccessLogWriter{})
	SharedLogger().Push(&AccessLog{
		Request: "Get /hello",
	}, []AccessLogWriter{})

	time.Sleep(10 * time.Second)
}
