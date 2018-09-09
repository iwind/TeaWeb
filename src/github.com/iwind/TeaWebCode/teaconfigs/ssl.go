package teaconfigs

type SSLConfig struct {
	On             bool   `yaml:"on" json:"on"`
	Certificate    string `yaml:"certificate" json:"certificate"`
	CertificateKey string `yaml:"certificateKey" json:"certificateKey"`
}
