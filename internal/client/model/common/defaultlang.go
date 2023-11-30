package common

import "imsdk/pkg/app"

func GetDefaultLang() string {
	var langConf struct {
		Languages []string `toml:"languages"`
		Default   string   `toml:"default"`
	}
	err := app.Config().Bind("global", "languages", &langConf)
	if err != nil {
		return ""
	}
	return langConf.Default
}
