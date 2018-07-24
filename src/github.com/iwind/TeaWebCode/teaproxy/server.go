package teaproxy

import (
	"github.com/iwind/TeaWebCode/teaconfigs"
	"time"
	"github.com/iwind/TeaWebCode/tealog"
	"github.com/iwind/TeaWebCode/teaconst"
	"fmt"
	"github.com/iwind/TeaGo/utils/string"
	"path/filepath"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"net/http"
	"io"
	"context"
	"net"
	"net/url"
)

type Server struct {
	config        *teaconfigs.ServerConfig
	globalWriters map[*teaconfigs.AccessLogConfig]tealog.AccessLogWriter
}

func NewServer(config *teaconfigs.ServerConfig) *Server {
	return &Server{
		config:        config,
		globalWriters: map[*teaconfigs.AccessLogConfig]tealog.AccessLogWriter{},
	}
}

func (server *Server) handle(writer http.ResponseWriter, req *http.Request) {
	// 初始化日志
	now := time.Now()
	log := &tealog.AccessLog{
		TeaVersion:    teaconst.TeaVersion,
		RequestURI:    req.RequestURI,
		RequestLength: req.ContentLength,
		RequestMethod: req.Method,
		Request:       req.Method + " " + req.RequestURI + " " + req.Proto,
		Referer:       req.Referer(),
		Scheme:        req.URL.Scheme,
		Proto:         req.Proto,
		Host:          req.Host,
		RequestPath:   req.URL.Path,
		UserAgent:     req.UserAgent(),
		Arg:           req.URL.Query(),
		Header:        req.Header,
		TimeISO8601:   now.Format("2006-01-02T15:04:05.000Z07:00"),
		TimeLocal:     now.Format("2/Jan/2006:15:04:05 -0700"),
		Msec:          fmt.Sprintf("%d.%d", now.Unix(), now.Nanosecond()/1000000),
	}

	// 写日志
	writers := []tealog.AccessLogWriter{}
	defer func() {
		log.RequestTime = time.Since(now).Seconds()

		// 服务器日志
		if len(writers) == 0 {
			if len(server.config.AccessLog) > 0 {
				for _, accessLogConfig := range server.config.AccessLog {
					writer, found := server.globalWriters[accessLogConfig]
					if found {
						writers = append(writers, writer)
					} else {
						writer, err := tealog.NewAccessLogWriter(accessLogConfig)
						if err != nil {
							logs.Error(err)
						} else {
							server.globalWriters[accessLogConfig] = writer
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
	request := NewRequest(req)
	host := server.config.Name[0]
	request.SetVariable("host", host)

	log.RemoteAddr = request.RemoteAddr()
	log.RemotePort = request.RemotePort()
	if req.URL.User != nil {
		log.RemoteUser = req.URL.User.Username()
	}

	// 当前Location定制的特性
	var cacheKey = ""
	if len(server.config.Locations) > 0 {
		// @TODO 提升性能
		// @TODO locations必须是有顺序的
		for _, location := range server.config.Locations {
			if location.Match(request.URL().Path) {
				// @TODO 日志
				logs.Println(location.AccessLog)
				if len(location.AccessLog) > 0 {

				}

				// 缓存
				if location.Cache != nil {
					cacheKey = stringutil.Md5(server.parseVariables(location.Cache.Key, request.Variables()))
					cachePath := location.Cache.Path
					if len(cachePath) == 0 {
						cachePath = "cache"
					}
					if !filepath.IsAbs(cachePath) {
						cachePath = Tea.Root + Tea.Ds + cachePath
					}

					cacheFile := cachePath + Tea.Ds + cacheKey + ".cache"
					if request.ReadCache(writer, cacheFile) {
						return
					}
					request.SetCacheFile(cacheFile)
				}
			}
		}
	}

	// @TODO 检查是否为代理
	if len(server.config.Backends) > 0 {
		server.proxyPass(writer, request, log)
	}
}

func (server *Server) proxyPass(writer http.ResponseWriter, request *Request, log *tealog.AccessLog) {
	// 检查后端
	if len(server.config.Backends) == 0 {
		http.Error(writer, "no backends available", http.StatusInternalServerError)
		log.Status = http.StatusInternalServerError
		log.StatusMessage = "no backends available"
		return
	}

	//@TODO 根据一定算法选择一个Backend
	backend := &teaconfigs.ServerBackendConfig{}
	for _, backendConfig := range server.config.Backends {
		backend = backendConfig
		break
	}

	// 主机名 @TODO 需要分析 *.xxx.com
	host := server.config.Name[0]

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
	request.URL().Scheme = "http"
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
	defer resp.Body.Close()
}

func (server *Server) parseVariables(s string, variables map[string]string) string {
	return variablesReg.ReplaceAllStringFunc(s, func(match string) string {
		value, found := variables[match[2:len(match)-1]]
		if found {
			return value
		}
		return ""
	})
}
