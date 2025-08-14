package conf

type Config interface {
	LoadConf() *Options
	SetConf(bool, int, int, int)
}

type Options struct {
	keepAlive bool
	headerReadTimeout int
	headerWriterTimeout int
	httpOperationsTimeout int
	idleConnTimeout int
}
