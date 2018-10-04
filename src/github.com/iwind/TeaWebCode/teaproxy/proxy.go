package teaproxy

import (
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/tealogs"
	"github.com/iwind/TeaGo/logs"
	"net/http"
	"strings"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

// 代理服务
type ProxyServer struct {
	config        *teaconfigs.ServerConfig
	globalWriters map[*teaconfigs.AccessLogConfig]tealogs.AccessLogWriter
}

// 获取新代理服务
func NewServer(config *teaconfigs.ServerConfig) *ProxyServer {
	return &ProxyServer{
		config:        config,
		globalWriters: map[*teaconfigs.AccessLogConfig]tealogs.AccessLogWriter{},
	}
}

// 处理请求
func (this *ProxyServer) handle(writer http.ResponseWriter, rawRequest *http.Request, listenerConfig *teaconfigs.ListenerConfig) {
	reqHost := rawRequest.Host
	colonIndex := strings.Index(reqHost, ":")
	domain := ""
	if colonIndex < 0 {
		domain = reqHost
	} else {
		domain = reqHost[:colonIndex]
	}
	server, serverName := listenerConfig.FindNamedServer(domain)
	if server == nil {
		http.Error(writer, "404 PAGE NOT FOUND", http.StatusNotFound)
		return
	}

	req := NewRequest(rawRequest)
	req.host = reqHost
	req.method = rawRequest.Method
	req.uri = rawRequest.URL.RequestURI()
	req.scheme = "http" // @TODO 支持 https
	req.serverName = serverName
	req.serverAddr = listenerConfig.Address
	req.root = server.Root
	req.index = server.Index
	req.charset = server.Charset

	// 查找Location
	err := req.configure(server, 0)
	if err != nil {
		req.serverError(writer)
		logs.Error(errors.New(reqHost + rawRequest.URL.String() + ": " + err.Error()))
		return
	}

	// 处理请求
	req.Call(writer)
}
