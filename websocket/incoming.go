package websocket

type payloadData []byte

func (pd *payloadData) UnmarshalJSON(data []byte) error {
	*pd = payloadData(data)
	return nil
}

func (pd *payloadData) ByteArr() []byte {
	return []byte(*pd)
}

type gatewayEvent struct {
	Op             uint        `json:"op"`
	Data           payloadData `json:"d"`
	SequenceNumber uint        `json:"s"`
	EventName      string      `json:"t"`
}

func (ge *gatewayEvent) GetOperationCode() uint {
	return ge.Op
}

type getGatewayResponse struct {
	URL    string `json:"url"`
	Shards uint   `json:"shards,omitempty"`
}

type helloPacket struct {
	HeartbeatInterval uint     `json:"heartbeat_interval"`
	Trace             []string `json:"_trace"`
}

type readyPacket struct {
	SessionID string   `json:"session_id"`
	Trace     []string `json:"_trace"`
}

type DiscordWSEvent interface {
	Name() string
	Data() []byte
	Unmarshal(interface{}) error
}

type Event struct {
	content *gatewayEvent
}

func (evt *Event) Name() string {
	return evt.content.EventName
}

func (evt *Event) Data() []byte {
	return evt.content.Data.ByteArr()
}

func (evt *Event) Unmarshal(v interface{}) error {
	return unmarshal(evt.Data(), v)
}
