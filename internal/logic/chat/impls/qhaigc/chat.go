package qhaigc

import (
	"encoding/json"
	"net/http"
	"time"
	"web-chat/internal/logic/chat"
	"web-chat/internal/svc"
)

type logicImpl struct {
	svcCtx *svc.Context
	client *http.Client
}

func (l *logicImpl) HandleSingleMessage() error {
	//TODO implement me
	panic("implement me")
}

func (l *logicImpl) PullModules() ([]string, error) {
	var (
		modules []string
		err     error
	)
	modules, err = l.handlePullModules()
	if err != nil {
		return nil, err
	}
	return modules, nil
}
func (l *logicImpl) handlePullModules() ([]string, error) {
	var (
		req *http.Request
		err error
	)
	req, err = http.NewRequest(http.MethodGet, PullModuleURL, nil)
	if err != nil {
		return nil, err
	}

	type Resp struct {
		Data []struct {
			Id string `json:"id"`
		} `json:"data"`
	}
	resp, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}
	var r Resp
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	modules := make([]string, 0, len(r.Data))
	for _, m := range r.Data {
		modules = append(modules, m.Id)
	}
	return modules, nil
}

func NewChatLogic(svcCtx *svc.Context) chat.Logic {
	return &logicImpl{
		svcCtx: svcCtx,
		client: &http.Client{Timeout: time.Second * 5},
	}
}
