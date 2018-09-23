package teaconfigs

import (
	"strings"
	"regexp"
)

const (
	LocationPatternTypePrefix = 1
	LocationPatternTypeExact  = 2
	LocationPatternTypeRegexp = 3
)

// 路径配置
// @TODO 匹配的时候去除路径中多于的斜杠（/）
type LocationConfig struct {
	On      bool   `yaml:"on" json:"on"`           // 是否开启 @TODO
	Pattern string `yaml:"pattern" json:"pattern"` // 匹配规则  @TODO

	patternType int // 规则类型：LocationPattern*

	prefix string // 前缀
	path   string // 精确的路径

	reg *regexp.Regexp // 匹配规则

	caseInsensitive bool // 大小写不敏感
	reverse         bool // 是否翻转规则，比如非前缀，非路径

	Async   bool         `yaml:"async" json:"async"`     // @TODO
	Notify  []string     `yaml:"notify" json:"notify"`   // @TODO
	LogOnly bool         `yaml:"logOnly" json:"logOnly"` // @TODO
	Cache   *CacheConfig `yaml:"cache" json:"cache"`     // @TODO
	Root    string       `yaml:"root" json:"root"`       // @TODO
	Charset string       `yaml:"charset" json:"charset"` // @TODO

	// 日志
	AccessLog []*AccessLogConfig `yaml:"accessLog" json:"accessLog"` // @TODO

	// 参考 http://nginx.org/en/docs/http/ngx_http_headers_module.html#add_header
	Headers []HeaderConfig `yaml:"headers" json:"headers"` // 头信息 @TODO

	// 参考：http://nginx.org/en/docs/http/ngx_http_access_module.html
	Allow []string `yaml:"allow" json:"allow"` // 允许的终端地址 @TODO
	Deny  []string `yaml:"deny" json:"deny"`   // 禁止的终端地址 @TODO

	Rewrite []*RewriteRule `yaml:"rewrite" json:"rewrite"` // 重写规则 @TODO
	Fastcgi *FastcgiConfig `yaml:"fastcgi" json:"fastcgi"` // Fastcgi配置 @TODO
	Proxy   string         `yaml:proxy" json:"proxy"`      //  代理配置 @TODO
}

func NewLocationConfig() *LocationConfig {
	return &LocationConfig{}
}

