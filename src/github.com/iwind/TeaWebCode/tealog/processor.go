package tealog

type Processor interface {
	Process(accessLog *AccessLog)
}
