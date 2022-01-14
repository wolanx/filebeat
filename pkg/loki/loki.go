package loki

import (
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/outputs"
)

//var (
//	logger = logp.NewLogger("output.loki")
//)

func init() {
	outputs.RegisterType("loki", newLoki)
}

func newLoki(_ outputs.IndexManager, _ beat.Info, observer outputs.Observer, cfg *common.Config) (outputs.Group, error) {
	config := defaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return outputs.Fail(err)
	}

	hosts, err := outputs.ReadHostList(cfg)
	if err != nil {
		return outputs.Fail(err)
	}

	clients := make([]outputs.NetworkClient, len(hosts))
	for i, host := range hosts {
		clients[i] = &lokiClient{
			host:     host,
			isHttp:   config.Protocol == "http",
			observer: observer,
		}
	}

	return outputs.SuccessNet(false, 200, 3, clients)
}