func (this *LocationConfig) Validate() error {
	// 分析pattern
	this.reverse = false
	this.caseInsensitive = false
	if len(this.Pattern) > 0 {
		spaceIndex := strings.Index(this.Pattern, " ")
		if spaceIndex < 0 {
			this.patternType = LocationPatternTypePrefix
			this.prefix = this.Pattern
		} else {
			cmd := this.Pattern[:spaceIndex]
			pattern := strings.TrimSpace(this.Pattern[spaceIndex+1:])
			if cmd == "*" { // 大小写非敏感
				this.patternType = LocationPatternTypePrefix
				this.prefix = pattern
				this.caseInsensitive = true
			} else if cmd == "!*" { // 大小写非敏感，翻转
				this.patternType = LocationPatternTypePrefix
				this.prefix = pattern
				this.caseInsensitive = true
				this.reverse = true
			} else if cmd == "!" {
				this.patternType = LocationPatternTypePrefix
				this.prefix = pattern
				this.reverse = true
			} else if cmd == "=" {
				this.patternType = LocationPatternTypeExact
				this.path = pattern
			} else if cmd == "=*" {
				this.patternType = LocationPatternTypeExact
				this.path = pattern
				this.caseInsensitive = true
			} else if cmd == "!=" {
				this.patternType = LocationPatternTypeExact
				this.path = pattern
				this.reverse = true
			} else if cmd == "!=*" {
				this.patternType = LocationPatternTypeExact
				this.path = pattern
				this.reverse = true
				this.caseInsensitive = true
			} else if cmd == "~" { // 正则
				this.patternType = LocationPatternTypeRegexp
				reg, err := regexp.Compile(pattern)
				if err != nil {
					return err
				}
				this.reg = reg
				this.path = pattern
			} else if cmd == "!~" {
				this.patternType = LocationPatternTypeRegexp
				reg, err := regexp.Compile(pattern)
				if err != nil {
					return err
				}
				this.reg = reg
				this.reverse = true
				this.path = pattern
			} else if cmd == "~*" { // 大小写非敏感小写
				this.patternType = LocationPatternTypeRegexp
				reg, err := regexp.Compile("(?i)" + pattern)
				if err != nil {
					return err
				}
				this.reg = reg
				this.caseInsensitive = true
				this.path = pattern
			} else if cmd == "!~*" {
				this.patternType = LocationPatternTypeRegexp
				reg, err := regexp.Compile("(?i)" + pattern)
				if err != nil {
					return err
				}
				this.reg = reg
				this.reverse = true
				this.caseInsensitive = true
				this.path = pattern
			} else {
				this.patternType = LocationPatternTypePrefix
				this.prefix = pattern
			}
		}
	} else {
		this.patternType = LocationPatternTypePrefix
		this.prefix = this.Pattern
	}

	// 校验缓存配置
	if this.Cache != nil {
		err := this.Cache.Validate()
		if err != nil {
			return err
		}
	}

	// 校验RewriteRule配置
	for _, rewriteRule := range this.Rewrite {
		err := rewriteRule.Validate()
		if err != nil {
			return err
		}
	}

	// 校验Fastcgi配置
	if this.Fastcgi != nil {
		err := this.Fastcgi.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// 模式类型
func (this *LocationConfig) PatternType() int {
	return this.patternType
}

// 模式字符串
// 去掉了模式字符
func (this *LocationConfig) PatternString() string {
	if this.patternType == LocationPatternTypePrefix {
		return this.prefix
	}
	return this.path
}

// 是否翻转
func (this *LocationConfig) IsReverse() bool {
	return this.reverse
}

// 是否大小写非敏感
func (this *LocationConfig) IsCaseInsensitive() bool {
	return this.caseInsensitive
}

// 判断是否匹配路径
func (this *LocationConfig) Match(path string) bool {
	if this.patternType == LocationPatternTypePrefix {
		if this.reverse {
			if this.caseInsensitive {
				return !strings.HasPrefix(strings.ToLower(path), strings.ToLower(this.prefix))
			} else {
				return !strings.HasPrefix(path, this.prefix)
			}
		} else {
			if this.caseInsensitive {
				return strings.HasPrefix(strings.ToLower(path), strings.ToLower(this.prefix))
			} else {
				return strings.HasPrefix(path, this.prefix)
			}
		}
	}

	if this.patternType == LocationPatternTypeExact {
		if this.reverse {
			if this.caseInsensitive {
				return strings.ToLower(path) != strings.ToLower(this.path)
			} else {
				return path != this.path
			}
		} else {
			if this.caseInsensitive {
				return strings.ToLower(path) == strings.ToLower(this.path)
			} else {
				return path == this.path
			}
		}
	}

	if this.patternType == LocationPatternTypeRegexp {
		if this.reg != nil {
			if this.reverse {
				return !this.reg.MatchString(path)
			} else {
				return this.reg.MatchString(path)
			}
		}

		return this.reverse
	}

	return false
}

// 组合参数为一个字符串
func (this *LocationConfig) SetPattern(pattern string, patternType int, caseInsensitive bool, reverse bool) {
	op := ""
	if patternType == LocationPatternTypePrefix {
		if caseInsensitive {
			op = "*"
			if reverse {
				op = "!*"
			}
		} else {
			if reverse {
				op = "!"
			}
		}
	} else if patternType == LocationPatternTypeExact {
		op = "="
		if caseInsensitive {
			op += "*"
		}
		if reverse {
			op = "!" + op
		}
	} else if patternType == LocationPatternTypeRegexp {
		op = "~"
		if caseInsensitive {
			op += "*"
		}
		if reverse {
			op = "!" + op
		}
	}
	if len(op) > 0 {
		pattern = op + " " + pattern
	}
	this.Pattern = pattern
}
