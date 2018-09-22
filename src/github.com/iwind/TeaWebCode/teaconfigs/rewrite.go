package teaconfigs

import "net/http"

// 重写规则定义
//
// 参考
// - http://nginx.org/en/docs/http/ngx_http_rewrite_module.html
// - https://httpd.apache.org/docs/current/mod/mod_rewrite.html
// - https://httpd.apache.org/docs/2.4/rewrite/flags.html
type RewriteRule struct {
	On      bool     `yaml:"on" json:"on"`           // 是否开启 @TODO
	Cond    []string `yaml:"cond" json:"cond"`       // 开启的条件 @TODO
	Pattern string   `yaml:"pattern" json:"pattern"` // 规则 @TODO
	Replace string   `yaml:"replace" json:"replace"` // 要替换成的URL @TODO
}

func (this *RewriteRule) Validate() error {
	return nil
}

// 对某个请求执行规则 @TODO
func (this *RewriteRule) Apply(req *http.Request) {

}
