package user

import (
	"fmt"
	"web-chat/api/http_model"
)

func (l *logicImpl) GetUserInfo(userID string) (*http_model.UserInfoResp, error) {
	if userID == "" {
		return nil, fmt.Errorf("user id is required")
	}
	entity, err := l.svcCtx.Dao.UserDao.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return &http_model.UserInfoResp{
		UUID:      entity.UUID,
		NickName:  entity.NickName,
		Email:     entity.Email,
		Phone:     entity.Phone,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}
