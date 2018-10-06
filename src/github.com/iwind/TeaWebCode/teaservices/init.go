package teaservices

import (
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaWebCode/teaservices/probes"
	_ "github.com/iwind/TeaWebCode/teaservices/probes/apps"
	"time"
)

func init() {
	logs.Println("start service probes")

	go func() {
		time.Sleep(1 * time.Second)

		new(probes.CPUProbe).Run()
		new(probes.MemoryProbe).Run()
		new(probes.NetworkProbe).Run()
		new(probes.DiskProbe).Run()
	}()
}
