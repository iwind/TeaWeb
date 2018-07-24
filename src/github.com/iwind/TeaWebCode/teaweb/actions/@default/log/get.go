package log

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/tealog"
	"github.com/pquerna/ffjson/ffjson"
)

type GetAction actions.Action

func (action *GetAction) Run() {
	accessLogs := tealog.SharedLogger().ReadLogs(0, 5)
	data, err := ffjson.Marshal(accessLogs)
	if err != nil {
		action.WriteString(err.Error())
	} else {
		action.Write(data)
	}
}
