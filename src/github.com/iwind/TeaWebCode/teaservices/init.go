package teaservices

import "github.com/iwind/TeaGo/logs"

func init() {
	logs.Println("register service probes")

	go func() {
		new(CPUProbe).Run()
		new(MemoryProbe).Run()
		new(NetworkProbe).Run()
		new(DiskProbe).Run()
	}()
}
