package typestruct

import "imsdk/internal/client/model/resource"

type Card struct {
	Icon    *resource.Image `json:"icon"`
	Text    []CardText      `json:"text"`
	Buttons []CardButtons   `json:"buttons"`
}

type CardText struct {
	Text  string `json:"t"`
	Color string `json:"color"`
	Value string `json:"value"`
}

type CardButtons struct {
	TXT          string `json:"txt"`
	EnableColor  string `json:"enableColor"`
	DisableColor string `json:"disableColor"`
	ButtonId     string `json:"buttonId"`
}
