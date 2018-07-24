package login

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaweb/configs"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
)

type IndexAction actions.Action

func (action *IndexAction) RunGet() {
	action.Show()
}

func (action *IndexAction) RunPost(params struct {
	Username string
	Password string
	Must     *actions.Must
	Auth     *helpers.UserShouldAuth
}) {
	params.Must.
		Field("username", params.Username).
		Require("请输入用户名").
		Field("password", params.Password).
		Require("请输入密码")

	config := configs.SharedAdminConfig()
	for _, user := range config.Users {
		if user.Username == params.Username && user.Password == params.Password {
			params.Auth.StoreUsername(user.Username)
			action.Next("/", nil, "").Success()
			return
		}
	}

	action.FailMessage(400, "登录失败，请检查用户名密码")
}
