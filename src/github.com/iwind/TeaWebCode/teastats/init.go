package teastats

import (
	"github.com/iwind/TeaWebCode/tealog"
)

func init() {
	tealog.SharedLogger().AddProcessor(new(Processor))
}
