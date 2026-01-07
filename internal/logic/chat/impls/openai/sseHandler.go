package openai

import (
	"encoding/json"
	"errors"
	"strings"
	"web-chat/api/http_model/chat"
	http2 "web-chat/pkg/http"
	"web-chat/pkg/logger"
)

type chatCompletionsStream struct {
	sr *http2.SSEReader
}

type ChatCompletionChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

func newOpenAIChatCompletionsStream(sr *http2.SSEReader) *chatCompletionsStream {
	return &chatCompletionsStream{sr: sr}
}

func (s *chatCompletionsStream) Close() error { return s.sr.Close() }

func (s *chatCompletionsStream) Next() (chat.StreamEvent, bool, error) {
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

		var e ChatCompletionChunk
		if err := json.Unmarshal([]byte(ev.Data), &e); err != nil {
			logger.L().Errorf("unmarshal response event error: %v", err)
			return chat.StreamEvent{Type: chat.EventError}, false, errors.New("invalid stream json")
		}

		if len(e.Choices) == 0 {
			continue
		}
		choice := e.Choices[0]
		if choice.Delta.Content != "" {
			return chat.StreamEvent{Type: chat.EventTextDelta, Delta: choice.Delta.Content}, false, nil
		}
		if choice.FinishReason != nil && *choice.FinishReason != "" {
			return chat.StreamEvent{Type: chat.EventDone}, true, nil
		}
	}
}
