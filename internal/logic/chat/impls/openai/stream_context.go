package openai

import (
	"sync"
	httpmodel "web-chat/api/http_model/chat"
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

func (l *logicImpl) setStreamContext(conversationID string, stream httpmodel.MessageSteam) {
	ps, ok := stream.(*persistedStream)
	if !ok {
		return
	}
	ps.ctx = &streamContext{conversationID: conversationID}
}

func (l *logicImpl) setStreamTitle(stream httpmodel.MessageSteam, title string) {
	ps, ok := stream.(*persistedStream)
	if !ok || ps.ctx == nil {
		return
	}
	ps.ctx.setTitle(title)
}
