package openai

import (
	"encoding/json"
	"errors"
	"strings"
	"web-chat/api/http_model/chat"
	http2 "web-chat/pkg/http"
	"web-chat/pkg/logger"
)

type responsesStream struct {
	sr *http2.SSEReader
}

type ResponseEvent struct {
	Type  string `json:"type"`
	Delta string `json:"delta,omitempty"`
}

func newOpenAIResponsesStream(sr *http2.SSEReader) *responsesStream {
	return &responsesStream{sr: sr}
}

func (s *responsesStream) Close() error { return s.sr.Close() }

func (s *responsesStream) Next() (chat.StreamEvent, bool, error) {
	for {
		ev, ok, err := s.sr.Next()
		if err != nil {
			return chat.StreamEvent{}, false, err
		}
		if !ok {
			return chat.StreamEvent{Type: chat.EventDone}, true, nil
		}
		if strings.TrimSpace(ev.Data) == "" {
			continue
		}

		var e ResponseEvent
		if err := json.Unmarshal([]byte(ev.Data), &e); err != nil {
			logger.L().Errorf("unmarshal response event error: %v", err)
			return chat.StreamEvent{Type: chat.EventError}, false, errors.New("invalid stream json")
		}

		switch e.Type {
		case "response.output_text.delta":
			if e.Delta != "" {
				return chat.StreamEvent{Type: chat.EventTextDelta, Delta: e.Delta}, false, nil
			}
		case "response.completed":
			return chat.StreamEvent{Type: chat.EventDone}, true, nil
		default:
		}
	}
}
