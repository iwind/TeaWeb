package ssl

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/teaproxy"
)

type UploadCertAction actions.Action

func (this *UploadCertAction) Run(params struct {
	Auth     *helpers.UserMustAuth
	Filename string
	CertFile *actions.File
}) {
	// @TODO 校验证书文件格式

	if params.CertFile == nil {
		this.Fail("请选择证书文件")
	}

	data, err := params.CertFile.Read()
	if err != nil {
		this.Fail(err.Error())
	}

	certFilename := stringutil.Rand(16) + params.CertFile.Ext
	configFile := files.NewFile(Tea.ConfigFile(certFilename))
	err = configFile.Write(data)
	if err != nil {
		this.Fail(err.Error())
	}

	server, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		configFile.Delete()
		this.Fail(err.Error())
	}

	if server.SSL == nil {
		server.SSL = new(teaconfigs.SSLConfig)
	}

	server.SSL.Certificate = certFilename
	server.WriteToFilename(params.Filename)

	if server.SSL.On && len(server.SSL.CertificateKey) > 0 {
		teaproxy.Restart()
	}

	this.Refresh().Success("保存成功")
}
