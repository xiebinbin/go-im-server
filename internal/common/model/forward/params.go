package forward

type MsgItem struct {
	ID         string                 `json:"id"`
	Amid       string                 `json:"amid"`
	Type       int8                   `json:"type"`
	ChatId     string                 `json:"chat_id"`
	SenderId   string                 `json:"from_address"`
	Content    string                 `json:"content"`
	Extra      string                 `json:"extra"`
	Action     interface{}            `json:"action,omitempty"`
	Status     int8                   `json:"status"`
	SendAt     int64                  `json:"send_time"`
	Number     int64                  `json:"number"`
	CreatedAt  int64                  `json:"create_time"`
	ReceiveIds []string               `json:"receive_ids,omitempty"`
	UpdatedAt  int64                  `json:"update_time"`
	Offline    map[string]interface{} `json:"offline"`
	Sequence   int64                  `json:"sequence"`
}
