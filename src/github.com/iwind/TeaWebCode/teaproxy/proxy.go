package teaproxy

import (
	"github.com/iwind/TeaWebCode/teaconfigs"
	"time"
	"github.com/iwind/TeaWebCode/tealog"
	"github.com/iwind/TeaWebCode/teaconst"
	"github.com/iwind/TeaGo/logs"
	"net/http"
	"io"
	"context"
	"net"
	"net/url"
	"github.com/iwind/TeaGo/types"
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

// 处理请求
func (this *ProxyServer) _handle(writer http.ResponseWriter, req *http.Request, listenerConfig *teaconfigs.ListenerConfig) {
	// scheme
	scheme := "http"
	if req.URL != nil {
		if len(req.URL.Scheme) == 0 {
			if listenerConfig.SSL != nil && listenerConfig.SSL.On {
				scheme = "https"
			}
		}
	}

	// 初始化日志
	now := time.Now()
	log := &tealog.AccessLog{
		TeaVersion:     teaconst.TeaVersion,
		RequestURI:     req.RequestURI,
		RequestLength:  req.ContentLength,
		RequestMethod:  req.Method,
		Request:        req.Method + " " + req.RequestURI + " " + req.Proto,
		Referer:        req.Referer(),
		Scheme:         scheme,
		Proto:          req.Proto,
		Host:           req.Host,
		ServerName:     req.Host,
		ServerPort:     listenerConfig.Port(),
		ServerProtocol: req.Proto,
		RequestPath:    req.URL.Path,
		UserAgent:      req.UserAgent(),
		Arg:            req.URL.Query(),
		Header:         req.Header,
		TimeISO8601:    now.Format("2006-01-02T15:04:05.000Z07:00"),
		TimeLocal:      now.Format("2/Jan/2006:15:04:05 -0700"),
		Msec:           float64(now.Unix()) + float64(now.Nanosecond())/1000000000,
		Timestamp:      now.Unix(),
	}

	// 写日志
	writers := []tealog.AccessLogWriter{}
	defer func() {
		log.RequestTime = time.Since(now).Seconds()

		// 服务器日志
		if len(writers) == 0 {
			if len(this.config.AccessLog) > 0 {
				for _, accessLogConfig := range this.config.AccessLog {
					writer, found := this.globalWriters[accessLogConfig]
					if found {
						writers = append(writers, writer)
					} else {
						writer, err := tealog.NewAccessLogWriter(accessLogConfig)
						if err != nil {
							logs.Error(err)
						} else {
							this.globalWriters[accessLogConfig] = writer
							writers = append(writers, writer)
						}
					}
				}
			}
		}

		tealog.SharedLogger().Push(log, writers)
	}()

	cookies := map[string]string{}
	for _, cookie := range req.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	log.Cookie = cookies

	// 主机名 @TODO 需要分析 *.xxx.com
	request := NewOldRequest(req)
	host := ""
	if len(this.config.Name) > 0 {
		host = this.config.Name[0]
	}
	request.SetVariable("host", host)

	log.RemoteAddr = request.RemoteAddr()
	log.RemotePort = types.Int(request.RemotePort())
	if req.URL.User != nil {
		log.RemoteUser = req.URL.User.Username()
	}

	// 当前Location定制的特性
	goNext := request.filterLocations(writer, this.config.Locations)
	if !goNext {
		return
	}

	// @TODO 检查是否为代理
	if len(this.config.Backends) > 0 {
		this.proxyPass(writer, request, log)
	}
}

func (this *ProxyServer) proxyPass(writer http.ResponseWriter, request *OldRequest, log *tealog.AccessLog) {
	// 检查后端
	if len(this.config.Backends) == 0 {
		http.Error(writer, "no backends available", http.StatusInternalServerError)
		log.Status = http.StatusInternalServerError
		log.StatusMessage = "no backends available"
		return
	}

	//@TODO 根据一定算法选择一个Backend
	backend := &teaconfigs.ServerBackendConfig{}
	for _, backendConfig := range this.config.Backends {
		backend = backendConfig
		break
	}

	log.BackendAddress = backend.Address

	// 主机名 @TODO 需要分析 *.xxx.com
	host := ""
	if len(this.config.Name) > 0 {
		host = this.config.Name[0]
	} else {
		host = backend.Address
	}

	// 设置代理相关的头部
	request.Header().Set("X-Real-IP", request.RemoteAddr())

	// 参考 https://tools.ietf.org/html/rfc7239
	request.Header().Set("X-Forwarded-For", request.RemoteAddr())
	request.Header().Set("X-Forwarded-Host", request.Host())
	request.Header().Set("X-Forwarded-By", request.RemoteAddr())
	request.Header().Set("X-Forwarded-Proto", request.Proto())

	// 其他头部信息
	request.Header().Set("Connection", "keep-alive")
	if len(host) > 0 {
		request.URL().Host = host
	}

	request.URL().Scheme = "http" //@TODO 支持https
	request.SetRequestURI("")

	// 域名
	if len(host) > 0 {
		request.SetHost(host)
	}

	//@TODO 处理超时等问题
	client := http.Client{
		Timeout: 30 * time.Second,

		// 处理跳转
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if via[0].URL.Host == req.URL.Host {
				http.Redirect(writer, request.HttpRequest(), req.URL.RequestURI(), http.StatusTemporaryRedirect)
			} else {
				http.Redirect(writer, request.HttpRequest(), req.URL.String(), http.StatusTemporaryRedirect)
			}
			return &RedirectError{}
		},
	}

	client.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// 后端地址
			addr = backend.Address

			// 握手配置
			return (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext(ctx, network, addr)
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	resp, err := client.Do(request.HttpRequest())
	if err != nil {
		urlError, ok := err.(*url.Error)
		if ok {
			if _, ok := urlError.Err.(*RedirectError); ok {
				return
			}
		}

		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Status = http.StatusInternalServerError
		log.StatusMessage = err.Error()

		log.BytesSent = int64(len(err.Error()))
		log.BodyBytesSent = log.BytesSent

		return
	}

	// 日志
	log.ContentType = resp.Header.Get("Content-Type")
	log.Args = request.HttpRequest().URL.RawQuery
	log.QueryString = log.Args
	log.Status = resp.StatusCode
	log.StatusMessage = resp.Status

	// 输出后端服务返回的Header
	for key, values := range resp.Header {
		if key == "Connection" {
			continue
		}

		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}

	// 状态码
	writer.WriteHeader(resp.StatusCode)

	// cache
	if request.ShouldCache() {
		request.WriteCache(resp)
	}

	// 输出内容
	written, err := io.Copy(writer, resp.Body)
	if err != nil {
		logs.Error(err)
	} else {
		log.BytesSent = written
		log.BodyBytesSent = written
	}
	resp.Body.Close()
}
