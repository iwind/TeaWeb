package teaconfigs

// 参考 http://nginx.org/en/docs/http/ngx_http_headers_module.html#add_header
type HeaderConfig struct {
	Name   string `yaml:"name"`   // @TODO
	Value  string `yaml:"value"`  // @TODO
	Always bool   `yaml:"always"` // @TODO
	Code   []string               // @TODO
}
