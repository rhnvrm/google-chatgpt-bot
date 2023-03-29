package gchat

type DeprecatedEvent struct {
	Message struct {
		Text string `json:"text"`
	}
}

type Response struct {
	Text string `json:"text"`
}
