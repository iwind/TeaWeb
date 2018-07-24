package teaconfigs

import (
	"github.com/iwind/TeaGo/logs"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

// 参考 http://nginx.org/en/docs/http/ngx_http_log_module.html#access_log
type AccessLogConfig struct {
	Target string                 `yaml:"target"`
	Off    bool                   `yaml:"off"`
	Config map[string]interface{} `yaml:"config"`
}

func (config *AccessLogConfig) Validate() {
	if len(config.Target) == 0 {
		logs.Error(errors.New("invalid access log target '" + config.Target + "'"))
	}
}
