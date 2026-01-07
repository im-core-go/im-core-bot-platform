package user

import (
	"fmt"
	"web-chat/api/http_model"
	"web-chat/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

func (l *logicImpl) Update(userID string, req *http_model.UserUpdateReq) error {
	if userID == "" {
		return fmt.Errorf("user id is required")
	}
	if req == nil {
		return fmt.Errorf("update request is nil")
	}
	lgr := logger.L()
	updateMap := make(map[string]interface{})
	if req.NickName != nil {
		updateMap["nick_name"] = *req.NickName
	}
	if req.Phone != nil {
		ok, err := l.svcCtx.Utils.Regexp.ValidatePhone(*req.Phone)
		if err != nil || !ok {
			lgr.Errorf("update phone invalid: %v", err)
			return fmt.Errorf("phone number is invalid")
		}
		updateMap["phone"] = *req.Phone
	}
	if req.Email != nil {
		ok, err := l.svcCtx.Utils.Regexp.ValidateEmail(*req.Email)
		if err != nil || !ok {
			lgr.Errorf("update email invalid: %v", err)
			return fmt.Errorf("email is invalid")
		}
		updateMap["email"] = *req.Email
	}
	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword(
			[]byte(*req.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			lgr.Errorf("update hash error: %v", err)
			return err
		}
		updateMap["password"] = string(hash)
	}
	if len(updateMap) == 0 {
		return fmt.Errorf("no fields to update")
	}
	if err := l.svcCtx.Dao.UserDao.UpdateUser(userID, updateMap); err != nil {
		lgr.Errorf("update user error: %v", err)
		return err
	}
	return nil
}
