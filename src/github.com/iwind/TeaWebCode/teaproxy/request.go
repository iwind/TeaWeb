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
	"github.com/tomasen/fcgi_client"
	"fmt"
	"github.com/iwind/TeaGo/types"
)

type Request struct {
	raw *http.Request

	scheme     string
	uri        string
	host       string
	method     string
	serverAddr string
	headers    http.Header // @TODO

	root    string
	backend *teaconfigs.ServerBackendConfig
	fastcgi *teaconfigs.FastcgiConfig
	proxy   *teaconfigs.ServerConfig
}

func NewRequest(rawRequest *http.Request) *Request {
	return &Request{
		raw: rawRequest,
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

	_, err = io.Copy(writer, fp)

	if err != nil {
		this.serverError(writer)
		logs.Error(err)
		return nil
	}
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
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return nil
	}
	defer resp.Body.Close()
	writer.WriteHeader(resp.StatusCode)

	io.Copy(writer, resp.Body)

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

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		logs.Error(err)
		return nil
	}
	return nil
}

func (this *Request) notFoundError(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusNotFound)
	writer.Write([]byte("404 PAGE NOT FOUND"))
}

func (this *Request) serverError(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte("500 INTERNAL SERVER ERROR"))
}
