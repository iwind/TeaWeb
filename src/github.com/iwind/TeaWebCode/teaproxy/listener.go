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
	config  *teaconfigs.ListenerConfig
	servers map[*teaconfigs.ServerConfig]*ProxyServer
	locker  *sync.Mutex
	server  *http.Server
}

func NewListener(config *teaconfigs.ListenerConfig) *Listener {
	listener := &Listener{
		config:  config,
		servers: map[*teaconfigs.ServerConfig]*ProxyServer{},
		locker:  &sync.Mutex{},
	}
	LISTENERS = append(LISTENERS, listener)
	return listener
}

func (this *Listener) Start() {
	httpHandler := http.NewServeMux()
	httpHandler.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		// @TODO 检查域名，通过域名取得对应的Server
		config := this.config.Servers[0]
		server, found := this.servers[config]
		if !found {
			server = NewServer(config)
			this.locker.Lock()
			this.servers[config] = server
			this.locker.Unlock()
		}

		server.handle(writer, req, this.config)
	})

	var err error

	this.server = &http.Server{Addr: this.config.Address, Handler: httpHandler}
	if this.config.SSL != nil && this.config.SSL.On {
		logs.Println("start ssl listener on", this.config.Address)
		err = this.server.ListenAndServeTLS(Tea.ConfigFile(this.config.SSL.Certificate), Tea.ConfigFile(this.config.SSL.CertificateKey))
	}

	if this.config.Http {
		logs.Println("start listener on", this.config.Address)
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
