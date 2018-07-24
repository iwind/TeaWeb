package teaconfigs

// 服务配置
type ServerConfig struct {
	Listen    []string               `yaml:"listen"`
	Name      []string               `yaml:"name"`
	Root      string                 `yaml:"root"` // @TODO
	Backends  []*ServerBackendConfig `yaml:"backends"`
	Locations []*LocationConfig      `yaml:"locations"`
	Charset   string                 `yaml:"charset"` // @TODO

	Async   bool     `yaml:"async"`   // @TODO
	Notify  []string `yaml:"notify"`  // @TODO
	LogOnly bool     `yaml:"logOnly"` // @TODO

	AccessLog []*AccessLogConfig `yaml:"accessLog"` // @TODO

	// 参考 http://nginx.org/en/docs/http/ngx_http_headers_module.html#add_header
	Headers []*HeaderConfig `yaml:"header"` // @TODO

	// 参考：http://nginx.org/en/docs/http/ngx_http_access_module.html
	Allow []string `yaml:"allow"` //@TODO
	Deny  []string `yaml:"deny"`  //@TODO
}

func (config *ServerConfig) Validate() error {
	// backends
	for _, backend := range config.Backends {
		err := backend.Validate()
		if err != nil {
			return err
		}
	}

	// locations
	for _, location := range config.Locations {
		err := location.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
