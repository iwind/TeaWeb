package teastats

import (
	"github.com/iwind/TeaWebCode/tealogs"
)

func init() {
	tealogs.SharedLogger().AddProcessor(new(Processor))
}
