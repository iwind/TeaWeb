package teaconfigs

type AccessLogPipeConfig struct {
	Path   string `yaml:"path"`
	Format string `yaml:"format"` // @TODO
	Buffer string `yaml:"buffer"` // @TODO
	Flush  string `yaml:"flush"`  // @TODO
}
