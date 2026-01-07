package user

import (
	"web-chat/api/http_model"
	"web-chat/internal/svc"
)

type Logic interface {
	Register(req *http_model.UserRegisterReq) error
	Login(req *http_model.LoginReq) (string, error)
<<<<<<< Updated upstream
	Update(userID int64, req *http_model.UserUpdateReq) error
	Logout(req *http_model.LogoutReq) error
=======
	LoginByCode(req *http_model.LoginCodeReq) (string, error)
	SendEmailCode(req *http_model.SendEmailCodeReq) error
	Update(userID string, req *http_model.UserUpdateReq) error
	Logout(req *http_model.LogoutReq) error
	GetUserInfo(userID string) (*http_model.UserInfoResp, error)
>>>>>>> Stashed changes
}

type logicImpl struct {
	svcCtx *svc.Context
}

func NewLogic(svcCtx *svc.Context) Logic {
	return &logicImpl{svcCtx: svcCtx}
}
