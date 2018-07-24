package teaproxy

import (
	"github.com/iwind/TeaWebCode/teaconfigs"
	"net/http"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/utils/string"
	"sync"
)

var variablesReg, _ = stringutil.RegexpCompile("\\$\\{\\w+}")

type Listener struct {
	config  *teaconfigs.ListenerConfig
	servers map[*teaconfigs.ServerConfig]*Server
	locker  *sync.Mutex
}

func NewListener(config *teaconfigs.ListenerConfig) *Listener {
	return &Listener{
		config:  config,
		servers: map[*teaconfigs.ServerConfig]*Server{},
		locker:  &sync.Mutex{},
	}
}

func (listener *Listener) Start() {
	httpHandler := http.NewServeMux()
	httpHandler.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		// @TODO 检查域名，通过域名取得对应的Server
		config := listener.config.Servers[0]
		server, found := listener.servers[config]
		if !found {
			server = NewServer(config)
			listener.locker.Lock()
			listener.servers[config] = server
			listener.locker.Unlock()
		}

		server.handle(writer, req)
	})

	logs.Println("start listener on", listener.config.Address)
	err := http.ListenAndServe(listener.config.Address, httpHandler)
	if err != nil {
		logs.Error(err)
		return
	}
}
