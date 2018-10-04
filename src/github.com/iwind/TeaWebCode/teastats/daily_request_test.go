package teastats

import (
	"testing"
	"github.com/iwind/TeaWebCode/tealogs"
	"time"
)

func TestDailyRequestParse(t *testing.T) {
	log := &tealogs.AccessLog{
		ServerId: "123456",
	}

	stat := new(DailyRequestsStat)
	stat.Process(log)

	time.Sleep(1 * time.Second)
}
