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
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/iwind/TeaGo/Tea"
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

func (this *Request) HttpRequest() *http.Request {
	return this.req
}

func (this *Request) SetVariable(name string, value string) {
	this.variables[name] = value
}

func (this *Request) Variables() map[string]string {
	return this.variables
}

func (this *Request) RemoteAddr() string {
	return this.remoteAddr
}

func (this *Request) RemotePort() string {
	return this.remotePort
}

func (this *Request) URL() *url.URL {
	return this.req.URL
}

func (this *Request) Header() http.Header {
	return this.req.Header
}

func (this *Request) Host() string {
	return this.req.Host
}

func (this *Request) SetHost(host string) {
	this.req.Host = host
}

func (this *Request) Proto() string {
	return this.req.Proto
}

func (this *Request) SetRequestURI(requestURI string) {
	this.req.RequestURI = requestURI
}

func (this *Request) SetCacheFile(cacheFile string) {
	this.cacheFile = cacheFile
}

func (this *Request) CacheFile() string {
	return this.cacheFile
}

func (this *Request) ShouldCache() bool {
	return len(this.CacheFile()) > 0
}

func (this *Request) WriteCache(resp *http.Response) {
	if len(this.cacheFile) == 0 {
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
	this.cacheMutex.Lock()
	defer this.cacheMutex.Unlock()

	cacheDir := filepath.Dir(this.CacheFile())
	_, err = os.Stat(cacheDir)
	if err != nil {
		err = os.MkdirAll(cacheDir, 0766)
		if err != nil {
			logs.Error(err)
			return
		}
	}

	fp, err := os.OpenFile(this.CacheFile(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
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
func (this *Request) ReadCache(writer http.ResponseWriter, cacheFile string) bool {
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

// @TODO 利用当前请求的方法格式化字符串
// 变量使用：${varName} 来表示
func (this *Request) Format(format string) string {
	return format
}

func (this *Request) filterLocations(writer http.ResponseWriter, locations []*teaconfigs.LocationConfig) (goNext bool) {
	if len(locations) == 0 {
		return true
	}

	requestPath := this.URL().Path

	var cacheKey = ""
	// @TODO 提升性能
	for _, location := range locations {
		if !location.On {
			continue
		}

		if location.Match(requestPath) {
			// @TODO 日志
			//logs.Println("accessLog:", location.AccessLog)
			if len(location.AccessLog) > 0 {

			}

			// 缓存
			if location.Cache != nil {
				cacheKey = stringutil.Md5(this.Format(location.Cache.Key))
				cachePath := location.Cache.Path
				if len(cachePath) == 0 {
					cachePath = "cache"
				}
				if !filepath.IsAbs(cachePath) {
					cachePath = Tea.Root + Tea.DS + cachePath
				}

				cacheFile := cachePath + Tea.DS + cacheKey + ".cache"
				if this.ReadCache(writer, cacheFile) {
					return true
				}
				this.SetCacheFile(cacheFile)
			}

			// @TODO 重写规则
			if len(location.Rewrite) > 0 {
				for _, rewriteRule := range location.Rewrite {
					if rewriteRule.Apply(requestPath, func(source string) string {
						return this.Format(source)
					}) {
						logs.Println("find rewrite")
						return false;
					}
				}
			}

			// @TODO fastcgi
		}
	}

	return true
}
