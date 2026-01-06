package chat

type StreamEventType string

const (
	EventTextDelta StreamEventType = "text.delta"
	EventImage     StreamEventType = "image"
	EventDone      StreamEventType = "done"
	EventError     StreamEventType = "error"
)

type StreamEvent struct {
	Type  StreamEventType `json:"type"`
	Delta string          `json:"delta,omitempty"`
}

type MessageSteam interface {
	Next() (StreamEvent, bool, error)
	Close() error
}
