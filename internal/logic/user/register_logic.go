package user

import (
	"context"
	"fmt"
	"web-chat/api/http_model"
	"web-chat/internal/model"
	"web-chat/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

func (l *logicImpl) Register(req *http_model.UserRegisterReq) error {
	var (
		err    error
		ok     bool
		entity model.User
	)
	lgr := logger.L()
	if ok, err = l.svcCtx.Utils.Regexp.ValidateEmail(req.Email); !ok || err != nil {
		lgr.Errorf("register email invalid: %v", err)
		return fmt.Errorf("email is invalid")
	}
	ok, err = l.svcCtx.Utils.Code.VerifyCode(context.Background(), req.EmailCode, req.Email)
	if err != nil {
		lgr.Errorf("register verify code error: %v", err)
		return err
	}
	if !ok {
		return fmt.Errorf("email code is invalid")
	}
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		lgr.Errorf("register hash error: %v", err)
		return err
	}
	if req.Phone != nil && *req.Phone != "" {
		if ok, err = l.svcCtx.Utils.Regexp.ValidatePhone(*req.Phone); !ok || err != nil {
			lgr.Errorf("register phone invalid: %v", err)
			return fmt.Errorf("phone number is invalid")
		}
		entity.Phone = req.Phone
	}
	entity.UUID = l.svcCtx.Utils.UUID.New()
	entity.Email = req.Email

	entity.Password = string(hash)
	entity.NickName = req.NickName

	err = l.svcCtx.Dao.UserDao.CreateUser(entity)
	if err != nil {
		lgr.Errorf("register create user error: %v", err)
		return err
	}
	return nil
}
