package teaconfigs

import (
	"github.com/iwind/TeaGo/Tea"
	"io/ioutil"
	"errors"
	"github.com/go-yaml/yaml"
	"path/filepath"
	"strings"
)

type ListenerConfig struct {
	Address string
	SSL     *SSLConfig
	Servers []*ServerConfig
}

func ParseConfigs() ([]*ListenerConfig, error) {
	listenerConfigMap := map[string]*ListenerConfig{}

	configsDir := Tea.ConfigDir()
	files, err := filepath.Glob(configsDir + Tea.DS + "*.proxy.conf")
	if err != nil {
		return nil, err
	}

	for _, configFile := range files {
		configData, err := ioutil.ReadFile(configFile)
		if err != nil {
			return nil, err
		}

		config := &ServerConfig{}
		err = yaml.Unmarshal(configData, config)
		if err != nil {
			return nil, err
		}

		if len(config.Listen) == 0 {
			return nil, errors.New("'listen' in config should be empty")
		}

		err = config.Validate()
		if err != nil {
			return nil, err
		}

		for _, address := range config.Listen {
			// 是否有端口
			if strings.Index(address, ":") == -1 {
				if config.SSL != nil && config.SSL.On {
					address += ":443"
				} else {
					address += ":80"
				}
			}

			listenerConfig, found := listenerConfigMap[address]

			if !found {
				listenerConfig = &ListenerConfig{
					Address: address,
					Servers: []*ServerConfig{config},
				}
				listenerConfigMap[address] = listenerConfig
			} else {
				listenerConfig.Servers = append(listenerConfig.Servers, config)
			}

			if config.SSL != nil {
				listenerConfig.SSL = config.SSL
			}
		}
	}

	listenerConfigArray := []*ListenerConfig{}
	for _, listenerConfig := range listenerConfigMap {
		listenerConfigArray = append(listenerConfigArray, listenerConfig)
	}

	return listenerConfigArray, nil
}
