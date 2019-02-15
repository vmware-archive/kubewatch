package icq

type Response struct {
	Response struct {
		StatusCode int              `json:"statusCode"`
		StatusText string           `json:"statusText"`
		RequestId  string           `json:"requestId"`
		Data       *MessageResponse `json:"data"`
	} `json:"response"`
}

type MessageResponse struct {
	SubCode struct {
		Error int `json:"error"`
	} `json:"subCode"`
	MessageID        string `json:"msgId"`
	HistoryMessageID int64  `json:"histMsgId"`
	State            string `json:"state"`
}

type WebhookRequest struct {
	Token   string   `json:"aimsid"`
	Updates []Update `json:"update"`
}

type Update struct {
	Update struct {
		Chat Chat   `json:"chat"`
		Date int    `json:"date"`
		From User   `json:"from"`
		Text string `json:"text"`
	} `json:"update"`
	UpdateID int `json:"update_id"`
}

type Chat struct {
	ID string `json:"id"`
}

type User struct {
	ID           string `json:"id"`
	LanguageCode string `json:"language_code"`
}
