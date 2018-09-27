package teaproxy

import (
	"testing"
	"github.com/iwind/TeaGo/assert"
	"net/http"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"bytes"
)

type testResponseWriter struct {
	a    *assert.Assertion
	data []byte
}

func testNewResponseWriter(a *assert.Assertion) *testResponseWriter {
	return &testResponseWriter{
		a: a,
	}
}

func (this *testResponseWriter) Header() http.Header {
	return http.Header{}
}

func (this *testResponseWriter) Write(data []byte) (int, error) {
	this.data = append(this.data, data ...)
	return len(data), nil
}

func (this *testResponseWriter) WriteHeader(statusCode int) {
}

func (this *testResponseWriter) Close() {
	this.a.Log(string(this.data))
}

func TestRequest_Call(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	request := NewRequest(nil)
	err := request.Call(writer)
	a.IsNotNil(err)
	if err != nil {
		a.Log(err.Error())
	}
}

func TestRequest_CallRoot(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	request := NewRequest(nil)
	request.root = Tea.ViewsDir() + "/@default"
	request.uri = "/layout.css"
	err := request.Call(writer)
	a.IsNil(err)
	writer.Close()
}

func TestRequest_CallBackend(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	req, err := http.NewRequest("GET", "/index.php?__ACTION__=/@wx", nil)
	if err != nil {
		a.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1"
	request := NewRequest(req)
	request.scheme = "http"
	request.host = "wx.balefm.cn"
	request.backend = &teaconfigs.ServerBackendConfig{
		Address: "127.0.0.1",
	}
	request.backend.Validate()
	err = request.Call(writer)
	a.IsNil(err)
	writer.Close()
}

func TestRequest_CallProxy(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	req, err := http.NewRequest("GET", "/index.php?__ACTION__=/@wx", nil)
	if err != nil {
		a.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1"
	request := NewRequest(req)
	request.scheme = "http"
	request.host = "wx.balefm.cn"

	proxy := teaconfigs.NewServerConfig()
	proxy.AddBackend(&teaconfigs.ServerBackendConfig{
		Address: "127.0.0.1:80",
	})
	proxy.AddBackend(&teaconfigs.ServerBackendConfig{
		Address: "127.0.0.1:81",
	})
	request.proxy = proxy

	err = request.Call(writer)
	a.IsNil(err)
	writer.Close()
}

func TestRequest_CallFastcgi(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	req, err := http.NewRequest("POST", "/index.php?__ACTION__=/@wx/box/version", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		a.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	request := NewRequest(req)
	request.scheme = "http"
	request.host = "wx.balefm.cn"
	request.serverAddr = "127.0.0.1:80"

	request.fastcgi = &teaconfigs.FastcgiConfig{
		Params: map[string]string{
			"SCRIPT_FILENAME": "/Users/liuxiangchao/Documents/Projects/pp/apps/baleshop.ppk/index.php",
			//"DOCUMENT_ROOT":   "/Users/liuxiangchao/Documents/Projects/pp/apps/baleshop.ppk",
		},
		Pass: "127.0.0.1:9000",
	}
	request.fastcgi.Validate()
	err = request.Call(writer)
	a.IsNil(err)
	writer.Close()
}
