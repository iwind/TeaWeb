package teaconfigs

import "github.com/iwind/TeaGo/files"

// 服务配置
type ServerConfig struct {
	Listen    []string               `yaml:"listen" json:"listen"`       // 监听地址
	Name      []string               `yaml:"name" json:"name"`           // 域名
	Root      string                 `yaml:"root" json:"root"`           // 根目录 @TODO
	Backends  []*ServerBackendConfig `yaml:"backends" json:"backends"`   // 后端服务器配置
	Locations []*LocationConfig      `yaml:"locations" json:"locations"` // 地址配置
	Charset   string                 `yaml:"charset" json:"charset"`     // 字符集 @TODO

	Async   bool     `yaml:"async" json:"async"`     // 请求是否异步处理 @TODO
	Notify  []string `yaml:"notify" json:"notify"`   // 请求转发地址 @TODO
	LogOnly bool     `yaml:"logOnly" json:"logOnly"` // 是否只记录日志 @TODO

	// 访问日志
	AccessLog []*AccessLogConfig `yaml:"accessLog" json:"accessLog"` // 访问日志

	// @TODO 支持ErrorLog

	// SSL
	// @TODO
	SSL *SSLConfig `yaml:"ssl" json:"ssl"`

	// 参考 http://nginx.org/en/docs/http/ngx_http_headers_module.html#add_header
	Headers []*HeaderConfig `yaml:"header" json:"headers"` // @TODO

	// 参考：http://nginx.org/en/docs/http/ngx_http_access_module.html
	Allow []string `yaml:"allow" json:"allow"` //@TODO
	Deny  []string `yaml:"deny" json:"deny"`   //@TODO
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func (this *ServerConfig) AddName(name ... string) {
	this.Name = append(this.Name, name ...)
}

func (this *ServerConfig) AddListen(address string) {
	this.Listen = append(this.Listen, address)
}

func (this *ServerConfig) AddBackend(config *ServerBackendConfig) {
	this.Backends = append(this.Backends, config)
}

func (this *ServerConfig) Validate() error {
	// backends
	for _, backend := range this.Backends {
		err := backend.Validate()
		if err != nil {
			return err
		}
	}

	// locations
	for _, location := range this.Locations {
		err := location.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *ServerConfig) WriteToFile(path string) error {
	writer, err := files.NewWriter(path)
	if err != nil {
		return err
	}
	_, err = writer.WriteYAML(this)
	return err
}
