package teaproxy

import "testing"

func TestStart(t *testing.T) {
	startProxies()
	Wait()
}
