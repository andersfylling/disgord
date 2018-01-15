package disgord

type GatewayPayload struct {
	Op             uint        `json:"op"`
	Data           interface{} `json:"d"`
	SequenceNumber uint        `json:"s,omitempty"`
	EventName      string      `json:"t,omitempty"`
}

type GetGatewayResponse struct {
	URL    string `json:"url"`
	Shards uint   `json:"shards,omitempty"`
}
