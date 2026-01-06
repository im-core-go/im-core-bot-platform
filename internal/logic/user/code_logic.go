package user

import (
	"context"
	"fmt"
	"web-chat/api/http_model"
	"web-chat/pkg/logger"
)

func (l *logicImpl) SendEmailCode(req *http_model.SendEmailCodeReq) error {
	if req == nil {
		return fmt.Errorf("send code request is nil")
	}
	ok, err := l.svcCtx.Utils.Regexp.ValidateEmail(req.Email)
	if err != nil {
		return err
	}
	if !ok {
		if okPhone, err := l.svcCtx.Utils.Regexp.ValidatePhone(req.Email); err == nil && okPhone {
			return fmt.Errorf("sms is not supported")
		}
		return fmt.Errorf("email is invalid")
	}
	if err := l.svcCtx.Utils.Code.SendEmailCode(context.Background(), req.Email); err != nil {
		logger.L().Errorf("send email code error: %v", err)
		return err
	}
	return nil
}
