package teaproxy

import (
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"sync"
)

func Start() {
	startProxies()
}

func Wait() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func startProxies() {
	listenerConfigs, err := teaconfigs.ParseConfigs()
	if err != nil {
		logs.Error(err)
		return
	}

	for _, config := range listenerConfigs {
		listener := NewListener(config)
		go listener.Start()
	}
}
