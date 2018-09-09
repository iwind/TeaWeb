package proxy

import (
	"github.com/iwind/TeaWebCode/teaweb/helpers"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/files"
)

type IndexAction struct {
	ParentAction
}

func (this *IndexAction) Run(params struct {
	Auth *helpers.UserMustAuth
}) {
	servers := []maps.Map{}

	dir := files.NewFile(Tea.ConfigDir())
	subFiles := dir.Glob("*.proxy.conf")
	files.Sort(subFiles, files.SortTypeModifiedTimeReverse)
	for _, configFile := range subFiles {
		reader, err := configFile.Reader()
		if err != nil {
			logs.Error(err)
			continue
		}

		config := &teaconfigs.ServerConfig{}
		err = reader.ReadYAML(config)
		if err != nil {
			continue
		}
		servers = append(servers, maps.Map{
			"config":   config,
			"filename": configFile.Name(),
		})
	}

	this.Data["servers"] = servers

	this.Show()
}
