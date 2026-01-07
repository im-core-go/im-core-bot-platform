package user

import (
	"web-chat/api/http_model"
	"web-chat/internal/svc"
)

type Logic interface {
	Register(req *http_model.UserRegisterReq) error
	Login(req *http_model.LoginReq) (string, error)
	LoginByCode(req *http_model.LoginCodeReq) (string, error)
	SendEmailCode(req *http_model.SendEmailCodeReq) error
	Update(userID string, req *http_model.UserUpdateReq) error
	Logout(req *http_model.LogoutReq) error
	GetUserInfo(userID string) (*http_model.UserInfoResp, error)
}

type logicImpl struct {
	svcCtx *svc.Context
}

func NewLogic(svcCtx *svc.Context) Logic {
	return &logicImpl{svcCtx: svcCtx}
}
