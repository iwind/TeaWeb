package teaproxy

import (
	"net/http"
	"strings"
	"github.com/iwind/TeaWebCode/teaconst"
	"fmt"
	"net/url"
	"net/http/httputil"
	"github.com/iwind/TeaGo/logs"
	"path/filepath"
	"os"
	"sync"
	"io"
	"bufio"
)

type Request struct {
	req        *http.Request
	cacheKey   string
	remoteAddr string
	remotePort string
	variables  map[string]string
	cacheFile  string
	cacheMutex *sync.Mutex
}

func NewRequest(req *http.Request) *Request {
	request := &Request{
		req:        req,
		cacheMutex: &sync.Mutex{},
	}

	index := strings.LastIndex(req.RemoteAddr, ":")
	if index == -1 {
		request.remoteAddr = req.RemoteAddr
	} else {
		request.remoteAddr = req.RemoteAddr[:index]
		request.remotePort = req.RemoteAddr[index+1:]
	}

	request.variables = map[string]string{
		"teaVersion": teaconst.TeaVersion,

		"remoteAddr": request.remoteAddr,
		"remotePort": request.remotePort,

		"requestURI":    req.RequestURI,
		"requestLength": fmt.Sprintf("%d", req.ContentLength),
		"requestMethod": req.Method,
		"scheme":        req.URL.Scheme,
	}

	return request
}

func (request *Request) HttpRequest() *http.Request {
	return request.req
}

func (request *Request) SetVariable(name string, value string) {
	request.variables[name] = value
}

func (request *Request) Variables() map[string]string {
	return request.variables
}

func (request *Request) RemoteAddr() string {
	return request.remoteAddr
}

func (request *Request) RemotePort() string {
	return request.remotePort
}

func (request *Request) URL() *url.URL {
	return request.req.URL
}

func (request *Request) Header() http.Header {
	return request.req.Header
}

func (request *Request) Host() string {
	return request.req.Host
}

func (request *Request) SetHost(host string) {
	request.req.Host = host
}

func (request *Request) Proto() string {
	return request.req.Proto
}

func (request *Request) SetRequestURI(requestURI string) {
	request.req.RequestURI = requestURI
}

func (request *Request) SetCacheFile(cacheFile string) {
	request.cacheFile = cacheFile
}

func (request *Request) CacheFile() string {
	return request.cacheFile
}

func (request *Request) ShouldCache() bool {
	return len(request.CacheFile()) > 0
}

func (request *Request) WriteCache(resp *http.Response) {
	if len(request.cacheFile) == 0 {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return
	}

	respBytes, err := httputil.DumpResponse(resp, true)
	if err != nil {
		logs.Errorf("[DumpResponse]%s", err.Error())
		return
	}

	// 加并发锁
	request.cacheMutex.Lock()
	defer request.cacheMutex.Unlock()

	cacheDir := filepath.Dir(request.CacheFile())
	_, err = os.Stat(cacheDir)
	if err != nil {
		err = os.MkdirAll(cacheDir, 0766)
		if err != nil {
			logs.Error(err)
			return
		}
	}

	fp, err := os.OpenFile(request.CacheFile(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		logs.Error(err)
		return
	}
	_, err = fp.Write(respBytes)
	if err != nil {
		logs.Error(err)
	}
	fp.Close()
}

// 从缓存中读取响应数据
func (request *Request) ReadCache(writer http.ResponseWriter, cacheFile string) bool {
	// @TODO 支持ETAG识别304 not modified
	// @TODO 支持内存缓存

	var resp *http.Response

	fp, err := os.OpenFile(cacheFile, os.O_RDONLY, 0666)
	if err != nil {
		return false
	}

	defer fp.Close()

	resp, err = http.ReadResponse(bufio.NewReader(fp), nil)
	if err != nil {
		return false
	}

	// 输出后端服务返回的Header
	for key, values := range resp.Header {
		if key == "Connection" {
			continue
		}

		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}

	io.Copy(writer, resp.Body)
	resp.Body.Close()
	return true
}
