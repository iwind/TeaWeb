package teaconfigs

// @TODO
// 参考：http://nginx.org/en/docs/http/ngx_http_fastcgi_module.html
type Fastcgi struct {
	Pass        string            //@TODO
	Index       string            //@TODO
	Params      map[string]string //@TODO
	ReadTimeout string            //@TODO
}
