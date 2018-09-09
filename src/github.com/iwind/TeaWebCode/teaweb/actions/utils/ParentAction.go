package utils

import (
	"github.com/iwind/TeaGo/actions"
)

type ParentAction actions.Action

func (this *ParentAction) URL(url string) string {
	return url
}
