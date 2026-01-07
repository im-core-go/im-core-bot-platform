package openai

import (
	"strings"
	httpmodel "web-chat/api/http_model/chat"
)

type persistedStream struct {
	inner      httpmodel.MessageSteam
	onComplete func(content string) error
	builder    strings.Builder
	done       bool
	ctx        *streamContext
}

func newPersistedStream(inner httpmodel.MessageSteam, onComplete func(content string) error) httpmodel.MessageSteam {
	return &persistedStream{
		inner:      inner,
		onComplete: onComplete,
	}
}

func (p *persistedStream) Next() (httpmodel.StreamEvent, bool, error) {
	ev, done, err := p.inner.Next()
	if err != nil {
		return ev, done, err
	}
	if ev.Type == httpmodel.EventTextDelta {
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
