package teastats

import (
	"testing"
	"github.com/iwind/TeaWebCode/tealog"
)

func TestTopCostStat_Process(t *testing.T) {
	log := &tealog.AccessLog{
		RequestTime: 1,
		Scheme:      "http",
		Host:        "google.com",
		RequestURI:  "/hello",
	}
	stat := new(TopCostStat)
	stat.Process(log)
}
