package conf

import (
	"github.com/basset-la/tools/env"
)

var properties *Props

func GetProps() Props {
	if properties != nil {
		return *properties
	}

	properties = &Props{}

	env.LoadProperties("config", properties)

	return *properties
}

type Props struct {
	SecretKey string `yaml:"secretKey"`
	Path      string `yaml:"path"`
	Database  struct {
		URI        string `json:"uri"`
		DB         string `json:"db"`
		Collection string `json:"collection"`
	} `yaml:"database"`
}
