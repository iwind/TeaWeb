package teaconfigs

import "testing"

func TestParseConfigs(t *testing.T) {
	configs, err := ParseConfigs()
	if err != nil {
		t.Error(err)
		return
	}
	for _, config := range configs {
		t.Log(config.Address, config.Servers)
	}

	if len(configs) > 0 {
		config := configs[0]
		if len(config.Servers) > 0 {
			for _, location := range config.Servers[0].Locations {
				t.Logf("Location: %#v", *location)
			}
		}
	}
}
