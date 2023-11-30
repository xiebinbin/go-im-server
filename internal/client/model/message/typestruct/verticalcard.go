package typestruct

import "imsdk/internal/client/model/resource"

type VerticalCardButtons = []VerticalCardButtonsItems
type VerticalCard struct {
	Title   VerticalCardTitle   `json:"title"`
	Text    VerticalCardText    `json:"text"`
	Image   *resource.Image     `json:"image"`
	Buttons VerticalCardButtons `json:"buttons"`
}

type VerticalCardButtonsItems struct {
	TXT          string `json:"txt"`
	EnableColor  string `json:"enableColor"`
	DisableColor string `json:"disableColor"`
	ButtonId     string `json:"buttonId"`
}

type VerticalCardTitle struct {
	Type  string `json:"t"`
	Color string `json:"color"`
	Value string `json:"value"`
}

type VerticalCardText struct {
	Type  string `json:"t"`
	Color string `json:"color"`
	Value string `json:"value"`
}