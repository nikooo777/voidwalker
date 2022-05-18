package configs

import (
	"time"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/tkanos/gonfig"
)

type DbConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Database string `json:"database"`
	Password string `json:"password"`
}
type Configs struct {
	RedirectBaseURL     string        `json:"redirect_base_url"`
	RedirectAdvancedURL string        `json:"redirect_advanced_url"`
	Chainquery          DbConfig      `json:"chainquery"`
	Speech              DbConfig      `json:"speech"`
	Voidwalker          DbConfig      `json:"voidwalker"`
	ChannelID           string        `json:"channel_id"`
	PublishAddress      string        `json:"publish_address"`
	ReflectorServer     string        `json:"reflector_server"`
	LbrynetTimeout      time.Duration `json:"lbrynet_timeout"`
	PreviousChannelIds  []string      `json:"previous_channel_ids"`
}

var Configuration *Configs

func Init(configPath string) error {
	if Configuration != nil {
		return nil
	}
	c := Configs{}
	err := gonfig.GetConf(configPath, &c)
	if err != nil {
		return errors.Err(err)
	}
	Configuration = &c
	return nil
}
