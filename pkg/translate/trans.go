package translate

import "imsdk/pkg/app"

type (
	GoogleV2Conf struct {
		Key string `toml:"key"`
	}
	config struct {
		GoogleV2Conf GoogleV2Conf `toml:"googleV2"`
	}
)

var conf config

func GetConf() config {
	app.Config().Bind("global", "translate", &conf)
	return conf
}
