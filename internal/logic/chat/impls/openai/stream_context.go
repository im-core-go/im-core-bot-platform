package openai

import (
	"github.com/im-core-go/im-core-bot-platform/internal/logic/chat"
	"sync"
)

type streamContext struct {
	conversationID string
	title          string
	mu             sync.Mutex
}

func (s *streamContext) setTitle(title string) {
	s.mu.Lock()
	s.title = title
	s.mu.Unlock()
}

func (s *streamContext) getTitle() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.title
}

func (l *logicImpl) setStreamContext(conversationID string, stream chat.MessageStream) {
	ps, ok := stream.(*persistedStream)
	if !ok {
		return
	}
	ps.ctx = &streamContext{conversationID: conversationID}
}

func (l *logicImpl) setStreamTitle(stream chat.MessageStream, title string) {
	ps, ok := stream.(*persistedStream)
	if !ok || ps.ctx == nil {
		return
	}
	ps.ctx.setTitle(title)
}
