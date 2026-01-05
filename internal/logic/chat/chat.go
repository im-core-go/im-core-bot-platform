package chat

type Logic interface {
	//	HandleSingleMessage() error
	PullModules() ([]string, error)
}
