package teastats

import (
	"testing"
	"github.com/iwind/TeaWebCode/tealogs"
)

func TestProcess(t *testing.T) {
	log := &tealogs.AccessLog{
		ServerId:    "123456",
		RequestTime: 0.023,
		RemoteAddr:  "183.131.156.10",
		UserAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
		RequestURI:  "/",

		Scheme: "http",
		Host:   "localhost",

		SentHeader: map[string][]string{
			"Content-Type": {"text/html; charset=utf-8"},
		},
	}
	log.Parse()
	new(Processor).Process(log)
}
