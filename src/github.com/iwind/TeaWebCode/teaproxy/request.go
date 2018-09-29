package teaproxy

import (
	"net/http"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/iwind/TeaGo/Tea"
	"strings"
	"os"
	"github.com/iwind/TeaGo/logs"
	"io"
	"time"
	"context"
	"net"
	"net/url"
	"fmt"
	"github.com/iwind/TeaGo/types"
	"regexp"
	"github.com/iwind/TeaWebCode/teaconst"
	"github.com/iwind/TeaWebCode/teaproxy/fcgiclient"
)

var requestVarReg = regexp.MustCompile("\\${[\\w.-]+}")

// 请求定义
type Request struct {
	raw *http.Request

	scheme     string
	uri        string
	host       string
	method     string
	serverName string // @TODO
	serverAddr string

	root    string
	backend *teaconfigs.ServerBackendConfig
	fastcgi *teaconfigs.FastcgiConfig
	proxy   *teaconfigs.ServerConfig

	// 执行请求
	filePath string

	requestFromTime       time.Time
	requestTime           float64 // @TODO
	responseBytesSent     int64   // @TODO
	responseBodyBytesSent int64   // @TODO
	responseStatus        int     // @TODO
	responseStatusMessage string  // @TODO
	requestTimeISO8601    string
	requestTimeLocal      string
	requestMsec           float64
	requestTimestamp      int64
}

// 获取新的请求
func NewRequest(rawRequest *http.Request) *Request {
	now := time.Now()
	return &Request{
		raw:                rawRequest,
		requestFromTime:    now,
		requestTimestamp:   now.Unix(),
		requestTimeISO8601: now.Format("2006-01-02T15:04:05.000Z07:00"),
		requestTimeLocal:   now.Format("2/Jan/2006:15:04:05 -0700"),
		requestMsec:        float64(now.Unix()) + float64(now.Nanosecond())/1000000000,
	}
}

func (this *Request) Call(writer http.ResponseWriter) error {
	if this.backend != nil {
		return this.callBackend(writer)
	}
	if this.proxy != nil {
		return this.callProxy(writer)
	}
	if this.fastcgi != nil {
		return this.callFastcgi(writer)
	}
	if len(this.root) > 0 {
		return this.callRoot(writer)
	}
	return errors.New("unable to handle the request")
}

// @TODO 支持eTag，cache等
func (this *Request) callRoot(writer http.ResponseWriter) error {
	if len(this.uri) == 0 {
		this.notFoundError(writer)
		return nil
	}

	filename := strings.Replace(this.uri, "/", Tea.DS, -1)
	filePath := ""
	if filename[0:1] == Tea.DS {
		filePath = this.root + filename
	} else {
		filePath = this.root + Tea.DS + filename
	}

	this.filePath = filePath

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			this.notFoundError(writer)
			return nil
		} else {
			this.serverError(writer)
			logs.Error(err)
			return nil
		}
	}

	fp, err := os.OpenFile(filePath, os.O_RDONLY, 444)
	if err != nil {
		this.serverError(writer)
		logs.Error(err)
		return nil
	}
	defer fp.Close()

	n, err := io.Copy(writer, fp)

	if err != nil {
		this.serverError(writer)
		logs.Error(err)
		return nil
	}

	this.responseStatus = http.StatusOK
	this.responseStatusMessage = "200 OK"
	this.responseBytesSent = n
	this.responseBodyBytesSent = n
	this.requestTime = time.Since(this.requestFromTime).Seconds()

	return nil
}

func (this *Request) callBackend(writer http.ResponseWriter) error {
	if len(this.backend.Address) == 0 {
		this.serverError(writer)
		logs.Error(errors.New("backend address should not be empty"))
		return nil
	}

	this.raw.URL.Scheme = this.scheme
	this.raw.URL.Host = this.host

	// 设置代理相关的头部
	// 参考 https://tools.ietf.org/html/rfc7239
	this.raw.Header.Set("X-Real-IP", this.raw.RemoteAddr)
	this.raw.Header.Set("X-Forwarded-For", this.raw.RemoteAddr)
	this.raw.Header.Set("X-Forwarded-Host", this.host)
	this.raw.Header.Set("X-Forwarded-By", this.raw.RemoteAddr)
	this.raw.Header.Set("X-Forwarded-Proto", this.raw.Proto)
	//this.raw.Header.Set("Connection", "keep-alive")

	// @TODO 使用client池
	client := http.Client{
		Timeout: 30 * time.Second,

		// 处理跳转
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if via[0].URL.Host == req.URL.Host {
				http.Redirect(writer, this.raw, req.URL.RequestURI(), http.StatusTemporaryRedirect)
			} else {
				http.Redirect(writer, this.raw, req.URL.String(), http.StatusTemporaryRedirect)
			}
			return &RedirectError{}
		},
	}

	client.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// 后端地址
			addr = this.backend.Address

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

	resp, err := client.Do(this.raw)
	if err != nil {
		urlError, ok := err.(*url.Error)
		if ok {
			if _, ok := urlError.Err.(*RedirectError); ok {
				return nil
			}
		}
		this.serverError(writer)
		logs.Error(err)
		return nil
	}
	defer resp.Body.Close()

	// 设置响应代码
	writer.WriteHeader(resp.StatusCode)

	// 设置Header
	for k, v := range resp.Header {
		if k == "Connection" {
			continue
		}
		for _, subV := range v {
			writer.Header().Add(k, subV)
		}
	}

	n, err := io.Copy(writer, resp.Body)
	if err != nil {
		this.serverError(writer)
		logs.Error(err)
		return nil
	}

	// 请求信息
	this.responseStatus = resp.StatusCode
	this.responseStatusMessage = resp.Status
	this.responseBytesSent = n
	this.responseBodyBytesSent = n
	this.requestTime = time.Since(this.requestFromTime).Seconds()

	return nil
}

