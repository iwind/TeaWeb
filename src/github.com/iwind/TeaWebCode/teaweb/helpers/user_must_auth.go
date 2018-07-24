package helpers

import (
	"github.com/iwind/TeaGo/actions"
)

type UserMustAuth struct {
	Username string
}

func (auth *UserMustAuth) RunAction(actionPtr interface{}, paramName string) (goNext bool) {
	var action = actionPtr.(actions.ActionWrapper).Object()
	var session = action.Session()
	var username = session.StringValue("username")
	if len(username) == 0 {
		auth.login(action)
		return false
	}

	auth.Username = username

	// 初始化内置方法
	action.ViewFunc("teaTitle", func() string {
		return action.Data["teaTitle"].(string)
	})

	// 初始化变量
	action.Data["teaTitle"] = "TeaWeb管理平台"
	action.Data["teaUsername"] = username
	action.Data["teaMenu"] = ""
	action.Data["teaModules"] = []map[string]interface{}{
		{
			"code":     "proxy",
			"menuName": "代理设置",
		},
		{
			"code":     "log",
			"menuName": "访问日志",
		},
		/**{
			"code":     "stat",
			"menuName": "统计",
		},
		{
			"code":     "monitor",
			"menuName": "监控",
		},**/
	}
	action.Data["teaSubMenus"] = []map[string]interface{}{}
	action.Data["teaTabbar"] = []map[string]interface{}{}

	return true
}

func (auth *UserMustAuth) login(action *actions.ActionObject) {
	action.RedirectURL("/login")
}
