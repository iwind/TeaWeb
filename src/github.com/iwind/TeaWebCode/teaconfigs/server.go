package teaconfigs

import (
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/Tea"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/iwind/TeaGo/logs"
)

// 服务配置
type ServerConfig struct {
	On bool `yaml:"on" json:"on"` // 是否开启 @TODO

	Id          string                 `yaml:"id" json:"id"`                   // ID
	Description string                 `yaml:"description" json:"description"` // 描述
	Name        []string               `yaml:"name" json:"name"`               // 域名
	Http        bool                   `yaml:"http" json:"http"`               // 是否支持HTTP
	Listen      []string               `yaml:"listen" json:"listen"`           // 监听地址
	Root        string                 `yaml:"root" json:"root"`               // 根目录 @TODO
	Backends    []*ServerBackendConfig `yaml:"backends" json:"backends"`       // 后端服务器配置
	Locations   []*LocationConfig      `yaml:"locations" json:"locations"`     // 地址配置
	Charset     string                 `yaml:"charset" json:"charset"`         // 字符集 @TODO

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

	Filename string `yaml:"filename" json:"filename"` // 配置文件名
}

// 从目录中加载配置
func LoadServerConfigsFromDir(dirPath string) []*ServerConfig {
	servers := []*ServerConfig{}

	dir := files.NewFile(dirPath)
	subFiles := dir.Glob("*.proxy.conf")
	files.Sort(subFiles, files.SortTypeModifiedTimeReverse)
	for _, configFile := range subFiles {
		reader, err := configFile.Reader()
		if err != nil {
			logs.Error(err)
			continue
		}

		config := &ServerConfig{}
		err = reader.ReadYAML(config)
		if err != nil {
			continue
		}
		config.Filename = configFile.Name()
		servers = append(servers, config)
	}

	return servers
}

// 取得一个新的服务配置
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		On: true,
	}
}

// 从配置文件中读取配置信息
func NewServerConfigFromFile(filename string) (*ServerConfig, error) {
	if len(filename) == 0 {
		return nil, errors.New("filename should not be empty")
	}
	reader, err := files.NewReader(Tea.ConfigFile(filename))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	config := &ServerConfig{}
	err = reader.ReadYAML(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// 添加域名
func (this *ServerConfig) AddName(name ... string) {
	this.Name = append(this.Name, name ...)
}

// 添加监听地址
func (this *ServerConfig) AddListen(address string) {
	this.Listen = append(this.Listen, address)
}

// 添加后端服务
func (this *ServerConfig) AddBackend(config *ServerBackendConfig) {
	this.Backends = append(this.Backends, config)
}

// 校验配置
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

// 将配置写入文件
func (this *ServerConfig) WriteToFile(path string) error {
	writer, err := files.NewWriter(path)
	if err != nil {
		return err
	}
	_, err = writer.WriteYAML(this)
	writer.Close()
	return err
}

// 将配置写入文件
func (this *ServerConfig) WriteToFilename(filename string) error {
	writer, err := files.NewWriter(Tea.ConfigFile(filename))
	if err != nil {
		return err
	}
	_, err = writer.WriteYAML(this)
	writer.Close()
	return err
}
