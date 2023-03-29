package gchat

type Event struct {
	Type    string `json:"type"`
	Message struct {
		Text         string `json:"text"`
		ArgumentText string `json:"argumentText"`
	}
}

type Response struct {
	Text string `json:"text"`
}