func (this *Request) callProxy(writer http.ResponseWriter) error {
	backend := this.proxy.NextBackend()
	this.backend = backend
	return this.callBackend(writer)
}

func (this *Request) callFastcgi(writer http.ResponseWriter) error {
	env := this.fastcgi.FilterParams(this.raw)
	if len(this.root) > 0 {
		if !env.Has("DOCUMENT_ROOT") {
			env["DOCUMENT_ROOT"] = this.root
		}
	}
	if !env.Has("REMOTE_ADDR") {
		env["REMOTE_ADDR"] = this.raw.RemoteAddr
	}
	if !env.Has("QUERY_STRING") {
		u, err := url.ParseRequestURI(this.uri)
		if err == nil {
			env["QUERY_STRING"] = u.RawQuery
		} else {
			env["QUERY_STRING"] = this.raw.URL.RawQuery
		}
	}
	if !env.Has("SERVER_NAME") {
		env["SERVER_NAME"] = this.host
	}
	if !env.Has("REQUEST_URI") {
		env["REQUEST_URI"] = this.uri
	}
	if !env.Has("HOST") {
		env["HOST"] = this.host
	}

	if len(this.serverAddr) > 0 {
		if !env.Has("SERVER_ADDR") {
			env["SERVER_ADDR"] = this.serverAddr
		}
		if !env.Has("SERVER_PORT") {
			portIndex := strings.LastIndex(this.serverAddr, ":")
			if portIndex >= 0 {
				env["SERVER_PORT"] = this.serverAddr[portIndex+1:]
			}
		}
	}

	// @TODO 使用连接池
	// @TODO 支持unix://...
	fcgi, err := fcgiclient.Dial("tcp", this.fastcgi.Pass)
	if err != nil {
		this.serverError(writer)
		logs.Error(err)
		return nil
	}

	// 请求相关
	if !env.Has("REQUEST_METHOD") {
		env["REQUEST_METHOD"] = this.method
	}
	if !env.Has("CONTENT_LENGTH") {
		env["CONTENT_LENGTH"] = fmt.Sprintf("%d", this.raw.ContentLength)
	}
	if !env.Has("CONTENT_TYPE") {
		env["CONTENT_TYPE"] = this.raw.Header.Get("Content-Type")
	}

	params := map[string]string{}
	for key, value := range env {
		params[key] = types.String(value)
	}
	resp, err := fcgi.Request(params, this.raw.Body)
	if err != nil {
		this.serverError(writer)
		logs.Error(err)
		return nil
	}

	// 设置响应码
	writer.WriteHeader(resp.StatusCode)

	// 设置Header
	for k, v := range resp.Header {
		if k == "Connection" {
			continue
		}
		for _, subV := range v {
			writer.Header().Add(k, subV)
		}
	}

	n, err := io.Copy(writer, resp.Body)
	if err != nil {
		logs.Error(err)
		return nil
	}

	// 请求信息
	this.responseStatus = resp.StatusCode
	this.responseStatusMessage = resp.Status
	this.responseBytesSent = n
	this.responseBodyBytesSent = n
	this.requestTime = time.Since(this.requestFromTime).Seconds()

	return nil
}

func (this *Request) notFoundError(writer http.ResponseWriter) {
	msg := "404 Page Not Found"

	writer.WriteHeader(http.StatusNotFound)
	writer.Write([]byte(msg))

	this.responseStatus = http.StatusNotFound
	this.responseStatusMessage = msg
	this.responseBodyBytesSent = int64(len(msg))
	this.responseBytesSent = this.responseBodyBytesSent
}

func (this *Request) serverError(writer http.ResponseWriter) {
	msg := "500 Internal Server Error"

	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte(msg))

	this.responseStatus = http.StatusInternalServerError
	this.responseStatusMessage = msg
	this.responseBodyBytesSent = int64(len(msg))
	this.responseBytesSent = this.responseBodyBytesSent
}

func (this *Request) requestRemoteAddr() string {
	return this.raw.RemoteAddr
}

func (this *Request) requestRemotePort() string {
	remoteAddr := this.requestRemoteAddr()
	index := strings.LastIndex(remoteAddr, ":")
	if index < 0 {
		return ""
	} else {
		return remoteAddr[index+1:]
	}
}

