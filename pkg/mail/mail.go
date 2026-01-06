package mail

type Manager interface {
	Send(msg Message, targets []string) error
}

type Message struct {
	Title    string
	Content  any
	Appendix string
}
