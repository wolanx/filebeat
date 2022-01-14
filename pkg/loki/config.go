package loki

type config struct {
	Protocol string `config:"protocol"` // http grpc
}

var (
	defaultConfig = config{}
)
