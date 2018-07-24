package teaconfigs

type AccessLogStdoutConfig struct {
	Format string `yaml:"format"`
	Buffer string `yaml:"buffer"` // @TODO
	Flush  string `yaml:"flush"`  // @TODO
}
