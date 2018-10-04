package teaplugins

var plugins = []*Plugin{}

func Register(plugin *Plugin) {
	plugins = append(plugins, plugin)
}

func TopBarWidgets() []*Widget {
	result := []*Widget{}
	for _, plugin := range plugins {
		for _, widget := range plugin.Widgets {
			if widget.TopBar {
				result = append(result, widget)
			}
		}
	}
	return result
}

func MenuBarWidgets() []*Widget {
	result := []*Widget{}
	for _, plugin := range plugins {
		for _, widget := range plugin.Widgets {
			if widget.MenuBar {
				result = append(result, widget)
			}
		}
	}
	return result
}

func HelperBarWidgets() []*Widget {
	result := []*Widget{}
	for _, plugin := range plugins {
		for _, widget := range plugin.Widgets {
			if widget.HelperBar {
				result = append(result, widget)
			}
		}
	}
	return result
}

func DashboardWidgets() []*Widget {
	result := []*Widget{}
	for _, plugin := range plugins {
		for _, widget := range plugin.Widgets {
			if widget.Dashboard {
				result = append(result, widget)
			}
		}
	}
	return result
}