func (this *Request) requestRemoteUser() string {
	username, _, ok := this.raw.BasicAuth()
	if !ok {
		return ""
	}
	return username
}

func (this *Request) requestURI() string {
	return this.uri
}

func (this *Request) requestPath() string {
	uri, err := url.ParseRequestURI(this.requestURI())
	if err != nil {
		return ""
	}
	return uri.Path
}

func (this *Request) requestLength() int64 {
	return this.raw.ContentLength
}

func (this *Request) requestMethod() string {
	return this.method
}

func (this *Request) requestFilename() string {
	return this.filePath
}

func (this *Request) requestProto() string {
	return this.raw.Proto
}

func (this *Request) requestReferer() string {
	return this.raw.Referer()
}

func (this *Request) requestUserAgent() string {
	return this.raw.UserAgent()
}

func (this *Request) requestContentType() string {
	return this.raw.Header.Get("Content-Type")
}

func (this *Request) requestString() string {
	return this.method + " " + this.requestURI() + " " + this.requestProto()
}

func (this *Request) requestCookiesString() string {
	var cookies = []string{}
	for _, cookie := range this.raw.Cookies() {
		cookies = append(cookies, url.QueryEscape(cookie.Name)+"="+url.QueryEscape(cookie.Value))
	}
	return strings.Join(cookies, "&")
}

func (this *Request) requestCookie(name string) string {
	cookie, err := this.raw.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Name
}

func (this *Request) requestQueryString() string {
	uri, err := url.ParseRequestURI(this.uri)
	if err != nil {
		return ""
	}
	return uri.RawQuery
}

func (this *Request) requestQueryParam(name string) string {
	uri, err := url.ParseRequestURI(this.uri)
	if err != nil {
		return ""
	}

	v, found := uri.Query()[name]
	if !found {
		return ""
	}
	return strings.Join(v, "&")
}

func (this *Request) requestServerPort() string {
	index := strings.LastIndex(this.serverAddr, ":")
	if index < 0 {
		return ""
	}
	return this.serverAddr[index+1:]
}

func (this *Request) requestHeadersString() string {
	var headers = []string{}
	for k, v := range this.raw.Header {
		for _, subV := range v {
			headers = append(headers, k+": "+subV)
		}
	}
	return strings.Join(headers, ";")
}

func (this *Request) requestHeader(key string) string {
	v, found := this.raw.Header[key]
	if !found {
		return ""
	}
	return strings.Join(v, ";")
}

// 利用请求参数格式化字符串
func (this *Request) format(source string) string {
	var varName = ""
	return requestVarReg.ReplaceAllStringFunc(source, func(s string) string {
		varName = s[2 : len(s)-1]

		switch varName {
		case "teaVersion":
			return teaconst.TeaVersion
		case "remoteAddr":
			return this.requestRemoteAddr()
		case "remotePort":
			return this.requestRemotePort()
		case "remoteUser":
			return this.requestRemoteUser()
		case "requestURI", "requestUri":
			return this.requestURI()
		case "requestPath":
			return this.requestPath()
		case "requestLength":
			return fmt.Sprintf("%d", this.requestLength())
		case "requestTime":
			return fmt.Sprintf("%.6f", this.requestTime)
		case "requestMethod":
			return this.requestMethod()
		case "requestFilename":
			return this.requestFilename()
		case "scheme":
			return this.scheme
		case "serverProtocol", "proto":
			return this.requestProto()
		case "bytesSent":
			return fmt.Sprintf("%d", this.responseBytesSent)
		case "bodyBytesSent":
			return fmt.Sprintf("%d", this.responseBodyBytesSent)
		case "status":
			return fmt.Sprintf("%d", this.responseStatus)
		case "statusMessage":
			return this.responseStatusMessage
		case "timeISO8601":
			return this.requestTimeISO8601
		case "timeLocal":
			return this.requestTimeLocal
		case "msec":
			return fmt.Sprintf("%.6f", this.requestMsec)
		case "timestamp":
			return fmt.Sprintf("%d", this.requestTimestamp)
		case "host":
			return this.host
		case "referer":
			return this.requestReferer()
		case "userAgent":
			return this.requestUserAgent()
		case "contentType":
			return this.requestContentType()
		case "request":
			return this.requestString()
		case "cookies":
			return this.requestCookiesString()
		case "args", "queryString":
			return this.requestQueryString()
		case "headers":
			return this.requestHeadersString()
		case "serverName":
			return this.serverName
		case "serverPort":
			return this.requestServerPort()
		}

		dotIndex := strings.Index(varName, ".")
		if dotIndex < 0 {
			return s
		}
		prefix := varName[:dotIndex]
		suffix := varName[dotIndex+1:]

		// cookie.
		if prefix == "cookie" {
			return this.requestCookie(suffix)
		}

		// arg.
		if prefix == "arg" {
			return this.requestQueryParam(suffix)
		}

		// header.
		if prefix == "header" || prefix == "http" {
			return this.requestHeader(suffix)
		}

		return s
	})
}
