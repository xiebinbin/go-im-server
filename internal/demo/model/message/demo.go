package message

import (
	"imsdk/internal/client/model/message/typestruct"
	"imsdk/internal/common/dao/message/detail"
	"imsdk/internal/demo/pkg/imsdk/resource"
)

func CreateType1Content(text string) detail.Type1Struct {
	data := detail.Type1Struct{
		Text: text,
	}
	return data
}

func CreateType2Content() detail.Type2Struct {
	data := detail.Type2Struct{
		Height:   3840,
		IsOrigin: 1,
		Size:     1718995,
		Width:    2160,
		BucketId: "49ba59abbe56e057",
		FileType: "jpg",
		ObjectId: "img/fd/fd6e2f8f49816b11f3eb6093f04aac0fc2a0c3d0",
	}
	return data
}

func CreateType5Content() detail.Type5Struct {
	data := detail.Type5Struct{
		Size:     1718995,
		BucketId: "49ba59abbe56e057",
		FileType: "jpg",
		ObjectId: "img/fd/fd6e2f8f49816b11f3eb6093f04aac0fc2a0c3d0",
	}
	return data
}

func CreateType19Content() detail.Type19Struct {
	var dataIcon detail.TypeIcon
	dataText := make([]detail.Type19Text, 0)
	dataButtons := make([]detail.Type19Buttons, 0)
	dataIcon = detail.TypeIcon{
		BucketId: "49ba59abbe56e057",
		FileType: "webp",
		Height:   126,
		Width:    126,
		Text:     "img/e4/e410959b8b1ac9b5043bd40d99d92b768350f990",
	}
	dataText = []detail.Type19Text{
		{
			T:     "txt",
			Color: "#000000",
			Value: "Would you like to send me your attached resume?",
		},
	}
	dataButtons = []detail.Type19Buttons{
		{
			TXT:          "Decline",
			EnableColor:  "#ABCDEF",
			DisableColor: "#EEEEEE",
			ButtonId:     "{\"id\": \"btn_1\",\"type\":19}",
		},
		{
			TXT:          "Agree",
			EnableColor:  "#ABCDEF",
			DisableColor: "#EEEEEE",
			ButtonId:     "{\"id\": \"btn_2\",\"type\":19}",
		},
	}
	data := detail.Type19Struct{
		Icon:    dataIcon,
		Text:    dataText,
		Buttons: dataButtons,
	}
	return data
}

func CreateType23Content() typestruct.VerticalCard {
	dataButtons := make(typestruct.VerticalCardButtons, 0)
	//dataIcon = detail.TypeIcon{
	//	BucketId: "49ba59abbe56e057",
	//	FileType: "webp",
	//	Height:   126,
	//	Width:    126,
	//	Text:     "img/e4/e410959b8b1ac9b5043bd40d99d92b768350f990",
	//}
	dataTitle := typestruct.VerticalCardTitle{
		Type:  "txt",
		Color: "#0D1324",
		Value: "How to add friends?",
	}

	dataText := typestruct.VerticalCardText{
		Type:  "txt",
		Color: "#000000",
		Value: "Would you like to send me your attached resume?",
	}
	dataButtons = []typestruct.VerticalCardButtonsItems{
		{
			TXT:          "Scan code",
			EnableColor:  "#00C6DB",
			DisableColor: "#eeeeee",
			ButtonId:     "https://sq",
		},
		{
			TXT:          "My QR Code",
			EnableColor:  "#00C6DB",
			DisableColor: "#eeeeee",
			ButtonId:     "https://mqr",
		},
		{
			TXT:          "Add friends with phone",
			EnableColor:  "#00C6DB",
			DisableColor: "#eeeeee",
			ButtonId:     "https://us",
		},
		{
			TXT:          "Mobile Contacts",
			EnableColor:  "#00C6DB",
			DisableColor: "#eeeeee",
			ButtonId:     "https://mc",
		},
	}
	data := typestruct.VerticalCard{
		Title:   dataTitle,
		Text:    dataText,
		Image:   resource.NewOriginImage(map[string]interface{}{}),
		Buttons: dataButtons,
	}
	return data
}

func CreateType20TempInfo(operator string, target []string) detail.Type20Struct {
	data := detail.Type20Struct{
		TemId:    "jobseeker-no-receive-resume",
		Operator: operator,
		Target:   target,
		Duration: 0,
		Number:   0,
	}
	//dataByte, _ := json.Marshal(data)
	//return string(dataByte)
	return data
}

func CreateType22TempInfo() detail.Type22Struct {
	data := detail.Type22Struct{
		Title: "Customize message title",
		Body:  "Customize message body \t Hello <a id=\"bt1\" color=\"#ABCDEF\">World</a> I'm fine",
	}
	//dataByte, _ := json.Marshal(data)
	//return string(dataByte)
	return data
}

func CreateType21Notification() detail.Type21Struct {
	data := detail.Type21Struct{
		BId: "{\"id\":\"https://qr.td.com.tr/mi/AWEwMWQ3NGJkNWJhMjQ4MDC7\",\"type\":21}",
		Data: []detail.Type21Data{
			{
				Lang: "tr",
				Title: "\"Gök ve yer benimle yan yana yaşıyor ve her şey benimle bir\", " +
					"aşkınlığın anlamı terk etmek ve yükselmek, dünyanın değerini terk " +
					"etmek ve daha yüksek ve daha geniş bir manevi aleme yükseltmektir.",
				TXT: "Yetişkinlerin sayılara karşı zaafı vardır. Bir arkadaşınızla tanıştırırsanız," +
					" size asla \"Sesi nasıl? Hangi oyunları oynamayı sever? Kelebek tahnitçiliği " +
					"topluyor mu?\" diye sormuyorlar ama \"Kaç yaşında? Kaç kardeşi var? Kilosu " +
					"var mı?\" Babası ne kadar kazanıyor?\" Bunu bildikleri için adamı tanıdıklarını " +
					"sanıyorlar.",
			},
			{
				Lang: "en",
				Title: "\"Heaven and earth live side by side with me, and all things are one with me\"," +
					" the meaning of transcendence lies in abandoning and upgrading, " +
					"abandoning the value of the world, and upgrading to a higher and " +
					"broader spiritual realm",
				TXT: "Adults have a soft spot for numbers. If you introduce a friend to them, " +
					"they never ask you \"How is his voice? What games does he like to play? " +
					"Does he collect butterfly taxidermy?\" but \"How old is he? " +
					"How many brothers does he have? Weight?\" How much? How much does his father earn?\"" +
					" They think that knowing this, they know the man.",
			},
			{
				Lang:  "cn",
				Title: "“天地与我并生，万物与我为一”,超越的意义，在于扬弃与提升，扬弃俗世的价值，而提升到更高更辽阔的精神领域中",
				TXT: "成人们对数字情有独钟。如果你为他们介绍一个朋友，他们从不会问你“他的嗓子怎么样？" +
					"他爱玩什么游戏？他会采集蝴蝶标本嘛？”而是问“他几岁了？有多少个兄弟？体重多少？" +
					"他的父亲挣多少钱？”他们认为知道了这些，就了解了这个人。",
			},
		},
	}
	//dataByte, _ := json.Marshal(data)
	//return string(dataByte)
	return data
}
