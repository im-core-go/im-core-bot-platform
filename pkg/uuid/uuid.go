package uuid

import "github.com/google/uuid"

type Wrap struct {
}

func NewWrap() *Wrap {
	return &Wrap{}
}
func (w *Wrap) New() string {
	return uuid.New().String()
}
