package stat

import "github.com/iwind/TeaGo/actions"

type IndexAction actions.Action

func (this *IndexAction) Run(params struct{}) {
	this.Show()
}
