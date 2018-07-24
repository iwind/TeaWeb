package teaconfigs

import "strings"

// 服务后端配置
type ServerBackendConfig struct {
	Name        []string `yaml:"name"`
	Address     string   `yaml:"address"`
	Weight      uint     `yaml:"weight"`      //@TODO
	IsBackup    bool     `yaml:"backup"`      //@TODO
	FailTimeout string   `yaml:"failTimeout"` //@TODO
	SlowStart   string   `yaml:"slowStart"`   //@TODO
	MaxFails    uint     `yaml:"maxFails"`    //@TODO
	MaxConns    uint     `yaml:"maxConns"`    //@TODO
	IsDown      bool     `yaml:"down"`        //@TODO
}

func (config *ServerBackendConfig) Validate() error {
	// 是否有端口
	if strings.Index(config.Address, ":") == -1 {
		// @TODO 如果是tls，则为443
		config.Address += ":80"
	}

	return nil
}