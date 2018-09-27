package teaconfigs

import (
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaWebCode/teaconst"
	"net/http"
	"path/filepath"
)

// Fastcgi配置
// 参考：http://nginx.org/en/docs/http/ngx_http_fastcgi_module.html
type FastcgiConfig struct {
	On          bool              `yaml:"on" json:"on"`                   // @TODO
	Pass        string            `yaml:"pass" json:"pass"`               //@TODO
	Index       string            `yaml:"index" json:"index"`             //@TODO
	Params      map[string]string `yaml:"params" json:"params"`           //@TODO
	ReadTimeout string            `yaml:"readTimeout" json:"readTimeout"` //@TODO

	paramsMap maps.Map
}

// 校验配置
func (this *FastcgiConfig) Validate() error {
	this.paramsMap = maps.NewMap(this.Params)
	if !this.paramsMap.Has("SCRIPT_FILENAME") {
		this.paramsMap["SCRIPT_FILENAME"] = ""
	}
	if !this.paramsMap.Has("SERVER_SOFTWARE") {
		this.paramsMap["SERVER_SOFTWARE"] = "teaweb/" + teaconst.TeaVersion
	}
	if !this.paramsMap.Has("REDIRECT_STATUS") {
		this.paramsMap["REDIRECT_STATUS"] = "200"
	}
	if !this.paramsMap.Has("GATEWAY_INTERFACE") {
		this.paramsMap["GATEWAY_INTERFACE"] = "CGI/1.1"
	}

	return nil
}

func (this *FastcgiConfig) FilterParams(req *http.Request) maps.Map {
	params := maps.NewMap(this.paramsMap)

	//@TODO 处理参数中的${varName}变量

	// 自动添加参数
	script := params.GetString("SCRIPT_FILENAME")
	if len(script) > 0 {
		if !params.Has("SCRIPT_NAME") {
			params["SCRIPT_NAME"] = filepath.Base(script)
		}
		if !params.Has("DOCUMENT_ROOT") {
			params["DOCUMENT_ROOT"] = filepath.Dir(script)
		}
		if !params.Has("PWD") {
			params["PWD"] = filepath.Dir(script)
		}
	}

	return params
}
