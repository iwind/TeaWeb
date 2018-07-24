package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
)

type IndexAction actions.Action

func (action *IndexAction) Run(params struct {
	Auth *helpers.UserMustAuth
}) {
	action.Data["teaMenu"] = "proxy"
	action.Show()
}
