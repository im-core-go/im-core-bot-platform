package openai

import (
	"github.com/im-core-go/im-core-bot-platform/internal/logic/chat"
	"strings"
)

type persistedStream struct {
	inner      chat.MessageStream
	onComplete func(content string) error
	builder    strings.Builder
	done       bool
	ctx        *streamContext
}

func newPersistedStream(inner chat.MessageStream, onComplete func(content string) error) chat.MessageStream {
	return &persistedStream{
		inner:      inner,
		onComplete: onComplete,
	}
}

func (p *persistedStream) Next() (chat.StreamEvent, bool, error) {
	ev, done, err := p.inner.Next()
	if err != nil {
		return ev, done, err
	}
	if ev.Type == chat.EventTextDelta {
		p.builder.WriteString(ev.Delta)
	}
	if done {
		p.flushOnce()
		if p.ctx != nil {
			ev.ConversationID = p.ctx.conversationID
			ev.Title = p.ctx.getTitle()
		}
	}
	return ev, done, nil
}

func (p *persistedStream) Close() error {
	p.flushOnce()
	return p.inner.Close()
}

func (p *persistedStream) flushOnce() {
	if p.done || p.onComplete == nil {
		return
	}
	p.done = true
	_ = p.onComplete(p.builder.String())
}
