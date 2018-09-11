package teaconfigs

import (
	"github.com/iwind/TeaGo/utils/string"
	"regexp"
	"strings"
)

// 路径配置
type LocationConfig struct {
	Path    string `yaml:"path" json:"path"`
	Pattern string `yaml:"pattern" json:"pattern"`
	reg     *regexp.Regexp

	Async   bool         `yaml:"async" json:"async"`     // @TODO
	Notify  []string     `yaml:"notify" json:"notify"`   // @TODO
	LogOnly bool         `yaml:"logOnly" json:"logOnly"` // @TODO
	Cache   *CacheConfig `yaml:"cache" json:"cache"`     // @TODO
	Root    string       `yaml:"root" json:"root"`       // @TODO
	Charset string       `yaml:"charset" json:"charset"` // @TODO

	// 日志
	AccessLog []*AccessLogConfig // @TODO

	// 参考 http://nginx.org/en/docs/http/ngx_http_headers_module.html#add_header
	Headers []HeaderConfig // @TODO

	// 参考：http://nginx.org/en/docs/http/ngx_http_access_module.html
	Allow []string `yaml:"allow" json:"allow"` //@TODO
	Deny  []string `yaml:"deny" json:"deny"`   //@TODO

	Rewrite []*Rewrite // 重写规则 @TODO
}

func (config *LocationConfig) Validate() error {
	if len(config.Pattern) > 0 {
		reg, err := stringutil.RegexpCompile(config.Pattern)
		if err != nil {
			return err
		}

		config.reg = reg
	}

	err := config.Cache.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (config *LocationConfig) Match(path string) bool {
	if config.reg != nil {
		return config.reg.MatchString(path)
	}

	return strings.HasPrefix(path, config.Path)
}
