package teaproxy

func Shutdown() {
	for _, listener := range LISTENERS {
		listener.Shutdown()
	}

	LISTENERS = []*Listener{}
}
