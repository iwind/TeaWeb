package teaconfigs

type SSLConfig struct {
	On             bool   `yaml:"on"`
	Certificate    string `yaml:"certificate"`
	CertificateKey string `yaml:"certificateKey"`
}
