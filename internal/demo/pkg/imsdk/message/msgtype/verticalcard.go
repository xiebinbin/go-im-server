package msgtype

import "imsdk/internal/demo/pkg/imsdk/resource"

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

func NewCardVerticalButtons(items ...VerticalCardButtonsItems) *VerticalCardButtons {
	return &items
}

func VerticalButtons(buttonId, txt, enableColor, disableColor string) *VerticalCardButtonsItems {
	return &VerticalCardButtonsItems{
		ButtonId:     buttonId,
		TXT:          txt,
		EnableColor:  enableColor,
		DisableColor: disableColor,
	}
}

func NewCardVertical(title VerticalCardTitle, text VerticalCardText, image *resource.Image, buttons *VerticalCardButtons) *VerticalCard {
	return &VerticalCard{
		Title:   title,
		Text:    text,
		Image:   image,
		Buttons: *buttons,
	}
}
