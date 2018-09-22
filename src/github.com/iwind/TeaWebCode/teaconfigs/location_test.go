package teaconfigs

import (
	"testing"
	"github.com/iwind/TeaGo/assert"
)

func TestLocationConfig_Match(t *testing.T) {
	location := NewLocationConfig()
	err := location.Validate()
	if err != nil {
		t.Fatal(err)
	}

	a := assert.NewAssertion(t).Quiet()

	location.Pattern = "/hell"
	a.IsNotError(location.Validate())
	a.IsTrue(location.Match("/hello"))

	location.Pattern = "/hello"
	a.IsNotError(location.Validate())
	a.IsTrue(location.Match("/hello"))

	location.Pattern = "~ ^/\\w+$"
	a.IsNotError(location.Validate())
	a.IsTrue(location.Match("/hello"))

	location.Pattern = "!~ ^/HELLO$"
	a.IsNotError(location.Validate())
	a.IsTrue(location.Match("/hello"))

	location.Pattern = "~* ^/HELLO$"
	a.IsNotError(location.Validate())
	a.IsTrue(location.Match("/hello"))

	location.Pattern = "!~* ^/HELLO$"
	a.IsNotError(location.Validate())
	a.IsFalse(location.Match("/hello"))

	location.Pattern = "= /hello"
	a.IsNotError(location.Validate())
	a.IsTrue(location.Match("/hello"))
}
