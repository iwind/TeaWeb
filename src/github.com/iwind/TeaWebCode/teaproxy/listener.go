package teaproxy

import (
	"github.com/iwind/TeaWebCode/teaconfigs"
	"net/http"
	"github.com/iwind/TeaGo/logs"
	"sync"
	"github.com/iwind/TeaGo/Tea"
)

// 监听服务定义
type Listener struct {
	listenerConfig *teaconfigs.ListenerConfig
	servers        map[*teaconfigs.ServerConfig]*ProxyServer
	locker         *sync.Mutex
	server         *http.Server
}

func NewListener(config *teaconfigs.ListenerConfig) *Listener {
	listener := &Listener{
		listenerConfig: config,
		servers:        map[*teaconfigs.ServerConfig]*ProxyServer{},
		locker:         &sync.Mutex{},
	}
	LISTENERS = append(LISTENERS, listener)
	return listener
}

func (this *Listener) Start() {
	httpHandler := http.NewServeMux()
	httpHandler.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		// @TODO 检查域名，通过域名取得对应的Server
		config := this.listenerConfig.Servers[0]
		server, found := this.servers[config]
		if !found {
			server = NewServer(config)
			this.locker.Lock()
			this.servers[config] = server
			this.locker.Unlock()
		}

		server.handle(writer, req, this.listenerConfig)
	})

	var err error

	this.server = &http.Server{Addr: this.listenerConfig.Address, Handler: httpHandler}
	if this.listenerConfig.SSL != nil && this.listenerConfig.SSL.On {
		logs.Println("start ssl listener on", this.listenerConfig.Address)
		err = this.server.ListenAndServeTLS(Tea.ConfigFile(this.listenerConfig.SSL.Certificate), Tea.ConfigFile(this.listenerConfig.SSL.CertificateKey))
	}

	if this.listenerConfig.Http {
		logs.Println("start listener on", this.listenerConfig.Address)
		err = this.server.ListenAndServe()
	}

	if err != nil {
		logs.Error(err)
		return
	}
}

func (this *Listener) Shutdown() error {
	return this.server.Shutdown(nil)
}
