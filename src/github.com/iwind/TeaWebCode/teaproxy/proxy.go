package teaproxy

import (
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/tealog"
	"github.com/iwind/TeaGo/logs"
	"net/http"
	"net/url"
	"strings"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

// 代理服务
type ProxyServer struct {
	config        *teaconfigs.ServerConfig
	globalWriters map[*teaconfigs.AccessLogConfig]tealog.AccessLogWriter
}

// 获取新代理服务
func NewServer(config *teaconfigs.ServerConfig) *ProxyServer {
	return &ProxyServer{
		config:        config,
		globalWriters: map[*teaconfigs.AccessLogConfig]tealog.AccessLogWriter{},
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
	server := listenerConfig.FindNamedServer(domain)
	if server == nil {
		http.Error(writer, "404 PAGE NOT FOUND", http.StatusNotFound)
		return
	}

	req := NewRequest(rawRequest)
	req.root = server.Root
	req.host = reqHost
	req.method = rawRequest.Method
	req.uri = rawRequest.URL.RequestURI()
	req.scheme = "http" // @TODO 支持 https

	// 查找Location
	err := this.filterRequest(server, req, 0)
	if err != nil {
		req.serverError(writer)
		logs.Error(errors.New(reqHost + rawRequest.URL.String() + ": " + err.Error()))
		return
	}

	// 处理请求
	req.Call(writer)
}

func (this *ProxyServer) filterRequest(server *teaconfigs.ServerConfig, req *Request, redirects int) error {
	if redirects > 8 {
		return errors.New("too many redirects")
	}
	redirects ++

	uri, err := url.ParseRequestURI(req.uri)
	if err != nil {
		return err
	}
	path := uri.Path

	req.root = server.Root

	// location的相关配置
	for _, location := range server.Locations {
		if location.Match(path) {
			if !location.On {
				continue
			}
			if len(location.Root) > 0 {
				req.root = location.Root
			}

			// rewrite相关配置
			if len(location.Rewrite) > 0 {
				for _, rule := range location.Rewrite {
					if !rule.On {
						continue
					}
					if rule.Apply(path, func(source string) string {
						return source
					}) {
						// @TODO 支持带host前缀的URL，比如：http://google.com/hello/world
						newURI, err := url.ParseRequestURI(rule.TargetURL())
						if err != nil {
							req.uri = rule.TargetURL()
							return nil
						}
						if len(newURI.RawQuery) > 0 {
							req.uri = newURI.Path + "?" + newURI.RawQuery
							if len(uri.RawQuery) > 0 {
								req.uri += "&" + uri.RawQuery
							}
						} else {
							req.uri = newURI.Path
							if len(uri.RawQuery) > 0 {
								req.uri += "?" + uri.RawQuery
							}
						}

						switch rule.TargetType() {
						case teaconfigs.RewriteTargetURL:
							return this.filterRequest(server, req, redirects)
						case teaconfigs.RewriteTargetProxy:
							proxyId := rule.TargetProxy()
							server, found := FindServer(proxyId)
							if !found {
								return errors.New("server with '" + proxyId + "' not found")
							}
							if !server.On {
								return errors.New("server with '" + proxyId + "' not available now")
							}
							return this.filterRequest(server, req, redirects)
						}
						return nil
					}
				}
			}

			// fastcgi
			if location.Fastcgi != nil && location.Fastcgi.On {
				req.fastcgi = location.Fastcgi
				return nil
			}

			// proxy
			if len(location.Proxy) > 0 {
				server, found := FindServer(location.Proxy)
				if !found {
					return errors.New("server with '" + location.Proxy + "' not found")
				}
				if !server.On {
					return errors.New("server with '" + location.Proxy + "' not available now")
				}
				return this.filterRequest(server, req, redirects)
			}

			// backends
			if len(location.Backends) > 0 {
				backend := location.NextBackend()
				if backend == nil {
					return errors.New("no backends available")
				}
				req.backend = backend
				return nil
			}

			// root
			if len(location.Root) > 0 {
				req.root = location.Root
				return nil
			}
		}
	}

	// server的相关配置
	if len(server.Rewrite) > 0 {
		for _, rule := range server.Rewrite {
			if !rule.On {
				continue
			}
			if rule.Apply(path, func(source string) string {
				return source
			}) {
				// @TODO 支持带host前缀的URL，比如：http://google.com/hello/world
				newURI, err := url.ParseRequestURI(rule.TargetURL())
				if err != nil {
					req.uri = rule.TargetURL()
					return nil
				}
				if len(newURI.RawQuery) > 0 {
					req.uri = newURI.Path + "?" + newURI.RawQuery
					if len(uri.RawQuery) > 0 {
						req.uri += "&" + uri.RawQuery
					}
				} else {
					if len(uri.RawQuery) > 0 {
						req.uri = newURI.Path + "?" + uri.RawQuery
					}
				}

				switch rule.TargetType() {
				case teaconfigs.RewriteTargetURL:
					return this.filterRequest(server, req, redirects)
				case teaconfigs.RewriteTargetProxy:
					proxyId := rule.TargetProxy()
					server, found := FindServer(proxyId)
					if !found {
						return errors.New("server with '" + proxyId + "' not found")
					}
					if !server.On {
						return errors.New("server with '" + proxyId + "' not available now")
					}
					return this.filterRequest(server, req, redirects)
				}
				return nil
			}
		}
	}

	// fastcgi
	if server.Fastcgi != nil && server.Fastcgi.On {
		req.fastcgi = server.Fastcgi
		return nil
	}

	// proxy
	if len(server.Proxy) > 0 {
		server, found := FindServer(server.Proxy)
		if !found {
			return errors.New("server with '" + server.Proxy + "' not found")
		}
		if !server.On {
			return errors.New("server with '" + server.Proxy + "' not available now")
		}
		return this.filterRequest(server, req, redirects)
	}

	// 转发到后端
	backend := server.NextBackend()
	if backend == nil {
		if len(req.root) == 0 {
			return errors.New("no backends available")
		}
	}
	req.backend = backend

	return nil
}
